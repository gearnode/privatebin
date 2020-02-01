package main

import (
	"github.com/gearnode/privatebin-cli/privatebin"
)

func main() {
	client, _ := privatebin.NewClient("https://privatebin.net")
	client.CreatePaste("test")
}
