package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/myOmikron/echotools/color"

	"github.com/monitoring-agency/q-scheduler/models"
)

func GeneratePrivateKey(config *models.Config) *rsa.PrivateKey {
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		color.Println(color.RED, "[Crypto Error]")
		fmt.Println("Cannot generate RSA key:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var privateKeyBytes = x509.MarshalPKCS1PrivateKey(privKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create(config.HTTP.TLSKeyPath)
	if err != nil {
		color.Println(color.RED, "[File Error]")
		fmt.Printf("Couldn't create key file %s\n", config.HTTP.TLSKeyPath)
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err := pem.Encode(privatePem, privateKeyBlock); err != nil {
		color.Println(color.RED, "[Crypt Error]")
		fmt.Println("Error while encoding private key")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return privKey
}

func LoadPrivateKey(config *models.Config) *rsa.PrivateKey {
	data, err := ioutil.ReadFile(config.HTTP.TLSKeyPath)
	if err != nil {
		color.Println(color.RED, "[File Error]")
		fmt.Printf("Couldn't load private key file %s \n", config.HTTP.TLSKeyPath)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	block, _ := pem.Decode(data)
	if block.Type != "RSA PRIVATE KEY" {
		color.Println(color.RED, "[Crypt Error]")
		fmt.Printf("Wrong type: RSA PRIVATE KEY expected, found: %s\n", block.Type)
		os.Exit(1)
	}
	privPemBytes := block.Bytes

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil {
			color.Println(color.RED, "[Crypt Error]")
			fmt.Println("Couldn't parse PKCS8 key")
			os.Exit(1)
		} else {
			color.Println(color.RED, "[Crypt Error]")
			fmt.Println("Couldn't parse PKCS1 key")
			os.Exit(1)
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	if privateKey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		color.Println(color.RED, "[Cast Error]")
		fmt.Println("Couldn't cast to rsa private key")
		os.Exit(1)
	}

	return privateKey
}
