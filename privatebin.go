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

package privatebin // import "gearno.de/privatebin"

import (
	"bytes"
	"compress/flate"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"gearno.de/base58"
	pv "gearno.de/privatebin/internal/version"
	"golang.org/x/crypto/pbkdf2"
)

const (
	apiVersion = 2

	iterations = 310_000
	keySize    = 256
	tagSize    = 128
	algorithm  = "aes"
	mode       = "gcm"
)

var (
	userAgent = "privatebin-cli/" + pv.Version + " (source; https://github.com/gearnode/privatebin)"
)

type (
	Client struct {
		endpoint               url.URL
		httpClient             *http.Client
		username               string
		password               string
		customHTTPHeaderFields map[string]string
		userAgent              string
	}

	Option func(c *Client)

	CreatePasteOptions struct {
		AttachmentName   string
		Formatter        string
		Expire           string
		OpenDiscussion   bool
		BurnAfterReading bool
		Compress         CompressionAlgorithm
		Password         []byte
	}

	createPasteRequest struct {
		V     int                    `json:"v"`
		AData [4]any                 `json:"adata"`
		Meta  createPasteRequestMeta `json:"meta"`
		CT    string                 `json:"ct"`
	}

	createPasteRequestMeta struct {
		Expire string `json:"expire"`
	}

	createPasteResponse struct {
		ID          string `json:"id"`
		Status      int    `json:"status"`
		Message     string `json:"message"`
		URL         string `json:"url"`
		DeleteToken string `json:"deletetoken"`
	}

	showPasteRequestMeta struct {
		Created    int `json:"created"`
		TimeToLive int `json:"time_to_live"`
	}

	showPasteResponse struct {
		Status        int                  `json:"status"`
		ID            string               `json:"id"`
		URL           string               `json:"url"`
		V             int                  `json:"v"`
		AData         AData                `json:"adata"`
		Meta          showPasteRequestMeta `json:"meta"`
		CT            string               `json:"ct"`
		Comments      []any                `json:"comments"`
		CommentCount  int                  `json:"comment_count"`
		CommentOffset int                  `json:"comment_offset"`
		Context       string               `json:"@context"`
	}
)

func WithBasicAuth(username, password string) Option {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

func WithCustomerHeaderField(k, v string) Option {
	return func(c *Client) {
		c.customHTTPHeaderFields[k] = v
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

func NewClient(endpoint url.URL, options ...Option) *Client {
	client := &Client{
		endpoint:               endpoint,
		httpClient:             http.DefaultClient,
		customHTTPHeaderFields: make(map[string]string),
		userAgent:              userAgent,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (c *Client) ShowPaste(
	ctx context.Context,
	urlWithMasterKey url.URL,
	password []byte,
) (any, error) {
	masterKey, err := base58.Decode(urlWithMasterKey.Fragment)
	if err != nil {
		return nil, fmt.Errorf("cannot decode master key: %w", err)
	}

	urlWithoutMasterKey := urlWithMasterKey
	urlWithoutMasterKey.Fragment = ""

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		urlWithoutMasterKey.String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("X-Requested-With", "JSONHttpRequest")

	if c.username != "" || c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot execute http request: %w", err)
	}
	defer res.Body.Close()

	var pasteResponse showPasteResponse

	err = json.NewDecoder(res.Body).Decode(&pasteResponse)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response body: %w", err)
	}

	masterKeyWithPassword := append(masterKey, password...)

	encryptedCipherText, err := base64.RawStdEncoding.DecodeString(pasteResponse.CT)
	if err != nil {
		return nil, fmt.Errorf("cannot base64 decode cipher text: %w", err)
	}

	authData, err := json.Marshal(pasteResponse.AData)
	if err != nil {
		return "", fmt.Errorf("cannot encode adata: %w", err)
	}

	key := pbkdf2.Key(
		masterKeyWithPassword,
		pasteResponse.AData.Spec.Salt,
		pasteResponse.AData.Spec.Iterations,
		pasteResponse.AData.Spec.KeySize/8,
		sha256.New,
	)

	if pasteResponse.AData.Spec.Algorithm != EncryptionAlgorithmAES {
		return nil, fmt.Errorf("unsupported encryption algorithm: %q", pasteResponse.AData.Spec.Algorithm)
	}

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("cannot create new cipher: %w", err)
	}

	if pasteResponse.AData.Spec.Mode != EncryptionModeGCM {
		return nil, fmt.Errorf("unsupported encryption mode: %q", pasteResponse.AData.Spec.Mode)
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", fmt.Errorf("cannot create new galois counter mode: %w", err)
	}

	cipherText, err := gcm.Open(nil, pasteResponse.AData.Spec.IV, encryptedCipherText, authData)
	if err != nil {
		return nil, err
	}

	if pasteResponse.AData.Spec.Compression == CompressionAlgorithmGZip {
		buf := bytes.NewBuffer(cipherText)
		fr := flate.NewReader(buf)
		defer fr.Close()
		cipherText, err = io.ReadAll(fr)
		if err != nil {
			return nil, fmt.Errorf("cannot read gzip: %w", err)
		}
	}

	paste := map[string]string{}
	err = json.Unmarshal(cipherText, &paste)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal paste content: %w", err)
	}

	return paste, nil
}

func (c *Client) CreatePaste(
	ctx context.Context,
	data []byte,
	opts CreatePasteOptions,
) (string, error) {
	paste := map[string]string{}
	if opts.AttachmentName != "" {
		ext := filepath.Ext(opts.AttachmentName)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		paste["attachment_name"] = "stdin"
		if opts.AttachmentName != "" {
			paste["attachment_name"] = opts.AttachmentName
		}

		paste["attachment"] = fmt.Sprintf(
			"data:%s;base64,%s",
			mimeType,
			base64.StdEncoding.EncodeToString(data),
		)
	} else {
		paste["paste"] = string(data)
	}

	pasteData, err := json.Marshal(&paste)
	if err != nil {
		return "", fmt.Errorf("cannot json marshal paste content: %w", err)
	}

	masterKey, err := generateRandomBytes(32)
	if err != nil {
		return "", fmt.Errorf("cannot generate random bytes: %w", err)
	}

	iv, err := generateRandomBytes(12)
	if err != nil {
		return "", fmt.Errorf("cannot generate iv: %w", err)
	}

	salt, err := generateRandomBytes(8)
	if err != nil {
		return "", fmt.Errorf("cannot generate salt: %w", err)
	}

	masterKeyWithPassword := append(masterKey, opts.Password...)

	key := pbkdf2.Key(masterKeyWithPassword, salt, iterations, keySize/8, sha256.New)

	compression := "none"
	if opts.Compress == CompressionAlgorithmGZip {
		compression = "zlib"

		var buf bytes.Buffer
		fw, err := flate.NewWriter(&buf, flate.BestCompression)
		if err != nil {
			return "", fmt.Errorf("cannot create new flate writer: %w", err)
		}

		if _, err := fw.Write(pasteData); err != nil {
			return "", fmt.Errorf("cannot write in flate buf: %w", err)
		}

		if err := fw.Close(); err != nil {
			return "", fmt.Errorf("cannot close flate writer: %w", err)
		}

		pasteData = buf.Bytes()
	}

	adata := [4]any{
		[8]any{
			base64.RawStdEncoding.EncodeToString(iv),
			base64.RawStdEncoding.EncodeToString(salt),
			iterations,
			keySize,
			tagSize,
			algorithm,
			mode,
			compression,
		},
		opts.Formatter,
		btoi(opts.OpenDiscussion),
		btoi(opts.BurnAfterReading),
	}

	authData, err := json.Marshal(adata)
	if err != nil {
		return "", fmt.Errorf("cannot encode adata: %w", err)
	}

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("cannot create new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", fmt.Errorf("cannot create new galois counter mode: %w", err)
	}

	cipherText := gcm.Seal(nil, iv, pasteData, authData)

	createPasteReq := &createPasteRequest{
		V:     apiVersion,
		AData: adata,
		Meta:  createPasteRequestMeta{Expire: opts.Expire},
		CT:    base64.RawStdEncoding.EncodeToString(cipherText),
	}

	var reqBody bytes.Buffer
	err = json.NewEncoder(&reqBody).Encode(createPasteReq)
	if err != nil {
		return "", fmt.Errorf("cannot marshal paste request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.endpoint.String(),
		&reqBody,
	)
	if err != nil {
		return "", fmt.Errorf("cannot create request: %w", err)
	}

	for k, v := range c.customHTTPHeaderFields {
		req.Header.Set(k, v)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(reqBody.Len()))
	req.Header.Set("X-Requested-With", "JSONHttpRequest")

	if c.username != "" || c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot execute http request: %w", err)
	}
	defer res.Body.Close()

	pasteResponse := createPasteResponse{}
	err = json.NewDecoder(res.Body).Decode(&pasteResponse)
	if err != nil {
		return "", fmt.Errorf("cannot decode response body: %w", err)
	}

	if pasteResponse.Status != 0 {
		return "", fmt.Errorf("status of the paste is not zero: %s", pasteResponse.Message)
	}

	pasteID, err := url.Parse(pasteResponse.URL)
	if err != nil {
		return "", fmt.Errorf("cannot parse paste url: %w", err)
	}

	pasteLink := url.URL{
		Scheme:   c.endpoint.Scheme,
		Host:     c.endpoint.Host,
		Path:     c.endpoint.Path,
		RawQuery: pasteID.RawQuery,
		Fragment: base58.Encode(masterKey),
	}

	return pasteLink.String(), nil
}

func btoi(v bool) int {
	if v {
		return 1
	}

	return 0
}

func itob(v int) bool {
	if v == 0 {
		return false
	}

	return true
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
