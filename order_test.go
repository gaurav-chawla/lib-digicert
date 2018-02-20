package lib_digicert

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestClient_Submit(t *testing.T) {

	token := os.Getenv("DIGICERT_API_TOKEN")
	orgId := os.Getenv("DIGICERT_ORG_ID")
	productNameId := os.Getenv("DIGICERT_PRODUCTNAME_ID")

	c := NewDigicertClient(token)

	c.WithDebug()

	oId, err := strconv.ParseInt(orgId, 10, 64)

	csr, csrPem, _, err := GenerateCSRAndKey(pkix.Name{
		CommonName:         "commonName",
		Organization:       []string{"myOrg"},
		Country:            []string{"US"},
		Locality:           []string{"San Fransisco"},
		OrganizationalUnit: []string{"Products"},
		Province:           []string{"CA"},
		StreetAddress:      []string{"123 Street 1"},
	}, []string{"testing.com"})

	if err != nil {
		t.Fatal("failed to generate csr", err)
	}

	order := &SubmitOrderInput{
		Certificate: OrderCertificateInput{
			SignatureHash:     "sha512",
			CommonName:        csr.Subject.CommonName,
			Csr:               *csrPem,
			OrganizationUnits: csr.Subject.OrganizationalUnit,
			Emails:            []string{"test@test.com"},
		},
		Organization: Organization{
			Id: oId,
		},
		ValidityYears: 1,
	}

	switch csr.SignatureAlgorithm {
	case x509.SHA384WithRSA, x509.ECDSAWithSHA384:
		order.Certificate.SignatureHash = "sha384"
	case x509.SHA512WithRSA, x509.ECDSAWithSHA512:
		order.Certificate.SignatureHash = "sha512"
	default:
		order.Certificate.SignatureHash = "sha256"
	}

	orderId, err := c.Submit(order, productNameId)

	if err != nil {
		t.Fatal("failed to submit the order", err)
	}

	//time.Sleep(3 * time.Second)

	certId, status, err := c.View(*orderId)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(*certId + "  " + *status)

	if *status == "issued" {
		b, err := c.Download(*certId)
		if err != nil {
			t.Fatal("failed to download certificate", err)
		}

		fmt.Println(string(b))
	}

	out, err := c.Revoke(*certId, "Revoke it")
	if err != nil {
		t.Fatal("failed to revoke certificate", err)
	}

	fmt.Println(out.Status)
	fmt.Println(out.Id)
}
