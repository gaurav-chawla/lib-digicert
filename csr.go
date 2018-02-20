package lib_digicert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"net"
)

//GenerateCSRAndKey generates CST and private key
func GenerateCSRAndKey(subject pkix.Name, hosts []string) (csr *x509.CertificateRequest, csrPEM, privKeyPEM *string, err error) {

	priv, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		err = fmt.Errorf("failed to generate key. Error: %s", err.Error())
		return
	}

	template := x509.CertificateRequest{
		Subject: subject,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	csrB, err := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	if err != nil {
		err = fmt.Errorf("failed to generate CSR. Error: %s", err.Error())
		return
	}

	csr, err = x509.ParseCertificateRequest(csrB)
	if err != nil {
		err = fmt.Errorf("failed to parse CSR. Error: %s", err.Error())
		return
	}

	var csrBuf bytes.Buffer
	err = pem.Encode(&csrBuf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrB})
	if err != nil {
		err = fmt.Errorf("failed to encode pem. Error: %s", err.Error())
		return
	}

	csrPEM = toString(string(csrBuf.Bytes()))
	// convert the private key to PEM
	privKeyPEM, err = getPrivateKeyPEM(priv)
	return csr, csrPEM, privKeyPEM, err
}

func getPrivateKeyPEM(privKey *rsa.PrivateKey) (*string, error) {
	var buf bytes.Buffer

	err := pem.Encode(&buf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})

	if err != nil {
		return nil, err
	}

	return toString(string(buf.Bytes())), nil
}
