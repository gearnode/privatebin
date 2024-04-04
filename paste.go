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
	"fmt"
	"mime"
	"path/filepath"
)

type (
	Paste struct {
		Data           []byte
		Attachement    []byte
		AttachmentName string
	}
)

func (p Paste) MarshalJSON() ([]byte, error) {
	output := map[string]string{}

	if len(p.Attachement) > 0 {
		ext := filepath.Ext(p.AttachmentName)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		if p.AttachmentName != "" {
			output["attachment_name"] = p.AttachmentName
		}

		output["attachment"] = fmt.Sprintf(
			"data:%s;base64,%s",
			mimeType,
			base64.StdEncoding.EncodeToString(p.Attachement),
		)
	}

	if len(p.Data) > 0 {
		output["paste"] = string(p.Data)
	}

	return json.Marshal(output)
}

func (p *Paste) UnmarshalJSON(data []byte) error {
	output := map[string]string{}
	err := json.Unmarshal(data, &output)
	if err != nil {
		return err
	}

	*p = Paste{
		[]byte(output["paste"]),
		[]byte(""),
		output["attachment_name"],
	}

	return nil
}
