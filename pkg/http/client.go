package http

import "github.com/go-resty/resty/v2"

type Client interface {
	Get(url string, query map[string]interface{}, header map[string]string) (*resty.Response, error)
	Post(url string, body map[string]interface{}, header map[string]string) (*resty.Response, error)
	Put(url string, body map[string]interface{}, header map[string]string) (*resty.Response, error)
	Patch(url string, body map[string]interface{}, header map[string]string) (*resty.Response, error)
	Delete(url string, body map[string]interface{}, header map[string]string) (*resty.Response, error)
	Request(method string, url string, body map[string]interface{}, header map[string]string) (*resty.Response, error)
}
