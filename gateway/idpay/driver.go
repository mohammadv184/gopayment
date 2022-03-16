package idpay

import httpClient "github.com/mohammadv184/gopayment/pkg/http"

// Driver configures the idpay driver
type Driver struct {
	MerchantID string
	Callback   string
	Sandbox    bool
}

// Const's for idpay
const (
	APIPurchaseURL       = "https://api.idpay.ir/v1.1/payment"
	APIPaymentURL        = "https://idpay.ir/p/ws/"
	APISandBoxPaymentURL = "https://idpay.ir/p/ws-sandbox/"
	APIVerifyURL         = "https://api.idpay.ir/v1.1/payment/verify"
)

var client httpClient.Client

func init() {
	client = httpClient.NewHTTP()
}

// GetDriverName returns driver name
func (Driver) GetDriverName() string {
	return "IDPay"
}

// SetClient sets the http client
func (Driver) SetClient(c httpClient.Client) {
	client = c
}
