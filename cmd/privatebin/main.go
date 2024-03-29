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

package main // import "gearno.de/cmd/privatebin"

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"gearno.de/privatebin"
	pv "gearno.de/privatebin/internal/version"
)

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}

func main() {
	ctx := context.Background()

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

		*cfgPath = path.Join(homeDir, ".config", "privatebin", "config.json")
	}

	cfg, err := loadCfgFile(*cfgPath)
	if err != nil {
		fail("cannot load configuration: %v", err)
	}

	binCfg, err := findBinCfg(cfg, *binName)
	if err != nil {
		fail("%v", err)
	}

	uri, err := url.Parse(binCfg.Host)
	if err != nil {
		fail("cannot parse %q bin %q host: %v", binCfg.Name, binCfg.Host, err)
	}

	client := privatebin.NewClient(
		*uri,
		privatebin.WithBasicAuth(
			binCfg.Auth.Username,
			binCfg.Auth.Password,
		),
	)

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

	var attachementName string
	var data []byte
	if *filename != "" {
		file, err := os.Open(*filename)
		if err != nil {
			fail("cannot open %q file: %v", *filename, err)
		}

		data, err = io.ReadAll(file)
		if err != nil {
			fail("cannot read %q file: %v", *filename, err)
		}

		if *attachment {
			attachementName = filepath.Base(*filename)
		}
	} else {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			fail("cannot read stdin: %v", err)
		}

		if *attachment {
			attachementName = "stdin"
		}
	}

	options := privatebin.CreatePasteOptions{
		AttachmentName:   attachementName,
		Formatter:        binCfg.Formatter,
		Expire:           binCfg.Expire,
		OpenDiscussion:   *binCfg.OpenDiscussion,
		BurnAfterReading: *binCfg.BurnAfterReading,
		Password:         []byte(*password),
		Compress:         privatebin.CompressionAlgoNone,
	}

	if *binCfg.GZip {
		options.Compress = privatebin.CompressionAlgoGZip
	}

	resp, err := client.CreatePaste(ctx, data, options)
	if err != nil {
		fail("cannot create the paste: %v", err)
	}

	fmt.Printf("%s\n", resp)
}
