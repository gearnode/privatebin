// Copyright (c) 2020-2026 Bryan Frimin <bryan@frimin.fr>.
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
	"fmt"
)

const (
	gcmStandardNonceSize = 12
	gcmStandardTagSize   = 16
)

// newGCMWithNonceOrTagSize creates a GCM cipher with custom nonce or tag sizes.
//
// Historical context:
//   - Before Go 1.19: There was a public function cipher.NewGCMWithNonceAndTagSize()
//     that allowed creation of GCM with both custom nonce AND tag sizes.
//   - Go 1.19: The public NewGCMWithNonceAndTagSize was removed as a breaking change.
//     However, the internal newGCMWithNonceAndTagSize function and gcmAble interface
//     were still accessible, allowing a workaround to achieve both custom parameters.
//   - Go 1.21+: The FIPS 140-3 refactoring moved everything into internal packages
//     (crypto/internal/fips140), making the gcmAble interface completely inaccessible.
//     No public or hackable interface remains to create GCM with both custom parameters.
//
// Related issue: https://github.com/golang/go/issues/42470 (still open)
//
// Current limitations due to Go 1.21+ changes:
// - Standard parameters (12-byte nonce, 16-byte tag): Fully supported via cipher.NewGCM()
// - Custom nonce size only: Supported via cipher.NewGCMWithNonceSize()
// - Custom tag size only: Supported via cipher.NewGCMWithTagSize()
// - Both custom nonce AND tag sizes: No longer possible - returns an error
//
// This implementation is now limited to maintain compatibility with newer Go versions.
// Applications requiring both custom nonce and tag sizes must either use Go < 1.19
// or find alternative cryptographic libraries.
//
// Design decision:
// I chose not to fork or implement a custom GCM mode because this project is not a
// full-time commitment, making it impossible to invest the time needed to properly
// maintain such critical cryptographic code. Using the latest Go version is prioritized
// to keep this project secure and usable by everyone. Third-party cryptographic libraries
// are not an option as none currently exist for this use case, and even if one did,
// establishing proper trust for cryptographic code in a security-focused project like
// this would be challenging.
func newGCMWithNonceOrTagSize(block cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	// For standard parameters, use the standard implementation
	if nonceSize == gcmStandardNonceSize && tagSize == gcmStandardTagSize {
		return cipher.NewGCM(block)
	}

	if tagSize == gcmStandardTagSize {
		return cipher.NewGCMWithNonceSize(block, nonceSize)
	}

	if nonceSize == gcmStandardNonceSize {
		return cipher.NewGCMWithTagSize(block, tagSize)
	}

	// For non-standard parameters in Go 1.19+, we cannot proceed
	return nil, fmt.Errorf("custom GCM parameters (nonce=%d, tag=%d) are not supported since Go 1.20+", nonceSize, tagSize)
}
