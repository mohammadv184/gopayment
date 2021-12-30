package idpay

import httpClient "github.com/mohammadv184/gopayment/pkg/http"

type Driver struct {
	MerchantID  string
	Callback    string
	Sandbox     bool
	Description string
}

// Const's for idpay
const (
	ApiPurchaseUrl       = "https://api.idpay.ir/v1.1/payment"
	ApiPaymentUrl        = "https://idpay.ir/p/ws/"
	ApiSandBoxPaymentUrl = "https://idpay.ir/p/ws-sandbox/"
	ApiVerifyUrl         = "https://api.idpay.ir/v1.1/payment/verify"
)

var client httpClient.Client

func init() {
	client = httpClient.NewHttp()
}

func (d Driver) GetDriverName() string {
	return "IDPay"
}
