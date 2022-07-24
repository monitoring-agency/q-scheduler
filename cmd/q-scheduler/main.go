package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/hellflame/argparse"
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/q-scheduler/models"
	"github.com/myOmikron/q-scheduler/modules/crypto"
	"os"

	"github.com/myOmikron/q-scheduler/server"
)

func main() {
	parser := argparse.NewParser("q-scheduler", "", nil)

	configPath := parser.String("", "config-path", &argparse.Option{
		Default:     "/etc/q-scheduler/config.toml",
		Help:        "Path to the configuration file of q-scheduler.",
		Inheritable: true,
	})

	startParser := parser.AddCommand("start", "Starts the scheduler", &argparse.ParserConfig{
		DisableDefaultShowHelp: true,
	})

	initParser := parser.AddCommand("init", "Initializes the scheduler", &argparse.ParserConfig{
		DisableDefaultShowHelp: true,
	})

	if err := parser.Parse(nil); err != nil {
		fmt.Println(err.Error())
		return
	}

	switch {
	case startParser.Invoked:
		server.StartServer(*configPath)
	case initParser.Invoked:
		config := models.GetConfig(*configPath)

		var privKey *rsa.PrivateKey = nil

		if _, err := os.Stat(config.HTTP.TLSKeyPath); err != nil {
			if os.IsNotExist(err) {
				fmt.Print("Generating private key .. ")
				privKey = crypto.GeneratePrivateKey(config)
				color.Println(color.GREEN, "Done")
			} else {
				color.Println(color.RED, "[Error]")
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		if privKey == nil {
			privKey = crypto.LoadPrivateKey(config)
		}
	}
}
