package http

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
	"sync"
)

var once sync.Once

// Http client
type Http struct {
	client *resty.Client
}

var singleton *Http

// NewHttp returns a new instance of Http
func NewHttp() *Http {
	once.Do(func() {
		singleton = &Http{
			client: resty.New(),
		}
	})
	return singleton
}

// Get sends a GET request
func (h *Http) Get(url string, query interface{}, header map[string]string) (*resty.Response, error) {
	return h.Request("GET", url, query, header)
}

// Post sends a POST request
func (h *Http) Post(url string, body interface{}, header map[string]string) (*resty.Response, error) {
	return h.Request("POST", url, body, header)
}

// Put sends a PUT request
func (h *Http) Put(url string, body interface{}, header map[string]string) (*resty.Response, error) {
	return h.Request("PUT", url, body, header)
}

// Patch sends a PATCH request
func (h *Http) Patch(url string, body interface{}, header map[string]string) (*resty.Response, error) {
	return h.Request("PATCH", url, body, header)
}

// Delete sends a DELETE request
func (h *Http) Delete(url string, body interface{}, header map[string]string) (*resty.Response, error) {
	return h.Request("DELETE", url, body, header)
}

// Request sends a http request
func (h *Http) Request(method string, url string, body interface{}, header map[string]string) (*resty.Response, error) {
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
