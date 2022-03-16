package zarinpal

import (
	"encoding/json"

	"github.com/mohammadv184/gopayment/helpers"

	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/invoice"
)

// Purchase sends a request to Zarinpal to purchase an invoice.
func (d *Driver) Purchase(invoice *invoice.Invoice) (string, error) {
	var reqBody = map[string]interface{}{
		"merchant_id":  d.MerchantID,
		"callback_url": d.Callback,
		"description":  invoice.GetDescription(),
		"amount":       invoice.GetAmount(),
		"metadata":     invoice.GetDetails(),
	}
	resp, _ := client.Post(APIPurchaseURL, reqBody, nil)
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

// PayURL returns the url to redirect the user to in order to pay the invoice.
func (*Driver) PayURL(invoice *invoice.Invoice) string {
	return APIPaymentURL + invoice.GetTransactionID()
}

// PayMethod returns the Request Method to be used to pay the invoice.
func (*Driver) PayMethod() string {
	return "GET"
}

// RenderRedirectForm renders the html form for redirect to payment page.
func (d *Driver) RenderRedirectForm(invoice *invoice.Invoice) (string, error) {
	return helpers.RenderRedirectTemplate(d.PayMethod(), d.PayURL(invoice), nil)
}
