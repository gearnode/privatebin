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

package main // import "gearno.de/cmd/privatebin"

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"gearno.de/privatebin"
	pv "gearno.de/privatebin/internal/version"
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
	GZip             *bool   `json:"gzip"`
	Formatter        string  `json:"formatter"`
}

type Cfg struct {
	Bin              []BinCfg `json:"bin"`
	Expire           string   `json:"expire"`
	OpenDiscussion   bool     `json:"open_discussion"`
	BurnAfterReading bool     `json:"burn_after_reading"`
	GZip             bool     `json:"gzip"`
	Formatter        string   `json:"formatter"`
}

func DefaultCfg() *Cfg {
	return &Cfg{
		Expire:    "1day",
		Formatter: "plaintext",
		GZip:      true,
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

		if binCfg.GZip == nil {
			binCfg.GZip = &cfg.GZip
		}

		cfg.Bin[i] = binCfg
	}

	return cfg, nil
}

func main() {
	cfgPath := flag.String("cfg-file", "", "the path of the configuration file (default \"~/.config/privatebin/config.json\")")
	binName := flag.String("bin", "", "the privatebin name to use")
	expire := flag.String("expire", "", "the time to live of the paste")
	openDiscussion := flag.Bool("open-discussion", false, "enable discussion on the paste")
	burnAfterReading := flag.Bool("burn-after-reading", false, "delete the paste after reading")
	gzip := flag.Bool("gzip", true, "gzip the paste data")
	formatter := flag.String("formatter", "", "the text formatter to use, can be plaintext, markdown or syntaxhighlighting")
	password := flag.String("password", "", "the paste password")
	filename := flag.String("filename", "", "read filepath instead of stdin")
	attachment := flag.Bool("attachment", false, "create the paste as an attachment")
	help := flag.Bool("help", false, "shows this help message")
	version := flag.Bool("version", false, "prints the privatebin cli version")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("privatebin cli version %s\n", pv.Version)
		os.Exit(1)
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

	if expire != nil {
		binCfg.Expire = *expire
	}

	if openDiscussion != nil {
		binCfg.OpenDiscussion = openDiscussion
	}

	if burnAfterReading != nil {
		binCfg.BurnAfterReading = burnAfterReading
	}

	if gzip != nil {
		binCfg.GZip = gzip
	}

	if *formatter != "" {
		binCfg.Formatter = *formatter
	}

	message := privatebin.PasteMessage{Attachment: *attachment}
	if *filename != "" {
		file, err := os.Open(*filename)
		if err != nil {
			fail("cannot open %q file: %v", *filename, err)
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			fail("cannot read %q file: %v", *filename, err)
		}

		message.Filename = filepath.Base(*filename)
		message.Data = data
	} else {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fail("cannot read stdin: %v", err)
		}
		message.Data = data
	}

	resp, err := client.CreatePaste(
		&message,
		binCfg.Expire,
		binCfg.Formatter,
		*binCfg.OpenDiscussion,
		*binCfg.BurnAfterReading,
		*binCfg.GZip,
		*password)

	if err != nil {
		fail("cannot create the paste: %v", err)
	}

	fmt.Printf("%s\n", resp.URL)
}
