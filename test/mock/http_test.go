package mock

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeResponse(t *testing.T) {
	var fakeResponse = FakeResponse{
		BodyProperty:   []byte("fake body"),
		StatusProperty: "200",
		HeaderProperty: map[string][]string{
			"Content-Type": {"application/json"},
		},
		StatusCodeProperty: 200,
		CookiesProperty: []*http.Cookie{
			{
				Name:  "fake-cookie-name",
				Value: "fake-cookie-value",
			},
		},
	}
	assert.Equal(t, fakeResponse.Body(), []byte("fake body"))
	assert.Equal(t, fakeResponse.Status(), "200")
	assert.Equal(t, fakeResponse.Header(), http.Header{"Content-Type": []string{"application/json"}})
	assert.Equal(t, fakeResponse.StatusCode(), 200)
	assert.Equal(t, fakeResponse.Cookies(), []*http.Cookie{
		{
			Name:  "fake-cookie-name",
			Value: "fake-cookie-value",
		},
	})

}
func TestMockHTTPClient(t *testing.T) {
	var mockHTTPClient = HTTPClient{}

	mockHTTPClient.
		On("Post", "", &map[string]interface{}{}, map[string]string{}).
		Return(&FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte("fake body"),
		}, nil)
	mockHTTPClient.
		On("Get", "", &map[string]interface{}{}, map[string]string{}).
		Return(&FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte("fake body"),
		}, nil)
	mockHTTPClient.
		On("Delete", "", &map[string]interface{}{}, map[string]string{}).
		Return(&FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte("fake body"),
		}, nil)
	mockHTTPClient.
		On("Put", "", &map[string]interface{}{}, map[string]string{}).
		Return(&FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte("fake body"),
		}, nil)
	mockHTTPClient.
		On("Patch", "", &map[string]interface{}{}, map[string]string{}).
		Return(&FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte("fake body"),
		}, nil)
	mockHTTPClient.
		On("Request", "", "", &map[string]interface{}{}, map[string]string{}).
		Return(&FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte("fake body"),
		}, nil)
	res, _ := mockHTTPClient.Post("", &map[string]interface{}{}, map[string]string{})
	assert.Equal(t, res.StatusCode(), 200)
	res, _ = mockHTTPClient.Get("", &map[string]interface{}{}, map[string]string{})
	assert.Equal(t, res.StatusCode(), 200)
	res, _ = mockHTTPClient.Put("", &map[string]interface{}{}, map[string]string{})
	assert.Equal(t, res.StatusCode(), 200)
	res, _ = mockHTTPClient.Patch("", &map[string]interface{}{}, map[string]string{})
	assert.Equal(t, res.StatusCode(), 200)
	res, _ = mockHTTPClient.Delete("", &map[string]interface{}{}, map[string]string{})
	assert.Equal(t, res.StatusCode(), 200)
	res, _ = mockHTTPClient.Request("", "", &map[string]interface{}{}, map[string]string{})
	assert.Equal(t, res.StatusCode(), 200)
}
