package http

import (
	nethttp "net/http"
)

// Client is an Adapter for resty client
type Client interface {
	Get(url string, query interface{}, header map[string]string) (Response, error)
	Post(url string, body interface{}, header map[string]string) (Response, error)
	Put(url string, body interface{}, header map[string]string) (Response, error)
	Patch(url string, body interface{}, header map[string]string) (Response, error)
	Delete(url string, body interface{}, header map[string]string) (Response, error)
	Request(method string, url string, body interface{}, header map[string]string) (Response, error)
}

// Response is an Adapter for resty response
type Response interface {
	Status() string
	StatusCode() int
	Body() []byte
	Header() nethttp.Header
	Cookies() []*nethttp.Cookie
}
