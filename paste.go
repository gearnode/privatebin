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
	"errors"
	"fmt"
	"mime"
	"net/url"
	"path/filepath"
	"strings"
)

type (
	Paste struct {
		Data           []byte
		Attachement    []byte
		AttachmentName string
		MimeType       string
	}
)

func (p Paste) MarshalJSON() ([]byte, error) {
	output := map[string]string{}

	if len(p.Attachement) > 0 {
		mimeType := p.MimeType
		if mimeType == "" {
			ext := filepath.Ext(p.AttachmentName)
			mimeType = mime.TypeByExtension(ext)
			if p.MimeType == "" {
				mimeType = "application/octet-stream"
			}
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

	var attachment []byte
	var mimeType string
	attachmentURL, ok := output["attachment"]
	if ok {
		parsedURL, err := url.Parse(attachmentURL)
		if err != nil {
			return fmt.Errorf("invalid attachment: error parsing url: %w", err)
		}

		parts := strings.Split(parsedURL.Opaque, ",")
		if len(parts) != 2 {
			return errors.New("invalid attachment: invalid data URL format")
		}

		if !strings.HasSuffix(parts[0], ";base64") {
			return errors.New("invalid attachment: missing or invalid base64 encoding")
		}

		mimeType = strings.TrimPrefix(parts[0], ";base64")
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		attachment, err = decode64(parts[1])
		if err != nil {
			return fmt.Errorf("invalid attachment: cannot base64 decode data: %w", err)
		}

	}

	*p = Paste{
		[]byte(output["paste"]),
		attachment,
		output["attachment_name"],
		mimeType,
	}

	return nil
}
