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
	"encoding/base64"
	"encoding/json"
)

type (
	AData struct {
		Spec             Spec
		Formatter        string
		OpenDiscussion   bool
		BurnAfterReading bool
	}

	Spec struct {
		IV          []byte
		Salt        []byte
		Iterations  int
		KeySize     int
		TagSize     int
		Algorithm   EncryptionAlgorithm
		Mode        EncryptionMode
		Compression CompressionAlgorithm
	}
)

func (adata AData) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		[4]any{
			adata.Spec,
			adata.Formatter,
			btoi(adata.OpenDiscussion),
			btoi(adata.BurnAfterReading),
		},
	)
}

func (adata *AData) UnmarshalJSON(data []byte) error {
	var values [4]json.RawMessage
	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}

	var spec Spec
	err = json.Unmarshal(values[0], &spec)
	if err != nil {
		return err
	}

	var (
		formatter                       string
		openDiscussion, burnAterReading int
	)

	err = json.Unmarshal(values[1], &formatter)
	if err != nil {
		return err
	}

	err = json.Unmarshal(values[2], &openDiscussion)
	if err != nil {
		return err
	}

	err = json.Unmarshal(values[3], &burnAterReading)
	if err != nil {
		return err
	}

	*adata = AData{spec, formatter, itob(openDiscussion), itob(burnAterReading)}

	return nil
}

func (spec Spec) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		[8]any{
			base64.RawStdEncoding.EncodeToString(spec.IV),
			base64.RawStdEncoding.EncodeToString(spec.Salt),
			spec.Iterations,
			spec.KeySize,
			spec.TagSize,
			spec.Algorithm,
			spec.Mode,
			spec.Compression,
		},
	)

}

func (spec *Spec) UnmarshalJSON(data []byte) error {
	var values [8]json.RawMessage
	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}

	var (
		encodedIv, encodedSalt       string
		iv, salt                     []byte
		iterations, keySize, tagSize int
		algorithm                    EncryptionAlgorithm
		mode                         EncryptionMode
		compression                  CompressionAlgorithm
	)

	err = json.Unmarshal(values[0], &encodedIv)
	if err != nil {
		return err
	}

	iv, err = base64.RawStdEncoding.DecodeString(encodedIv)
	if err != nil {
		return err
	}

	err = json.Unmarshal(values[1], &encodedSalt)
	if err != nil {
		return err
	}

	salt, err = base64.RawStdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return err
	}

	err = json.Unmarshal(values[2], &iterations)
	if err != nil {
		return err
	}

	err = json.Unmarshal(values[3], &keySize)
	if err != nil {
		return err
	}

	err = json.Unmarshal(values[4], &tagSize)
	if err != nil {
		return err
	}

	err = json.Unmarshal(values[5], &algorithm)
	if err != nil {
		return err
	}

	err = json.Unmarshal(values[6], &mode)
	if err != nil {
		return err
	}

	err = json.Unmarshal(values[7], &compression)
	if err != nil {
		return err
	}

	*spec = Spec{iv, salt, iterations, keySize, tagSize, algorithm, mode, compression}

	return nil
}
