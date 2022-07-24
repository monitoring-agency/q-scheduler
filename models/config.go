package models

import (
	"errors"
	"fmt"
	"github.com/myOmikron/echotools/color"
	"github.com/pelletier/go-toml"
	"io/fs"
	"io/ioutil"
	"os"
)

type HTTP struct {
	ListenAddress string
	TLSKeyPath    string
	TLSCertPath   string
}

type Database struct {
	Path string
}

type Config struct {
	HTTP     HTTP
	Database Database
}

func (c *Config) checkConfig() error {
	if c.HTTP.ListenAddress == "" {
		return errors.New("parameter ListenAddress in section HTTP must not be empty")
	}

	if c.Database.Path == "" {
		return errors.New("parameter Path in section Database must not be empty")
	}

	return nil
}

func GetConfig(configPath string) *Config {
	config := &Config{}

	if configBytes, err := ioutil.ReadFile(configPath); errors.Is(err, fs.ErrNotExist) {
		color.Printf(color.RED, "Config was not found at %s\n", configPath)
		b, _ := toml.Marshal(config)
		fmt.Print(string(b))
		os.Exit(1)
	} else {
		if err := toml.Unmarshal(configBytes, config); err != nil {
			panic(err)
		}
	}

	// Check for config errors
	if err := config.checkConfig(); err != nil {
		color.Println(color.RED, "[Config Error]")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return config
}
