package lib_digicert

import (
	"fmt"
	"net/http"
	"strconv"
)

type SubmitOrderInput struct {
	Certificate                 `json:"certificate,omitempty"`
	Organization                `json:"organization,omitempty"`
	ValidityYears               int    `json:"validity_years,omitempty"`
	CustomExpirationDate        string `json:"custom_expiration_date,omitempty"`
	Comments                    string `json:"comments,omitempty"`
	DisableRenewalNotifications bool   `json:"disable_renewal_notifications,omitempty"`
	RenewalOfOrderID            int    `json:"renewal_of_order_id,omitempty"`
	PaymentMethod               string `json:"payment_method,omitempty"`
	DisableCt                   bool   `json:"disable_ct,omitempty"`
}

type ServerPlatform struct {
	Id int64 `json:"id,omitempty"`
}

type Certificate struct {
	CommonName        string   `json:"common_name,omitempty"`
	Csr               string   `json:"csr,omitempty"`
	OrganizationUnits []string `json:"organization_units,omitempty"`
	ServerPlatform    `json:"server_platform,omitempty"`
	SignatureHash     string `json:"signature_hash,omitempty"`
	ProfileOption     string `json:"profile_option,omitempty"`
}
type Organization struct {
	Id int64 `json:"id,omitempty"`
}

type CertificateOrderResponse struct {
	Id          int64               `json:"id,omitempty"`
	Certificate CertificateResponse `json:"certificate,omitempty"`
	Status      string              `json:"status,omitempty"`
}

type CertificateResponse struct {
	Id int64 `json:"id,omitempty"`
}

type Requests struct {
	Id     int64  `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}

type OrderResponse struct {
	Id       int64      `json:"id,omitempty"`
	Requests []Requests `json:"requests,omitempty"`
}

//Submit submits the order to generate the certificate and returns the orderId
func (c *DigicertClient) Submit(order *SubmitOrderInput, productNameID string) (*OrderResponse, error) {

	res, err := c.request(order, "/order/certificate/"+productNameID, http.MethodPost, &OrderResponse{})

	if err != nil {
		return nil, err
	}

	orderResp := *res.(*OrderResponse)

	return &orderResp, nil
}

//View returns certificateId and order status for an order
func (c *DigicertClient) View(orderId string) (*string, *string, error) {

	res, err := c.request(nil, fmt.Sprintf("/order/certificate/%s", orderId), http.MethodGet, &CertificateOrderResponse{})
	if err != nil {
		return nil, nil, err
	}

	order := *res.(*CertificateOrderResponse)

	i := strconv.FormatInt(order.Certificate.Id, 10)

	return &i, &order.Status, nil
}
