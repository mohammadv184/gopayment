package http

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
	"sync"
)

var once sync.Once

type Http struct {
	client *resty.Client
}

var singleton *Http

func NewHttp() *Http {
	once.Do(func() {
		singleton = &Http{
			client: resty.New(),
		}
	})
	return singleton
}
func (h *Http) Get(url string, query map[string]interface{}, header map[string]string) (*resty.Response, error) {
	return h.Request("GET", url, query, header)
}
func (h *Http) Post(url string, body map[string]interface{}, header map[string]string) (*resty.Response, error) {
	return h.Request("POST", url, body, header)
}
func (h *Http) Put(url string, body map[string]interface{}, header map[string]string) (*resty.Response, error) {
	return h.Request("PUT", url, body, header)
}
func (h *Http) Patch(url string, body map[string]interface{}, header map[string]string) (*resty.Response, error) {
	return h.Request("PATCH", url, body, header)
}
func (h *Http) Delete(url string, body map[string]interface{}, header map[string]string) (*resty.Response, error) {
	return h.Request("DELETE", url, body, header)
}
func (h *Http) Request(method string, url string, body map[string]interface{}, header map[string]string) (*resty.Response, error) {
	method = strings.ToUpper(method)
	req := h.client.R()
	req.SetHeaders(header)
	if method == "GET" {
		for k, v := range body {
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
