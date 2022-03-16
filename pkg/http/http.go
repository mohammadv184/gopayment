package http

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
)

var once sync.Once

// HTTP client
type HTTP struct {
	client *resty.Client
}

var singleton *HTTP

// NewHTTP returns a new instance of HTTP
func NewHTTP() *HTTP {
	once.Do(func() {
		singleton = &HTTP{
			client: resty.New(),
		}
	})
	return singleton
}

// Get sends a GET request
func (h *HTTP) Get(url string, query interface{}, header map[string]string) (Response, error) {
	return h.Request("GET", url, query, header)
}

// Post sends a POST request
func (h *HTTP) Post(url string, body interface{}, header map[string]string) (Response, error) {
	return h.Request("POST", url, body, header)
}

// Put sends a PUT request
func (h *HTTP) Put(url string, body interface{}, header map[string]string) (Response, error) {
	return h.Request("PUT", url, body, header)
}

// Patch sends a PATCH request
func (h *HTTP) Patch(url string, body interface{}, header map[string]string) (Response, error) {
	return h.Request("PATCH", url, body, header)
}

// Delete sends a DELETE request
func (h *HTTP) Delete(url string, body interface{}, header map[string]string) (Response, error) {
	return h.Request("DELETE", url, body, header)
}

// Request sends a http request
func (h *HTTP) Request(method, url string, body interface{}, header map[string]string) (Response, error) {
	method = strings.ToUpper(method)
	req := h.client.R()
	req.SetHeaders(header)
	if method == "GET" {
		for k, v := range body.(map[string]interface{}) {
			req.SetQueryParam(k, fmt.Sprintf("%v", v))
		}
	} else {
		req.SetBody(body)
	}
	resp, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
