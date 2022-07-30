package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/myOmikron/echotools/color"

	"github.com/myOmikron/q-scheduler/models"
)

func GenerateCSR(privKey *rsa.PrivateKey, config *models.Config) bytes.Buffer {
	subj := pkix.Name{
		CommonName:         config.HTTP.DNS,
		Country:            []string{""},
		Province:           []string{""},
		Locality:           []string{""},
		Organization:       []string{"Q"},
		OrganizationalUnit: []string{"Scheduler"},
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, privKey)
	buff := bytes.Buffer{}

	if err := pem.Encode(&buff, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}); err != nil {
		color.Println(color.RED, "[Encoding Error]")
		fmt.Println("Error while encoding csr:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return buff
}
