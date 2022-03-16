package mock

import (
	"net/http"

	httpClient "github.com/mohammadv184/gopayment/pkg/http"
	"github.com/stretchr/testify/mock"
)

// HTTPClient is a mock of http.Client
type HTTPClient struct {
	mock.Mock
}

// Get is a mocking a method
func (m *HTTPClient) Get(url string, query interface{}, header map[string]string) (httpClient.Response, error) {
	args := m.Called(url, query, header)
	return args.Get(0).(httpClient.Response), args.Error(1)

}

// Post is a mocking a method
func (m *HTTPClient) Post(url string, body interface{}, header map[string]string) (httpClient.Response, error) {
	args := m.Called(url, body, header)
	return args.Get(0).(httpClient.Response), args.Error(1)

}

// Put is a mocking a method
func (m *HTTPClient) Put(url string, body interface{}, header map[string]string) (httpClient.Response, error) {
	args := m.Called(url, body, header)
	return args.Get(0).(httpClient.Response), args.Error(1)

}

// Patch is a mocking a method
func (m *HTTPClient) Patch(url string, body interface{}, header map[string]string) (httpClient.Response, error) {
	args := m.Called(url, body, header)
	return args.Get(0).(httpClient.Response), args.Error(1)

}

// Delete is a mocking a method
func (m *HTTPClient) Delete(url string, body interface{}, header map[string]string) (httpClient.Response, error) {
	args := m.Called(url, body, header)
	return args.Get(0).(httpClient.Response), args.Error(1)

}

// Request is a mocking a method
func (m *HTTPClient) Request(method, url string, body interface{}, header map[string]string) (httpClient.Response, error) {
	args := m.Called(method, url, body, header)
	return args.Get(0).(httpClient.Response), args.Error(1)
}

// FakeResponse is a fake http.Response
type FakeResponse struct {
	StatusCodeProperty int
	BodyProperty       []byte
	StatusProperty     string
	HeaderProperty     http.Header
	CookiesProperty    []*http.Cookie
}

// StatusCode returns the StatusCodeProperty
func (m *FakeResponse) StatusCode() int {
	return m.StatusCodeProperty
}

//Body returns the body of the response
func (m *FakeResponse) Body() []byte {
	return m.BodyProperty
}

// Header returns the http.Header
func (m *FakeResponse) Header() http.Header {
	return m.HeaderProperty
}

// Status returns the response status
func (m *FakeResponse) Status() string {
	return m.StatusProperty
}

// Cookies returns cookies
func (m *FakeResponse) Cookies() []*http.Cookie {
	return m.CookiesProperty
}
