package payping

import (
	httpClient "github.com/mohammadv184/gopayment/pkg/http"
	"github.com/mohammadv184/gopayment/traits"
)

type Driver struct {
	Token       string
	Callback    string
	Description string
	traits.HasDetail
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
