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
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"go.gearno.de/privatebin/v2"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"

	userAgent         = "privatebin-cli/" + version + " (source; https://go.gearno.de/privatebin)"
	cfgPath           string
	binName           string
	extraHeaderFields []string
	client            *privatebin.Client
	binCfg            *BinCfg
	output            string

	ctx           = context.Background()
	clientOptions = []privatebin.Option{
		privatebin.WithUserAgent(userAgent),
	}

	expire           string
	openDiscussion   bool
	burnAfterReading bool
	gzip             bool
	formatter        string
	password         string
	filename         string
	attachment       bool

	insecure      bool
	confirmBurn   bool
	skipTLSVerify bool
	proxy         string

	rootCmd = &cobra.Command{
		Use:     "privatebin",
		Version: fmt.Sprintf("%s-%s (%s)", version, commit, date),
		Short:   "A streamlined CLI for effortlessly creating and managing PrivateBin pastes",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			switch output {
			case "":
			case "json":
			default:
				return fmt.Errorf("invalid output: %q, valid options are '', 'json'", output)
			}

			if cfgPath == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("cannot get user home directory: %w", err)
				}

				cfgPath = path.Join(homeDir, ".config", "privatebin", "config.json")
			}

			cfg, err := loadCfgFile(cfgPath)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("cannot load configuration: %w", err)
				}
				cfg = defaultConfig()
			}

			binCfg, err = findBinCfg(cfg, binName)
			if err != nil {
				binCfg = &BinCfg{
					Expire:            cfg.Expire,
					OpenDiscussion:    &cfg.OpenDiscussion,
					BurnAfterReading:  &cfg.BurnAfterReading,
					GZip:              &cfg.GZip,
					Formatter:         cfg.Formatter,
					SkipTLSVerify:     &cfg.SkipTLSVerify,
					Proxy:             cfg.Proxy,
					ExtraHeaderFields: cfg.ExtraHeaderFields,
				}
			}

			clientOptions = append(
				clientOptions,
				privatebin.WithBasicAuth(
					binCfg.Auth.Username,
					binCfg.Auth.Password,
				),
			)

			for k, v := range binCfg.ExtraHeaderFields {
				clientOptions = append(
					clientOptions,
					privatebin.WithCustomHeaderField(k, v),
				)
			}

			for _, value := range extraHeaderFields {
				parts := strings.SplitN(value, ":", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid header field format: '%s', expected 'key: value'", value)
				}

				clientOptions = append(
					clientOptions,
					privatebin.WithCustomHeaderField(
						strings.TrimSpace(parts[0]),
						strings.TrimSpace(parts[1]),
					),
				)
			}

			if (binCfg.SkipTLSVerify != nil && *binCfg.SkipTLSVerify) || skipTLSVerify {
				tlsConfig := &tls.Config{
					InsecureSkipVerify: true,
				}

				clientOptions = append(
					clientOptions,
					privatebin.WithTLSConfig(tlsConfig),
				)
			}

			proxyAddr := binCfg.Proxy
			if proxy != "" {
				proxyAddr = proxy
			}

			if proxyAddr != "" {
				proxyURL, err := url.Parse(proxyAddr)
				if err != nil {
					return fmt.Errorf("cannot parse proxy url %q: %w", proxyAddr, err)
				}

				clientOptions = append(
					clientOptions,
					privatebin.WithProxyURL(*proxyURL),
				)
			}

			host, err := url.Parse(binCfg.Host)
			if err != nil {
				return fmt.Errorf("cannot parse %q bin %q host: %w", binCfg.Name, binCfg.Host, err)
			}

			client = privatebin.NewClient(*host, clientOptions...)
			return nil
		},
	}

	showCmd = &cobra.Command{
		Use:          "show",
		Short:        "Show a paste",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			link, err := url.Parse(args[0])
			if err != nil {
				return fmt.Errorf("cannot parse paste url: %w", err)
			}

			if link.Scheme+"://"+link.Host != strings.TrimRight(binCfg.Host, "/") {
				if !insecure {
					return fmt.Errorf("untrusted privatebin instance use --insecure flag or add it to the configuration")
				}
			}

			options := privatebin.ShowPasteOptions{
				Password:    []byte(password),
				ConfirmBurn: confirmBurn,
			}

			result, err := client.ShowPaste(ctx, *link, options)
			if err != nil {
				return fmt.Errorf("cannot show paste: %w", err)
			}

			switch output {
			case "":
				fmt.Fprintf(os.Stdout, "%s\n", result.Paste.Data)
			case "json":
				var comments []map[string]string
				for _, comment := range result.Comments {
					comments = append(
						comments,
						map[string]string{
							"comment_id": comment.CommentID,
							"paste_id":   comment.PasteID,
							"parent_id":  comment.ParentID,
							"nickname":   comment.Nickname,
							"text":       comment.Text,
						},
					)
				}

				json.NewEncoder(os.Stdout).Encode(
					map[string]any{
						"paste_id": result.PasteID,
						"paste": map[string]string{
							"attachment_name": result.Paste.AttachmentName,
							"attachment":      base64.StdEncoding.EncodeToString(result.Paste.Attachment),
							"data":            base64.StdEncoding.EncodeToString(result.Paste.Data),
						},
						"comment_count": result.CommentCount,
						"comments":      comments,
					},
				)
			}
			return nil
		},
	}

	createCmd = &cobra.Command{
		Use:          "create",
		Short:        "Create a paste",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if binCfg.Host == "" {
				return fmt.Errorf("no privatebin instance configured, please create a configuration file or use the --config flag")
			}

			if cmd.Flags().Changed("expire") {
				binCfg.Expire = expire
			}

			if cmd.Flags().Changed("open-discussion") {
				binCfg.OpenDiscussion = &openDiscussion
			}

			if cmd.Flags().Changed("burn-after-reading") {
				binCfg.BurnAfterReading = &burnAfterReading
			}

			if cmd.Flags().Changed("gzip") {
				binCfg.GZip = &gzip
			}

			if cmd.Flags().Changed("formatter") {
				binCfg.Formatter = formatter
			}

			var (
				attachementName string
				data            []byte
				err             error
			)

			if cmd.Flags().Changed("filename") {
				file, err := os.Open(filename)
				if err != nil {
					return fmt.Errorf("cannot open %q file: %w", filename, err)
				}

				data, err = io.ReadAll(file)
				if err != nil {
					return fmt.Errorf("cannot read %q file: %w", filename, err)
				}

				if cmd.Flags().Changed("attachment") {
					attachementName = filepath.Base(filename)
				}
			} else {
				data, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("cannot read stdin: %w", err)
				}

				if cmd.Flags().Changed("attachment") {
					attachementName = "stdin"
				}
			}

			options := privatebin.CreatePasteOptions{
				AttachmentName:   attachementName,
				Formatter:        binCfg.Formatter,
				Expire:           binCfg.Expire,
				OpenDiscussion:   *binCfg.OpenDiscussion,
				BurnAfterReading: *binCfg.BurnAfterReading,
				Password:         []byte(password),
				Compress:         privatebin.CompressionAlgorithmNone,
			}

			if *binCfg.GZip {
				options.Compress = privatebin.CompressionAlgorithmGZip
			}

			result, err := client.CreatePaste(ctx, data, options)
			if err != nil {
				return fmt.Errorf("cannot create the paste: %w", err)
			}

			switch output {
			case "":
				fmt.Fprintf(os.Stdout, "%s\n", result.PasteURL.String())
			case "json":
				json.NewEncoder(os.Stdout).Encode(
					map[string]any{
						"paste_id":     result.PasteID,
						"paste_url":    result.PasteURL.String(),
						"delete_token": result.DeleteToken,
					},
				)
			}

			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "the command output format")
	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "the config file (default is $HOME/.config/privatebin/config.json)")
	rootCmd.PersistentFlags().StringVarP(&binName, "bin", "b", "", "the name of the privatebin instance to use (default \"\")")
	rootCmd.PersistentFlags().StringSliceVarP(&extraHeaderFields, "header", "H", []string{}, "extra HTTP header fields to include in the request sent")
	rootCmd.PersistentFlags().StringVar(&proxy, "proxy", "", "proxy URL to use for requests (e.g. socks5://127.0.0.1:9050 for TOR)")

	createCmd.Flags().StringVar(&expire, "expire", "", "the time to live of the paste")
	createCmd.Flags().BoolVar(&openDiscussion, "open-discussion", false, "enable discussion on the paste")
	createCmd.Flags().BoolVar(&burnAfterReading, "burn-after-reading", false, "delete the paste after reading")
	createCmd.Flags().BoolVar(&gzip, "gzip", true, "gzip the paste data")
	createCmd.Flags().StringVar(&formatter, "formatter", "", "the text formatter to use, can be plaintext, markdown or syntaxhighlighting")
	createCmd.Flags().StringVar(&password, "password", "", "the paste password")
	createCmd.Flags().StringVar(&filename, "filename", "", "read filepath instead of stdin")
	createCmd.Flags().BoolVar(&attachment, "attachment", false, "create the paste as an attachment")
	createCmd.Flags().BoolVar(&skipTLSVerify, "skip-tls-verify", false, "skip TLS certificate verification")

	showCmd.Flags().BoolVar(&insecure, "insecure", false, "allow reading paste from untrusted instance")
	showCmd.Flags().BoolVar(&confirmBurn, "confirm-burn", false, "confirm paste opening, it will be deleted immediately afterwards")
	showCmd.Flags().StringVar(&password, "password", "", "the paste password")
	showCmd.Flags().BoolVar(&skipTLSVerify, "skip-tls-verify", false, "skip TLS certificate verification")

	rootCmd.AddCommand(showCmd, createCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
