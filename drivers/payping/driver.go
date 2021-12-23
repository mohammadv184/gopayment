package payping

import (
	httpClient "github.com/mohammadv184/gopayment/pkg/http"
)

type Driver struct {
	Token       string
	Callback    string
	Description string
}

const (
	ApiPurchaseUrl = "https://api.payping.ir/v2/pay"
	ApiPaymentUrl  = "https://api.payping.ir/v2/pay/gotoipg/"
	ApiVerifyUrl   = "https://api.payping.ir/v2/pay/verify"
)

var client httpClient.Client

func init() {
	client = httpClient.NewHttp()
}
func (d Driver) GetDriverName() string {
	return "PayPing"
}
