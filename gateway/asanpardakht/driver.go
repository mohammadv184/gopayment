package asanpardakht

import (
	httpClient "github.com/mohammadv184/gopayment/pkg/http"
)

// Driver configures the AsanPardakht driver
type Driver struct {
	MerchantConfigID string
	Callback         string
	Username         string
	Password         string
}

// Const's for AsanPardakht
const (
	APIPurchaseURL   = "https://ipgrest.asanpardakht.ir/v1/Token"
	APIPaymentURL    = "https://asan.shaparak.ir"
	APIVerifyURL     = "https://ipgrest.asanpardakht.ir/v1/verify"
	APISettlementURL = "https://ipgrest.asanpardakht.ir/v1/Settlement"
	APITranResultURL = "https://ipgrest.asanpardakht.ir/v1/TranResult"
)

var client httpClient.Client

func init() {
	client = httpClient.NewHTTP()
}

// GetDriverName returns driver name
func (Driver) GetDriverName() string {
	return "AsanPardakht"
}

// SetClient sets the http client
func (Driver) SetClient(c httpClient.Client) {
	client = c
}
