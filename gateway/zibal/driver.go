package zibal

import (
	httpClient "github.com/mohammadv184/gopayment/pkg/http"
)

// Driver config struct for zibal driver
type Driver struct {
	Merchant string `json:"merchant"`
	Callback string `json:"callback"`
}

// Const's for zibal driver
const (
	APIPurchaseURL = "https://gateway.zibal.ir/v1/request"
	APIVerifyURL   = "https://gateway.zibal.ir/v1/verify"
	APIPaymentURL  = "https://gateway.zibal.ir/start/"
)

var client httpClient.Client

func init() {
	client = httpClient.NewHTTP()
}

// GetDriverName returns driver name
func (Driver) GetDriverName() string {
	return "Zibal"
}

// SetClient sets the http client
func (Driver) SetClient(c httpClient.Client) {
	client = c
}
