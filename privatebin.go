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
	"bytes"
	"compress/flate"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go.gearno.de/encoding/base58"
	"golang.org/x/crypto/pbkdf2"
)

const (
	apiVersion = 2

	iterationCount = 600_000
	keySize        = 256
	tagSize        = 128
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

	ShowPasteOptions struct {
		Password    []byte
		ConfirmBurn bool
	}

	CreatePasteResult struct {
		PasteID     string
		PasteURL    url.URL
		DeleteToken string
	}

	ShowPasteResult struct {
		PasteID      string
		CommentCount int
		Paste        Paste
		Comments     []Comment
	}

	Comment struct {
		CommentID string
		PasteID   string
		ParentID  string
		Nickname  string
		Text      string
	}

	createPasteRequest struct {
		V     int                    `json:"v"`
		AData AData                  `json:"adata"`
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
		Status        int                        `json:"status"`
		Message       string                     `json:"message"`
		ID            string                     `json:"id"`
		URL           string                     `json:"url"`
		V             int                        `json:"v"`
		AData         AData                      `json:"adata"`
		Meta          showPasteRequestMeta       `json:"meta"`
		CT            string                     `json:"ct"`
		Comments      []showPasteResponseComment `json:"comments"`
		CommentCount  int                        `json:"comment_count"`
		CommentOffset int                        `json:"comment_offset"`
		Context       string                     `json:"@context"`
	}

	showPasteResponseCommentMeta struct {
		Icon    string `json:"icon"`
		Created int    `json:"created"`
	}

	showPasteResponseComment struct {
		ID       string                       `json:"id"`
		PasteID  string                       `json:"pasteid"`
		ParentID string                       `json:"parentid"`
		URL      string                       `json:"url"`
		V        int                          `json:"v"`
		CT       string                       `json:"ct"`
		AData    Spec                         `json:"adata"`
		Meta     showPasteResponseCommentMeta `json:"meta"`
	}
)

func WithBasicAuth(username, password string) Option {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

func WithCustomHeaderField(k, v string) Option {
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
		httpClient:             defaultPooledClient(),
		customHTTPHeaderFields: make(map[string]string),
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (c *Client) ShowPaste(
	ctx context.Context,
	urlWithMasterKey url.URL,
	opts ShowPasteOptions,
) (*ShowPasteResult, error) {
	fragment := urlWithMasterKey.Fragment
	if strings.HasPrefix(urlWithMasterKey.Fragment, "-") {
		fragment = urlWithMasterKey.Fragment[1:]

		if !opts.ConfirmBurn {
			return nil, fmt.Errorf("cannot read a paste that is set to be burned after reading")
		}
	}

	masterKey, err := base58.Decode(fragment)
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

	if pasteResponse.Status != 0 {
		return nil, fmt.Errorf("cannot load paste: server respond with %d status: %s", pasteResponse.Status, pasteResponse.Message)
	}

	authData, err := json.Marshal(pasteResponse.AData)
	if err != nil {
		return nil, fmt.Errorf("cannot encode adata: %w", err)
	}

	masterKeyWithPassword := append(masterKey, opts.Password...)
	cipherText, err := decrypt(masterKeyWithPassword, pasteResponse.CT, authData, pasteResponse.AData.Spec)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt data: %w", err)
	}

	var paste Paste
	err = json.Unmarshal(cipherText, &paste)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal paste content: %w", err)
	}

	var comments []Comment
	for i, comment := range pasteResponse.Comments {
		authData, err := json.Marshal(comment.AData)
		if err != nil {
			return nil, fmt.Errorf("cannot encode comment (#%d) adata: %w", i, err)
		}

		data, err := decrypt(masterKeyWithPassword, comment.CT, authData, comment.AData)
		if err != nil {
			return nil, fmt.Errorf("cannot decrypt comment (#%d): %w", i, err)
		}

		var message map[string]string
		err = json.Unmarshal(data, &message)
		if err != nil {
			return nil, fmt.Errorf("cannot decode comment (#%d): %w", i, err)
		}

		comments = append(
			comments,
			Comment{
				CommentID: comment.ID,
				PasteID:   comment.PasteID,
				ParentID:  comment.ParentID,
				Nickname:  message["nickname"],
				Text:      message["comment"],
			},
		)

	}

	return &ShowPasteResult{
		PasteID:      pasteResponse.ID,
		CommentCount: pasteResponse.CommentCount,
		Paste:        paste,
		Comments:     comments,
	}, nil
}

func (c *Client) CreatePaste(
	ctx context.Context,
	data []byte,
	opts CreatePasteOptions,
) (*CreatePasteResult, error) {
	var paste Paste

	if opts.AttachmentName != "" {
		paste = Paste{nil, data, opts.AttachmentName, ""}
	} else {
		paste = Paste{data, nil, "", ""}
	}

	pasteData, err := json.Marshal(&paste)
	if err != nil {
		return nil, fmt.Errorf("cannot json marshal paste content: %w", err)
	}

	masterKey, err := generateRandomBytes(32)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random bytes: %w", err)
	}

	iv, err := generateRandomBytes(12)
	if err != nil {
		return nil, fmt.Errorf("cannot generate iv: %w", err)
	}

	salt, err := generateRandomBytes(8)
	if err != nil {
		return nil, fmt.Errorf("cannot generate salt: %w", err)
	}

	masterKeyWithPassword := append(masterKey, opts.Password...)
	key := pbkdf2.Key(masterKeyWithPassword, salt, iterationCount, keySize/8, sha256.New)

	if opts.Compress == CompressionAlgorithmGZip {
		var buf bytes.Buffer
		fw, err := flate.NewWriter(&buf, flate.BestCompression)
		if err != nil {
			return nil, fmt.Errorf("cannot create new flate writer: %w", err)
		}

		if _, err := fw.Write(pasteData); err != nil {
			return nil, fmt.Errorf("cannot write in flate buf: %w", err)
		}

		if err := fw.Close(); err != nil {
			return nil, fmt.Errorf("cannot close flate writer: %w", err)
		}

		pasteData = buf.Bytes()
	}

	adata := AData{
		Spec{
			iv,
			salt,
			iterationCount,
			keySize,
			tagSize,
			EncryptionAlgorithmAES,
			EncryptionModeGCM,
			opts.Compress,
		},
		opts.Formatter,
		opts.OpenDiscussion,
		opts.BurnAfterReading,
	}

	authData, err := json.Marshal(adata)
	if err != nil {
		return nil, fmt.Errorf("cannot encode adata: %w", err)
	}

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cannot create new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, fmt.Errorf("cannot create new galois counter mode: %w", err)
	}

	cipherText := gcm.Seal(nil, iv, pasteData, authData)

	createPasteReq := &createPasteRequest{
		V:     apiVersion,
		AData: adata,
		Meta:  createPasteRequestMeta{Expire: opts.Expire},
		CT:    base64.StdEncoding.EncodeToString(cipherText),
	}

	var reqBody bytes.Buffer
	err = json.NewEncoder(&reqBody).Encode(createPasteReq)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal paste request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.endpoint.String(),
		&reqBody,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	for k, v := range c.customHTTPHeaderFields {
		req.Header.Set(k, v)
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(reqBody.Len()))
	req.Header.Set("X-Requested-With", "JSONHttpRequest")

	if c.username != "" || c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot execute http request: %w", err)
	}
	defer res.Body.Close()

	pasteResponse := createPasteResponse{}
	err = json.NewDecoder(res.Body).Decode(&pasteResponse)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response body: %w", err)
	}

	if pasteResponse.Status != 0 {
		return nil, fmt.Errorf("cannot create paste: server respond with %d status: %s", pasteResponse.Status, pasteResponse.Message)
	}

	pasteID, err := url.Parse(pasteResponse.URL)
	if err != nil {
		return nil, fmt.Errorf("cannot parse paste url: %w", err)
	}

	fragment := base58.Encode(masterKey)
	if opts.BurnAfterReading {
		fragment = "-" + fragment
	}

	pasteLink := url.URL{
		Scheme:   c.endpoint.Scheme,
		Host:     c.endpoint.Host,
		Path:     c.endpoint.Path,
		RawQuery: pasteID.RawQuery,
		Fragment: fragment,
	}

	return &CreatePasteResult{
		PasteID:     pasteResponse.ID,
		PasteURL:    pasteLink,
		DeleteToken: pasteResponse.DeleteToken,
	}, nil
}

func decrypt(masterKey []byte, ct string, adata []byte, spec Spec) ([]byte, error) {
	encryptedCipherText, err := decode64(ct)
	if err != nil {
		return nil, fmt.Errorf("cannot base64 decode cipher text: %w", err)
	}

	key := pbkdf2.Key(
		masterKey,
		spec.Salt,
		spec.Iterations,
		spec.KeySize/8,
		sha256.New,
	)

	var (
		cipherBlock cipher.Block
		gcm         cipher.AEAD
	)

	switch spec.Algorithm {
	case EncryptionAlgorithmAES:
		cipherBlock, err = aes.NewCipher(key)
		if err != nil {
			return nil, fmt.Errorf("cannot create new cipher: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported encryption algorithm: %q", spec.Algorithm)
	}

	switch spec.Mode {
	case EncryptionModeGCM:
		gcm, err = newGCMWithNonceAndTagSize(
			cipherBlock,
			len(spec.IV),
			spec.TagSize/8,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot create new galois counter mode: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported encryption mode: %q", spec.Mode)
	}

	cipherText, err := gcm.Open(nil, spec.IV, encryptedCipherText, adata)
	if err != nil {
		return nil, err
	}

	switch spec.Compression {
	case CompressionAlgorithmNone:
	case CompressionAlgorithmGZip:
		buf := bytes.NewBuffer(cipherText)
		fr := flate.NewReader(buf)
		defer fr.Close()
		cipherText, err = io.ReadAll(fr)
		if err != nil {
			return nil, fmt.Errorf("cannot read gzip: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported compression mode: %q", spec.Compression)
	}

	return cipherText, nil
}
