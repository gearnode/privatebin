package privatebin

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"math/rand"
	"net/url"
	"time"
)

type Client struct {
	URL url.URL
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

func NewClient(host string) (*Client, error) {
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return &Client{URL: *url}, nil
}

func (c *Client) CreatePaste(message []byte) (*CreatePasteResponse, error) {
	masterKey, err := generateRandomBytes(32)
	if err != nil {
		return nil, err
	}

	pasteData, err := encrypt(masterKey, message)
	if err != nil {
		return nil, err
	}

	createPasteReq := &CreatePasteRequest{
		V:     2,
		AData: pasteData.adata(),
		Meta:  CreatePasteRequestMeta{Expire: "1week"},
		CT:    base64.RawStdEncoding.EncodeToString(pasteData.Data),
	}

	fmt.Printf("%#v\n", createPasteReq)

	return nil, nil
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
	Data []byte
}

func (p *PasteData) adata() []interface{} {
	return []interface{}{
		p.SpecArray(),
		"plaintext",
		0,
		0,
	}
}

func encrypt(masterKey []byte, message []byte) (*PasteData, error) {
	iv, err := generateRandomBytes(12)
	if err != nil {
		return nil, err
	}

	salt, err := generateRandomBytes(8)
	if err != nil {
		return nil, err
	}

	paste := &PasteData{
		PasteSpec: &PasteSpec{
			IV:          base64.RawStdEncoding.EncodeToString(iv),
			Salt:        base64.RawStdEncoding.EncodeToString(salt),
			Iterations:  100000,
			KeySize:     256,
			TagSize:     128,
			Algorithm:   "aes",
			Mode:        "gcm",
			Compression: "none",
		},
	}

	key := pbkdf2.Key(masterKey, salt, paste.Iterations, 32, sha256.New)

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
