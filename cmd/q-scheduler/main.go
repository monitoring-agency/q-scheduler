package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hellflame/argparse"
	"github.com/myOmikron/echotools/color"

	"github.com/monitoring-agency/q-scheduler/models"
	"github.com/monitoring-agency/q-scheduler/modules/crypto"
	"github.com/monitoring-agency/q-scheduler/server"
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
			fmt.Print("Loading private key .. ")
			privKey = crypto.LoadPrivateKey(config)
			color.Println(color.GREEN, "Done")
		}

		fmt.Print("Checking certificate .. ")
		if _, err := os.Stat(config.HTTP.TLSCertPath); err != nil {
			if os.IsNotExist(err) {
				color.Println(color.YELLOW, "Not found")
				fmt.Print("Generating certificate signing request .. ")
				csrBuff := crypto.GenerateCSR(privKey, config)
				color.Println(color.GREEN, "Done")

				if err := ioutil.WriteFile("cert.csr", csrBuff.Bytes(), 0600); err != nil {
					color.Println(color.RED, "[File Error]")
					fmt.Println("Could not write cert.csr:")
					fmt.Println(err.Error())
					os.Exit(1)
				}

				color.Println(color.PURPLE, "Generated \"cert.csr\". Copy the file to the core and execute:")
				color.Println(color.BLUE, "/usr/local/bin/q-core sign /path/to/cert.csr")
			} else {
				color.Println(color.RED, "[Error]")
				fmt.Println(err.Error())
				os.Exit(1)
			}
		} else {
			color.Println(color.GREEN, "Done")
		}
	}
}
