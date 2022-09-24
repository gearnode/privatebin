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

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gearnode/privatebin"
)

type AuthCfg struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BinCfg struct {
	Name             string  `json:"name"`
	Host             string  `json:"host"`
	Auth             AuthCfg `json:"auth"`
	Expire           string  `json:"expire"`
	OpenDiscussion   *bool   `json:"open_discussion"`
	BurnAfterReading *bool   `json:"burn_after_reading"`
	Formatter        string  `json:"formatter"`
	UserAgent        string   `json:"user_agent"`
}

type Cfg struct {
	Bin              []BinCfg `json:"bin"`
	Expire           string   `json:"expire"`
	OpenDiscussion   bool     `json:"open_discussion"`
	BurnAfterReading bool     `json:"burn_after_reading"`
	Formatter        string   `json:"formatter"`
	UserAgent        string   `json:"user_agent"`
}

func DefaultCfg() *Cfg {
	return &Cfg{
		Expire:    "1day",
		Formatter: "plaintext",
		UserAgent: "Privatebin cli +https://github.com/gearnode/privatebin",
	}
}

func (cfg *Cfg) FindBinCfg(name string) (*BinCfg, error) {
	for _, bin := range cfg.Bin {
		if bin.Name == name {
			return &bin, nil
		}
	}

	return nil, fmt.Errorf("cannot find %q bin configuration", name)
}

func fail(format string, args ...interface{}) {
	fmt.Printf("error: "+format+"\n", args...)
	os.Exit(1)
}

func loadCfgFile(path string) (*Cfg, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %v", err)
	}

	value, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %v", err)
	}

	cfg := DefaultCfg()
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

		if binCfg.UserAgent == "" {
			binCfg.UserAgent = cfg.UserAgent
		}

		cfg.Bin[i] = binCfg
	}

	return cfg, nil
}

func main() {
	cfgPath := flag.String("cfg-file", "",
		"the path of the configuration file (default "+
			"\"~/.config/privatebin/config.json\")")
	binName := flag.String("bin", "", "the privatebin name to use")
	expire := flag.String("expire", "",
		"the time to live of the paste")
	openDiscussion := flag.Bool("open-discussion", false,
		"enable discussion on the paste")
	burnAfterReading := flag.Bool("burn-after-reading", false,
		"delete the paste after reading")
	formatter := flag.String("formatter", "",
		"the text formatter to use, can be plaintext, markdown"+
			" or syntaxhighlighting")
	password := flag.String("password", "", "the paste password")
	help := flag.Bool("help", false, "shows this help message")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *cfgPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fail("cannot get user home directory: %v", err)
		}

		*cfgPath = path.Join(homeDir, ".config", "privatebin",
			"config.json")
	}

	cfg, err := loadCfgFile(*cfgPath)
	if err != nil {
		fail("cannot load configuration: %v", err)
	}

	binCfg, err := cfg.FindBinCfg(*binName)
	if err != nil {
		fail("%v", err)
	}

	uri, err := url.Parse(binCfg.Host)
	if err != nil {
		fail("cannot parse %q bin host: %v",
			binCfg.Name, binCfg.Host)
	}

	client := privatebin.NewClient(uri,
		binCfg.Auth.Username,
		binCfg.Auth.Password)

	var data []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fail("cannot read stdin: %v", err)
	}

	if expire != nil {
		binCfg.Expire = *expire
	}

	if openDiscussion != nil {
		binCfg.OpenDiscussion = openDiscussion
	}

	if burnAfterReading != nil {
		binCfg.BurnAfterReading = burnAfterReading
	}

	if *formatter != "" {
		binCfg.Formatter = *formatter
	}

	resp, err := client.CreatePaste(
		strings.Join(data, "\n"),
		binCfg.Expire,
		binCfg.Formatter,
		binCfg.UserAgent,
		*binCfg.OpenDiscussion,
		*binCfg.BurnAfterReading,
		*password)

	if err != nil {
		fail("cannot create the paste: %v", err)
	}

	fmt.Printf("%s\n", resp.URL)
}
