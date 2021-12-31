package payping

import (
	httpClient "github.com/mohammadv184/gopayment/pkg/http"
)

// Driver configures the payping driver
type Driver struct {
	Token       string
	Callback    string
	Description string
}

// Const's for payping
const (
	APIPurchaseURL = "https://api.payping.ir/v2/pay"
	APIPaymentURL  = "https://api.payping.ir/v2/pay/gotoipg/"
	APIVerifyURL   = "https://api.payping.ir/v2/pay/verify"
)

var client httpClient.Client

func init() {
	client = httpClient.NewHTTP()
}

// GetDriverName returns the name of the driver
func (d Driver) GetDriverName() string {
	return "PayPing"
}

// SetClient sets the http client
func (d Driver) SetClient(c httpClient.Client) {
	client = c
}
