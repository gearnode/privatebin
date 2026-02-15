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
	EncryptionModeUnknow EncryptionMode = iota
	EncryptionModeGCM
)

type EncryptionMode uint8

func (em EncryptionMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(em.String())
}

func (em *EncryptionMode) UnmarshalJSON(data []byte) error {
	var (
		v EncryptionMode
		s string
	)

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	switch s {
	case "gcm":
		v = EncryptionModeGCM
	default:
		v = EncryptionModeUnknow
	}

	*em = v

	return nil
}

func (em EncryptionMode) String() string {
	switch em {
	case EncryptionModeGCM:
		return "gcm"
	default:
		return "unknown"
	}
}
