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
	"encoding/json"
)

const (
	EncryptionAlgorithmUnknow EncryptionAlgorithm = iota
	EncryptionAlgorithmAES
)

type EncryptionAlgorithm uint8

func (ea EncryptionAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(ea.String())
}

func (ea *EncryptionAlgorithm) UnmarshalJSON(data []byte) error {
	var (
		v EncryptionAlgorithm
		s string
	)

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	switch s {
	case "aes":
		v = EncryptionAlgorithmAES
	default:
		v = EncryptionAlgorithmUnknow
	}

	*ea = v

	return nil
}

func (ea EncryptionAlgorithm) String() string {
	switch ea {
	case EncryptionAlgorithmAES:
		return "aes"
	default:
		return "unknown"
	}
}
