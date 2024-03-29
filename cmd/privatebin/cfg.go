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
		Name             string  `json:"name"`
		Host             string  `json:"host"`
		Auth             AuthCfg `json:"auth"`
		Expire           string  `json:"expire"`
		OpenDiscussion   *bool   `json:"open_discussion"`
		BurnAfterReading *bool   `json:"burn_after_reading"`
		GZip             *bool   `json:"gzip"`
		Formatter        string  `json:"formatter"`
	}

	Cfg struct {
		Bin              []BinCfg `json:"bin"`
		Expire           string   `json:"expire"`
		OpenDiscussion   bool     `json:"open_discussion"`
		BurnAfterReading bool     `json:"burn_after_reading"`
		GZip             bool     `json:"gzip"`
		Formatter        string   `json:"formatter"`
	}
)

func defaultConfig() *Cfg {
	return &Cfg{
		Expire:    "1day",
		Formatter: "plaintext",
		GZip:      true,
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
		return nil, fmt.Errorf("cannot open file: %v", err)
	}

	value, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %v", err)
	}

	cfg := defaultConfig()
	if err := json.Unmarshal(value, cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal file: %v", err)
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

		cfg.Bin[i] = binCfg
	}

	return cfg, nil
}
