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
	ApiPurchaseUrl = "https://api.payping.ir/v2/pay"
	ApiPaymentUrl  = "https://api.payping.ir/v2/pay/gotoipg/"
	ApiVerifyUrl   = "https://api.payping.ir/v2/pay/verify"
)

var client httpClient.Client

func init() {
	client = httpClient.NewHttp()
}

// GetDriverName returns the name of the driver
func (d Driver) GetDriverName() string {
	return "PayPing"
}
