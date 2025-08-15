// Copyright (c) 2020-2024 Bryan Frimin <bryan@frimin.fr>.
//
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
// REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
// INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
// LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
// OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
// PERFORMANCE OF THIS SOFTWARE.

package privatebin

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"runtime"
	"time"
)

func btoi(v bool) int {
	if v {
		return 1
	}

	return 0
}

func itob(v int) bool {
	return v != 0
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func decode64(s string) ([]byte, error) {
	if len(s)%4 == 0 {
		return base64.StdEncoding.DecodeString(s)
	}

	return base64.RawStdEncoding.DecodeString(s)
}

func defaultPooledClient(tlsConfig *tls.Config) *http.Client {
	dial := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dial.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		TLSClientConfig:       tlsConfig,
	}

	return &http.Client{Transport: transport}
}

// Golang standard library does not expose GCM with custom nonce and
// tag size, even if it supported. Following code is a backport from
// the Golang crypto module to allowing it.
//
// References:
// - https://go-review.googlesource.com/c/go/+/116435
// - https://github.com/golang/go/issues?q=NewGCMWithNonceAndTagSize
// - https://github.com/golang/go/issues/42470

const (
	gcmBlockSize      = 16
	gcmMinimumTagSize = 12 // NIST SP 800-38D recommends tags with 12 or more bytes.
)

type gcmAble interface {
	NewGCM(nonceSize, tagSize int) (cipher.AEAD, error)
}

func newGCMWithNonceAndTagSize(block cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	if tagSize < gcmMinimumTagSize || tagSize > gcmBlockSize {
		return nil, errors.New("cipher: incorrect tag size given to GCM")
	}

	if nonceSize <= 0 {
		return nil, errors.New("cipher: the nonce can't have zero length, or the security of the key will be immediately compromised")
	}

	if block, ok := block.(gcmAble); ok {
		return block.NewGCM(nonceSize, tagSize)
	}
	// Attempt to use cipher.NewGCM() before giving up.
	// This matches the encrypt logic and works fine for some privatebin implementations.
	return cipher.NewGCM(block)
}
