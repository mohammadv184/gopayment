package http

import "github.com/go-resty/resty/v2"

// Client is an Adapter for resty client
type Client interface {
	Get(url string, query interface{}, header map[string]string) (*resty.Response, error)
	Post(url string, body interface{}, header map[string]string) (*resty.Response, error)
	Put(url string, body interface{}, header map[string]string) (*resty.Response, error)
	Patch(url string, body interface{}, header map[string]string) (*resty.Response, error)
	Delete(url string, body interface{}, header map[string]string) (*resty.Response, error)
	Request(method string, url string, body interface{}, header map[string]string) (*resty.Response, error)
}
