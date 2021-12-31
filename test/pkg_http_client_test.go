package test

import (
	"encoding/json"
	"testing"

	"github.com/mohammadv184/gopayment/pkg/http"
	"github.com/stretchr/testify/suite"
)

type HTTPClientTestSuite struct {
	suite.Suite
	Client http.Client
}

func (s *HTTPClientTestSuite) SetupTest() {
	s.Client = http.NewHTTP()
}
func (s *HTTPClientTestSuite) TestGet() {
	reqBody := map[string]interface{}{
		"foo": "bar",
	}
	headers := map[string]string{
		"Content-Type": "application/json",
		"Foo":          "bar",
	}
	req, err := s.Client.Get("https://httpbin.org/get", reqBody, headers)
	s.Require().Nil(err)
	var respBody map[string]interface{}
	err = json.Unmarshal(req.Body(), &respBody)
	s.Require().Nil(err)

	s.Equal(reqBody, respBody["args"])

}
func (s *HTTPClientTestSuite) TestPost() {
	reqBody := map[string]interface{}{
		"foo": "bar",
	}
	headers := map[string]string{
		"Content-Type": "application/json",
		"Foo":          "bar",
	}
	req, err := s.Client.Post("https://httpbin.org/post", reqBody, headers)
	s.Require().Nil(err)
	var respBody map[string]interface{}
	err = json.Unmarshal(req.Body(), &respBody)
	s.Require().Nil(err)
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(respBody["data"].(string)), &jsonData)
	s.Require().Nil(err)

	s.Equal(reqBody, jsonData)

}
func (s *HTTPClientTestSuite) TestPut() {
	reqBody := map[string]interface{}{
		"foo": "bar",
	}
	headers := map[string]string{
		"Content-Type": "application/json",
		"Foo":          "bar",
	}
	req, err := s.Client.Put("https://httpbin.org/put", reqBody, headers)
	s.Require().Nil(err)
	var respBody map[string]interface{}
	err = json.Unmarshal(req.Body(), &respBody)
	s.Require().Nil(err)
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(respBody["data"].(string)), &jsonData)
	s.Require().Nil(err)

	s.Equal(reqBody, jsonData)

}
func (s *HTTPClientTestSuite) TestPatch() {
	reqBody := map[string]interface{}{
		"foo": "bar",
	}
	headers := map[string]string{
		"Content-Type": "application/json",
		"Foo":          "bar",
	}
	req, err := s.Client.Patch("https://httpbin.org/patch", reqBody, headers)
	s.Require().Nil(err)
	var respBody map[string]interface{}
	err = json.Unmarshal(req.Body(), &respBody)
	s.Require().Nil(err)
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(respBody["data"].(string)), &jsonData)
	s.Require().Nil(err)

	s.Equal(reqBody, jsonData)

}
func (s *HTTPClientTestSuite) TestDelete() {
	reqBody := map[string]interface{}{
		"foo": "bar",
	}
	headers := map[string]string{
		"Content-Type": "application/json",
		"Foo":          "bar",
	}
	req, err := s.Client.Delete("https://httpbin.org/delete", reqBody, headers)
	s.Require().Nil(err)
	var respBody map[string]interface{}
	err = json.Unmarshal(req.Body(), &respBody)
	s.Require().Nil(err)
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(respBody["data"].(string)), &jsonData)
	s.Require().Nil(err)

	s.Equal(reqBody, jsonData)

}
func (s *HTTPClientTestSuite) TestRequest() {
	reqBody := map[string]interface{}{
		"foo": "bar",
	}
	headers := map[string]string{
		"Content-Type": "application/json",
		"Foo":          "bar",
	}
	req, err := s.Client.Request("get", "https://httpbin.org/get", reqBody, headers)
	s.Require().Nil(err)
	var respBody map[string]interface{}
	err = json.Unmarshal(req.Body(), &respBody)
	s.Require().Nil(err)

	s.Equal(reqBody, respBody["args"])

	_, err = s.Client.Request("get", "test", reqBody, headers)
	s.Require().NotNil(err)

}
func TestHttpClientTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPClientTestSuite))
}
