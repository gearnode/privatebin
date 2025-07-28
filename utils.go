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

func defaultPooledClient() *http.Client {
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
	}

	return &http.Client{Transport: transport}
}

// Golang standard library does not expose GCM with custom nonce and
// tag size in Go 1.24 and later. The gcmAble interface was removed.
// For compatibility, we'll fallback to using the standard GCM implementation
// and validate parameters.
//
// References:
// - https://go-review.googlesource.com/c/go/+/116435
// - https://github.com/golang/go/issues?q=NewGCMWithNonceAndTagSize
// - https://github.com/golang/go/issues/42470

const (
	gcmBlockSize         = 16
	gcmMinimumTagSize    = 12 // NIST SP 800-38D recommends tags with 12 or more bytes.
	gcmStandardNonceSize = 12 // Standard GCM nonce size
	gcmStandardTagSize   = 16 // Standard GCM tag size
)

func newGCM(cipherBlock cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	if tagSize < gcmMinimumTagSize || tagSize > gcmBlockSize {
		return nil, errors.New("cipher: incorrect tag size given to GCM")
	}

	if nonceSize <= 0 {
		return nil, errors.New("cipher: the nonce can't have zero length, or the security of the key will be immediately compromised")
	}

	// In Go 1.24+, we need to use the standard cipher.NewGCM functions
	// Check if we can use standard GCM (nonce=12, tag=16)
	if nonceSize == gcmStandardNonceSize && tagSize == gcmStandardTagSize {
		return cipher.NewGCM(cipherBlock)
	}

	// For non-standard nonce sizes, use NewGCMWithNonceSize if available
	if tagSize == gcmStandardTagSize {
		return cipher.NewGCMWithNonceSize(cipherBlock, nonceSize)
	}

	// For non-standard tag sizes, use NewGCMWithTagSize if available
	if nonceSize == gcmStandardNonceSize {
		return cipher.NewGCMWithTagSize(cipherBlock, tagSize)
	}

	// If both nonce and tag sizes are non-standard, this is not supported
	// in Go 1.24+ through the standard library
	return nil, errors.New("cipher: custom nonce and tag sizes are not supported together in Go 1.24+")
}
