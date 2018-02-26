package lib_digicert

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestClient_Submit(t *testing.T) {

	domain := "testing.com"
	commonName := "cn." + domain

	token := os.Getenv("DIGICERT_API_TOKEN")
	orgId := os.Getenv("DIGICERT_ORG_ID")
	productNameId := os.Getenv("DIGICERT_PRODUCTNAME_ID")

	c := NewDigicertClient(token)

	oId, err := strconv.ParseInt(orgId, 10, 64)

	csr, csrPem, _, err := GenerateCSRAndKey(pkix.Name{
		CommonName:         commonName,
		Organization:       []string{"myOrg"},
		Country:            []string{"US"},
		Locality:           []string{"San Fransisco"},
		OrganizationalUnit: []string{"Products"},
		Province:           []string{"CA"},
		StreetAddress:      []string{"123 Street 1"},
	}, []string{domain})

	if err != nil {
		t.Fatal("failed to generate csr", err)
	}

	today := time.Now()
	customExp1Day := today.Add(time.Hour * 24 * 1)

	order := &SubmitOrderInput{
		Certificate: Certificate{
			SignatureHash:     "sha512",
			CommonName:        csr.Subject.CommonName,
			Csr:               *csrPem,
			OrganizationUnits: csr.Subject.OrganizationalUnit,
		},
		Organization: Organization{
			Id: oId,
		},
		ValidityYears:        1,
		CustomExpirationDate: customExp1Day.Format("2006-01-02"),
	}

	switch csr.SignatureAlgorithm {
	case x509.SHA384WithRSA, x509.ECDSAWithSHA384:
		order.Certificate.SignatureHash = "sha384"
	case x509.SHA512WithRSA, x509.ECDSAWithSHA512:
		order.Certificate.SignatureHash = "sha512"
	default:
		order.Certificate.SignatureHash = "sha256"
	}

	orderResp, err := c.Submit(order, productNameId)

	if err != nil {
		t.Fatal("failed to submit the order", err)
	}

	orderId := orderResp.Id

	//time.Sleep(3 * time.Second)

	issued := false
	var certId *string
	var status *string

	for i := 0; i < 10; i++ {
		certId, status, err = c.View(fmt.Sprint(orderId))

		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(*certId + "  " + *status)
		if *status == "issued" {
			issued = true
			break
		}
		time.Sleep(12 * time.Second)
	}

	if issued {
		b, err := c.Download(*certId)
		if err != nil {
			t.Fatal("failed to download certificate", err)
		}

		fmt.Println(string(b))

		out, err := c.Revoke(*certId, "Revoke it")
		if err != nil {
			t.Fatal("failed to revoke certificate", err)
		}

		fmt.Println(out.Status)
		fmt.Println(out.Id)
	} else {
		fmt.Println("certificate status is not changed to issued. Status: ", *status)
	}
}
