// Copyright (c) 2020-2022 Bryan Frimin <bryan@frimin.fr>.
//
// Permission to use, copy, modify, and/or distribute this software for
// any purpose with or without fee is hereby granted, provided that the
// above copyright notice and this permission notice appear in all
// copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL
// WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE
// AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL
// DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR
// PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
// TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
// PERFORMANCE OF THIS SOFTWARE.

package privatebin // import "gearno.de/privatebin"

import (
	"bytes"
	"compress/flate"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"gearno.de/base58"
	pv "gearno.de/privatebin/internal/version"
	"golang.org/x/crypto/pbkdf2"
)

const (
	PrivateBinAPIVersion = 2
)

type Client struct {
	URL      url.URL
	Username string
	Password string
}

type CreatePasteRequest struct {
	V     int                    `json:"v"`
	AData []interface{}          `json:"adata"`
	Meta  CreatePasteRequestMeta `json:"meta"`
	CT    string                 `json:"ct"`
}

type CreatePasteRequestMeta struct {
	Expire string `json:"expire"`
}

type CreatePasteResponse struct {
	ID          string `json:"id"`
	Status      int    `json:"status"`
	Message     string `json:"message"`
	URL         string `json:"url"`
	DeleteToken string `json:"deletetoken"`
}

type PasteSpec struct {
	IV          string
	Salt        string
	Iterations  int
	KeySize     int
	TagSize     int
	Algorithm   string
	Mode        string
	Compression string
}

func (ps *PasteSpec) SpecArray() []interface{} {
	return []interface{}{
		ps.IV,
		ps.Salt,
		ps.Iterations,
		ps.KeySize,
		ps.TagSize,
		ps.Algorithm,
		ps.Mode,
		ps.Compression,
	}
}

func NewClient(uri *url.URL, username, password string) *Client {
	return &Client{
		URL:      *uri,
		Username: username,
		Password: password,
	}
}

type PasteContent struct {
	Paste          string `json:"paste"`
	Attachment     string `json:"attachment,omitempty"`
	AttachmentName string `json:"attachment_name,omitempty"`
}

type PasteMessage struct {
	Attachment bool
	Filename   string
	Data       []byte
}

func (c *Client) CreatePaste(
	message *PasteMessage,
	expire, formatter string,
	openDiscussion, burnAfterReading, gzip bool,
	password string,
) (*CreatePasteResponse, error) {
	masterKey, err := generateRandomBytes(32)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot generate random bytes: %w", err)
	}

	pasteContent := PasteContent{}
	if message.Attachment {
		pasteContent.AttachmentName = message.Filename
		ext := filepath.Ext(message.Filename)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		pasteContent.AttachmentName = message.Filename
		if pasteContent.AttachmentName == "" {
			pasteContent.AttachmentName = "stdin"
		}
		pasteContent.Attachment = fmt.Sprintf(
			"data:%s;base64,%s",
			mimeType,
			base64.StdEncoding.EncodeToString(message.Data),
		)
	} else {
		pasteContent.Paste = string(message.Data)
	}

	pasteContentStr, err := json.Marshal(&pasteContent)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot marshal paste content: %w", err)
	}

	pasteData, err :=
		encrypt(masterKey, password, pasteContentStr, formatter,
			openDiscussion, burnAfterReading, gzip)
	if err != nil {
		return nil, fmt.Errorf("cannot encrypt data: %w", err)
	}

	createPasteReq := &CreatePasteRequest{
		V:     PrivateBinAPIVersion,
		AData: pasteData.adata(),
		Meta:  CreatePasteRequestMeta{Expire: expire},
		CT: base64.RawStdEncoding.
			EncodeToString(pasteData.Data),
	}

	body, err := json.Marshal(createPasteReq)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot marshal paste request: %w", err)
	}

	req, err := http.NewRequest("POST",
		c.URL.String(),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("cannot create http request: %w", err)
	}

	req.Header.Set("User-Agent", "privatebin-cli/"+pv.Version+" (source; https://github.com/gearnode/privatebin)")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	req.Header.Set("X-Requested-With", "JSONHttpRequest")
	req.SetBasicAuth(c.Username, c.Password)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot execute http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"pastebin server responds with %q status code",
			res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot read response body: %w", err)
	}

	pasteResponse := CreatePasteResponse{}
	err = json.Unmarshal(resBody, &pasteResponse)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot unmarshal response: %w", err)
	}

	if pasteResponse.Status != 0 {
		return nil, fmt.Errorf(
			"status of the paste is not zero: %s",
			pasteResponse.Message)
	}

	pasteId, err := url.Parse(pasteResponse.URL)
	if err != nil {
		return nil, fmt.Errorf("cannot parse paste url: %w", err)
	}

	var uri url.URL
	uri.Scheme = c.URL.Scheme
	uri.Host = c.URL.Host
	uri.Path = c.URL.Path
	uri.RawQuery = pasteId.RawQuery
	uri.Fragment = base58.Encode(masterKey)

	pasteResponse.URL = uri.String()

	return &pasteResponse, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

type PasteData struct {
	*PasteSpec
	Data             []byte
	Formatter        string
	OpenDiscussion   bool
	BurnAfterReading bool
}

func (p *PasteData) adata() []interface{} {
	var b2i = map[bool]int8{false: 0, true: 1}

	return []interface{}{
		p.SpecArray(),
		p.Formatter,
		b2i[p.OpenDiscussion],
		b2i[p.BurnAfterReading],
	}
}

func encrypt(
	masterKey []byte,
	password string,
	message []byte,
	formatter string,
	openDiscussion, burnAfterReading, gzipCompress bool,
) (*PasteData, error) {
	iv, err := generateRandomBytes(12)
	if err != nil {
		return nil, err
	}

	salt, err := generateRandomBytes(8)
	if err != nil {
		return nil, err
	}

	pasteSpec := &PasteSpec{
		IV:          base64.RawStdEncoding.EncodeToString(iv),
		Salt:        base64.RawStdEncoding.EncodeToString(salt),
		Iterations:  310_000,
		KeySize:     256,
		TagSize:     128,
		Algorithm:   "aes",
		Mode:        "gcm",
		Compression: "none",
	}

	if gzipCompress {
		pasteSpec.Compression = "zlib"

		var buf bytes.Buffer
		fw, err := flate.NewWriter(&buf, flate.BestCompression)
		if err != nil {
			return nil, err
		}

		if _, err := fw.Write(message); err != nil {
			return nil, err
		}

		if err := fw.Close(); err != nil {
			return nil, err
		}

		message = buf.Bytes()
	}

	paste := &PasteData{
		Formatter:        formatter,
		OpenDiscussion:   openDiscussion,
		BurnAfterReading: burnAfterReading,
		PasteSpec:        pasteSpec,
	}

	masterKey = append(masterKey, []byte(password)...)
	key := pbkdf2.Key(masterKey,
		salt, paste.Iterations, 32, sha256.New)

	adata, err := json.Marshal(paste.adata())
	if err != nil {
		return nil, err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	data := gcm.Seal(nil, iv, message, adata)

	paste.Data = data

	return paste, nil
}
