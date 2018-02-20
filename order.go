package lib_digicert

import (
	"fmt"
	"net/http"
	"strconv"
)

type OrderCertificateInput struct {
	CommonName        string   `json:"common_name"`
	Emails            []string `json:"emails"`
	Csr               string   `json:"csr"`
	OrganizationUnits []string `json:"organization_units,omitempty"`
	SignatureHash     string   `json:"signature_hash"`
}

type Organization struct {
	Id int64 `json:"id"`
}

type SubmitOrderInput struct {
	Certificate   OrderCertificateInput `json:"certificate"`
	Organization  Organization          `json:"organization"`
	ValidityYears int                   `json:"validity_years"`
}

type CertificateOrderResponse struct {
	Id          int64               `json:"id"`
	Certificate CertificateResponse `json:"certificate"`
	Status      string              `json:"status"`
}

type CertificateResponse struct {
	Id int64 `json:"id"`
}

type OrderResponse struct {
	Id int64 `json:"id"`
}

//Submit submits the order to generate the certificate and returns the orderId
func (c *DigicertClient) Submit(order *SubmitOrderInput, productNameID string) (*string, error) {

	res, err := c.Request(order, "/order/certificate/"+productNameID, http.MethodPost, &OrderResponse{})

	if err != nil {
		return nil, err
	}

	orderResp := *res.(*OrderResponse)

	i := strconv.FormatInt(orderResp.Id, 10)

	return &i, nil
}

//View returns certificateId and order status for an order
func (c *DigicertClient) View(orderId string) (*string, *string, error) {

	res, err := c.Request(nil, fmt.Sprintf("/order/certificate/%s", orderId), http.MethodGet, &CertificateOrderResponse{})
	if err != nil {
		return nil, nil, err
	}

	order := *res.(*CertificateOrderResponse)

	i := strconv.FormatInt(order.Certificate.Id, 10)

	return &i, &order.Status, nil
}
