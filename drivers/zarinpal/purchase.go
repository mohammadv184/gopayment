package zarinpal

import (
	"encoding/json"
	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/invoice"
)

// Purchase sends a request to Zarinpal to purchase an invoice.
func (d *Driver) Purchase(invoice *invoice.Invoice) (string, error) {
	var reqBody = map[string]interface{}{
		"merchant_id":  d.MerchantID,
		"callback_url": d.Callback,
		"description":  d.Description,
		"amount":       invoice.GetAmount(),
		"metadata":     invoice.GetDetails(),
	}
	resp, _ := client.Post(ApiPurchaseUrl, reqBody, nil)
	if resp.StatusCode() != 100 {
		return "", e.ErrPurchaseFailed{
			Message: resp.Status() + " purchase failed",
		}
	}
	var res map[string]interface{}
	err := json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return "", err
	}
	return res["data"].(map[string]interface{})["authority"].(string), nil
}

// PayUrl returns the url to redirect the user to in order to pay the invoice.
func (d *Driver) PayUrl(invoice *invoice.Invoice) string {
	return ApiPaymentUrl + invoice.GetTransactionID()
}

// PayMethod returns the Request Method to be used to pay the invoice.
func (d *Driver) PayMethod() string {
	return "GET"
}
