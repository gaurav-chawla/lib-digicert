package lib_digicert

import (
	"fmt"
	"net/http"
)

type Requester struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type RevokeCertificateResponse struct {
	Id        int64     `json:"id"`
	Date      string    `json:"date"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Comments  string    `json:"comments,omitempty"`
	Requester Requester `json:"requester"`
}

//Download returns the certificate zip file bytes by certificateId
func (c *DigicertClient) Download(certificateId string) ([]byte, error) {

	headers := map[string]string{
		"Accept": "*/*",
	}

	res, err := c.SimpleRequest(fmt.Sprintf("/certificate/%s/download/format/default", certificateId), http.MethodGet, headers)

	if err != nil {
		return nil, err
	}

	return res, nil
}

//Revoke submits the revoke certificate request
func (c *DigicertClient) Revoke(certificateId, comments string) (*RevokeCertificateResponse, error) {

	res, err := c.Request(struct {
		Comments string `json:"comments,omitempty"`
	}{
		Comments: comments,
	}, fmt.Sprintf("/certificate/%s/revoke", certificateId), http.MethodPut, &RevokeCertificateResponse{})

	if err != nil {
		return nil, err
	}

	return res.(*RevokeCertificateResponse), nil
}
