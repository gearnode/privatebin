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

const (
	helpMessage = `Usage of bin/privatebin:
  -bin string
        the privatebin name to use
  -cfg-file string
        the path of the configuration file (default "~/.config/privatebin/config.json")
  -help
        shows this help message
  -version
        prints the privatebin cli version

Commands:
  create [-attachment] [-burn-after-reading] [-expire=<value>] [-filename=<value>] [-formatter=<value>] [-gzip] [-open-discussion] [-password=<value>]
        Create a new paste.

  show [-insecure] [-password=<value>]
        Show a paste content.
`
)

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}

func main() {
	ctx := context.Background()

	help := flag.Bool("help", false, "shows this help message")
	cfgPath := flag.String("cfg-file", "", "the path of the configuration file (default \"~/.config/privatebin/config.json\")")
	binName := flag.String("bin", "", "the privatebin name to use")
	version := flag.Bool("version", false, "prints the privatebin cli version")

	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stdout, helpMessage)
		os.Exit(0)
	}

	if *version {
		fmt.Fprintf(os.Stdout, "privatebin cli version %s\n", pv.Version)
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

	if len(os.Args) == 1 {
		fmt.Fprint(os.Stderr, helpMessage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		handleCreate(ctx, binCfg, client)
	case "show":
		handleShow(ctx, binCfg, client)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		fmt.Fprintf(os.Stderr, "Use 'privatebin -h' for a list of commands.\n")
		os.Exit(2)
	}
}

func handleCreate(ctx context.Context, binCfg *BinCfg, client *privatebin.Client) {
	createCmd := flag.NewFlagSet("privatebin create", flag.ExitOnError)

	expire := createCmd.String("expire", "", "the time to live of the paste")
	openDiscussion := createCmd.Bool("open-discussion", false, "enable discussion on the paste")
	burnAfterReading := createCmd.Bool("burn-after-reading", false, "delete the paste after reading")
	gzip := createCmd.Bool("gzip", true, "gzip the paste data")
	formatter := createCmd.String("formatter", "", "the text formatter to use, can be plaintext, markdown or syntaxhighlighting")
	password := createCmd.String("password", "", "the paste password")
	filename := createCmd.String("filename", "", "read filepath instead of stdin")
	attachment := createCmd.Bool("attachment", false, "create the paste as an attachment")

	if err := createCmd.Parse(flag.Args()[1:]); err != nil {
		fmt.Println("Failed to parse create command flags:", err)
		os.Exit(3)
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

	if gzip != nil {
		binCfg.GZip = gzip
	}

	if *formatter != "" {
		binCfg.Formatter = *formatter
	}

	var (
		attachementName string
		data            []byte
		err             error
	)
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
		Compress:         privatebin.CompressionAlgorithmNone,
	}

	if *binCfg.GZip {
		options.Compress = privatebin.CompressionAlgorithmGZip
	}

	resp, err := client.CreatePaste(ctx, data, options)
	if err != nil {
		fail("cannot create the paste: %v", err)
	}

	fmt.Printf("%s\n", resp)
}

func handleShow(ctx context.Context, binCfg *BinCfg, client *privatebin.Client) {
	showCmd := flag.NewFlagSet("privatebin show", flag.ExitOnError)

	insecure := showCmd.Bool("insecure", false, "")
	password := showCmd.String("password", "", "the paste password")

	if err := showCmd.Parse(flag.Args()[1:]); err != nil {
		fmt.Println("Failed to parse create command flags:", err)
		os.Exit(3)
	}

	value := showCmd.Arg(0)
	link, err := url.Parse(value)
	if err != nil {
		fail("cannot parse url: %v", err)
	}

	fmt.Printf("XXX %v\n", *insecure)

	resp, err := client.ShowPaste(ctx, *link, []byte(*password))
	if err != nil {
		fail("cannot show the paste: %v", err)
	}

	fmt.Printf("%v\n", resp)
}
