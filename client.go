package lib_digicert

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

const digiCertBaseURL = "https://www.digicert.com/services/v2"

type DigicertClient struct {
	connector *http.Client
	token     string
	baseUrl   string
	logger    func(format string, args ...interface{})
	debug     bool
}

type Errors struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewDigicertClient(token string) *DigicertClient {
	tr := &http.Transport{
		DisableCompression: false,
		DisableKeepAlives:  false,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			InsecureSkipVerify: false,
		},
		Proxy: http.ProxyFromEnvironment,
	}

	c := &DigicertClient{
		connector: &http.Client{Transport: tr},
		token:     token,
		baseUrl:   digiCertBaseURL,
	}

	c.logger = func(format string, args ...interface{}) {
		if c.debug {
			log.Printf(format+"\n", args...)
		}
	}

	return c
}

func (c *DigicertClient) WithBaseURL(baseUrl string) *DigicertClient {
	c.baseUrl = baseUrl
	return c
}

func (c *DigicertClient) WithDebug() *DigicertClient {
	c.debug = true
	return c
}

func (c *DigicertClient) WithoutDebug() *DigicertClient {
	c.debug = false
	return c
}

func (c *DigicertClient) WithLogger(logFn func(format string, args ...interface{})) *DigicertClient {
	c.debug = true
	c.logger = logFn
	return c
}

func (c DigicertClient) call(r *http.Request) (*http.Response, error) {
	r.Header.Set("X-DC-DEVKEY", c.token)

	if c.debug {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			return nil, fmt.Errorf("failed to dump the request payload: %s", err.Error())
		}
		c.logger("api: %s payload: %s", r.URL, dump)
	}

	resp, err := c.connector.Do(r)

	if err != nil {
		return nil, fmt.Errorf("failed to request the digicert api with error: %s", err.Error())
	}

	if resp == nil {
		return nil, fmt.Errorf("received empty response from digicert api")
	}

	if c.debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("failed to dump the response payload: %s", err.Error())
		}
		c.logger("api: %s response: %s status: %d", r.URL, dump, resp.StatusCode)
	}

	if resp.StatusCode >= 300 {
		if resp.Body != nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to request api: %s with err: %s", r.URL, err.Error())
			}
			var errs Errors
			err = json.Unmarshal(body, &errs)
			if err != nil {
				return nil, fmt.Errorf("failed to parse error %q: %s", body, err.Error())
			}
			return nil, fmt.Errorf("request failed. Status: %s, Error: %s", resp.Status, errs)
		}
		return nil, fmt.Errorf("request failed. Status: %s", resp.Status)
	}
	return resp, nil
}

func (c *DigicertClient) request(requestPayload interface{}, apiPath, method string, response interface{}) (interface{}, error) {

	var req []byte
	var err error
	if requestPayload != nil {
		req, err = json.Marshal(requestPayload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload. Err %s", err.Error())
		}
	}

	c.logger("request payload", req)

	r, err := http.NewRequest(method, c.baseUrl+apiPath, bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("request failed. Error: %s", err.Error())
	}

	r.Header.Set("Accept", "application/json")
	if (r.Method == http.MethodPost || r.Method == http.MethodPut) && r.Body != nil {
		r.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.call(r)
	if err != nil {
		return nil, fmt.Errorf("request failed. Error: %s", err.Error())
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	defer resp.Body.Close()

	err = json.Unmarshal(buf.Bytes(), response)
	if err != nil {
		c.logger("unable to unmarshal struct")
		return nil, err
	}

	return response, nil
}

func (c *DigicertClient) simpleRequest(apiPath, method string, headers map[string]string) ([]byte, error) {

	r, err := http.NewRequest(method, c.baseUrl+apiPath, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed. Error: %s", err.Error())
	}

	for k, v := range headers {
		r.Header.Set(k, v)
	}

	resp, err := c.call(r)
	if err != nil {
		return nil, fmt.Errorf("request failed. Error: %s", err.Error())
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	return buf.Bytes(), nil
}

func stringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func toString(s string) *string {
	return &s
}
