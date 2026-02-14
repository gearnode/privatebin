// Copyright (c) 2020-2025 Bryan Frimin <bryan@frimin.fr>.
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

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type (
	AuthCfg struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	BinCfg struct {
		Name              string            `json:"name"`
		Host              string            `json:"host"`
		Auth              AuthCfg           `json:"auth"`
		Expire            string            `json:"expire"`
		OpenDiscussion    *bool             `json:"open-discussion"`
		BurnAfterReading  *bool             `json:"burn-after-reading"`
		GZip              *bool             `json:"gzip"`
		SkipTLSVerify     *bool             `json:"skip-tls-verify"`
		Formatter         string            `json:"formatter"`
		Proxy             string            `json:"proxy"`
		ExtraHeaderFields map[string]string `json:"extra-header-fields"`
	}

	Cfg struct {
		Bin               []BinCfg          `json:"bin"`
		Expire            string            `json:"expire"`
		OpenDiscussion    bool              `json:"open-discussion"`
		BurnAfterReading  bool              `json:"burn-after-reading"`
		GZip              bool              `json:"gzip"`
		SkipTLSVerify     bool              `json:"skip-tls-verify"`
		Formatter         string            `json:"formatter"`
		Proxy             string            `json:"proxy"`
		ExtraHeaderFields map[string]string `json:"extra-header-fields"`
	}
)

func defaultConfig() *Cfg {
	return &Cfg{
		Expire:            "1day",
		Formatter:         "plaintext",
		GZip:              true,
		ExtraHeaderFields: make(map[string]string),
	}
}

func findBinCfg(cfg *Cfg, name string) (*BinCfg, error) {
	for _, bin := range cfg.Bin {
		if bin.Name == name {
			return &bin, nil
		}
	}

	return nil, fmt.Errorf("cannot find %q bin configuration", name)
}

func loadCfgFile(path string) (*Cfg, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	value, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	cfg := defaultConfig()
	if err := json.Unmarshal(value, cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal file: %w", err)
	}

	for i, binCfg := range cfg.Bin {
		if binCfg.Expire == "" {
			binCfg.Expire = cfg.Expire
		}

		if binCfg.OpenDiscussion == nil {
			binCfg.OpenDiscussion = &cfg.OpenDiscussion
		}

		if binCfg.BurnAfterReading == nil {
			binCfg.BurnAfterReading = &cfg.BurnAfterReading
		}

		if binCfg.Formatter == "" {
			binCfg.Formatter = cfg.Formatter
		}

		if binCfg.GZip == nil {
			binCfg.GZip = &cfg.GZip
		}

		if binCfg.SkipTLSVerify == nil {
			binCfg.SkipTLSVerify = &cfg.SkipTLSVerify
		}

		if binCfg.Proxy == "" {
			binCfg.Proxy = cfg.Proxy
		}

		if binCfg.ExtraHeaderFields == nil {
			binCfg.ExtraHeaderFields = cfg.ExtraHeaderFields
		}

		cfg.Bin[i] = binCfg
	}

	return cfg, nil
}
