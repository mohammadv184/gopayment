package zarinpal

import (
	httpClient "github.com/mohammadv184/gopayment/pkg/http"
)

type Driver struct {
	MerchantID  string
	Callback    string
	Description string
}

const (
	ApiPurchaseUrl = "https://api.zarinpal.com/pg/v4/payment/request.json"
	ApiVerifyUrl   = "https://api.zarinpal.com/pg/v4/payment/verify.json"
	ApiPaymentUrl  = "https://www.zarinpal.com/pg/StartPay/"
)

var client httpClient.Client

func init() {
	client = httpClient.NewHttp()
}
func (d Driver) GetDriverName() string {
	return "ZarinPal"
}
