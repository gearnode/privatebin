package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/gearnode/privatebin-cli/privatebin"
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
	OpenDiscussion   bool    `json:"open_discussion"`
	BurnAfterReading bool    `json:"burn_after_reading"`
	Formatter        string  `json:"formatter"`
}

type Cfg struct {
	Bin              []BinCfg `json:"bin"`
	Expire           string   `json:"expire"`
	OpenDiscussion   bool     `json:"open_discussion"`
	BurnAfterReading bool     `json:"burn_after_reading"`
	Formatter        string   `json:"formatter"`
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

	var cfg Cfg
	if err := json.Unmarshal(value, &cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal file: %v", err)
	}

	return &cfg, nil
}

func main() {
	cfgPath := flag.String("cfg-file", "", "the path of the configuration file")
	binName := flag.String("bin", "", "the privatebin name to use")
	help := flag.Bool("help", false, "Shows this help message")

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

		*cfgPath = path.Join(homeDir, "privatebin.json")
	}

	cfg, err := loadCfgFile(*cfgPath)
	if err != nil {
		fail("cannot load configuration: %v", err)
	}

	_, err = cfg.FindBinCfg(*binName)
	if err != nil {
		fail("%v", err)
	}

	var uri url.URL

	client, _ := privatebin.NewClient(uri.String())
	resp, _ := client.CreatePaste("test")
	fmt.Printf("%s\n", resp.URL)
}
