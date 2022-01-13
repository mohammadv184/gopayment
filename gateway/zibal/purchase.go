package zibal

import (
	"encoding/json"
	"strconv"

	e "github.com/mohammadv184/gopayment/errors"

	"github.com/mohammadv184/gopayment/helpers"

	"github.com/mohammadv184/gopayment/invoice"
)

// Purchase sends a request to zibal to purchase an invoice.
func (d *Driver) Purchase(invoice *invoice.Invoice) (string, error) {
	var reqBody = map[string]interface{}{
		"merchant":    d.Merchant,
		"callbackUrl": d.Callback,
		"description": invoice.GetDescription(),
		"orderId":     invoice.GetUUID(),
		"amount":      invoice.GetAmount(),
	}
	if d := invoice.GetDetail("phone"); d != "" {
		reqBody["mobile"] = d
	}
	resp, _ := client.Post(APIPurchaseURL, reqBody, nil)
	var res map[string]interface{}
	err := json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 || res["result"].(float64) != 100 {
		return "", e.ErrPurchaseFailed{
			Message: res["message"].(string),
		}
	}

	return strconv.Itoa(int(res["trackId"].(float64))), nil
}

// PayURL returns the url to redirect the user to in order to pay the invoice.
func (d *Driver) PayURL(invoice *invoice.Invoice) string {
	return APIPaymentURL + invoice.GetTransactionID()
}

// PayMethod returns the Request Method to be used to pay the invoice.
func (d *Driver) PayMethod() string {
	return "GET"
}

// RenderRedirectForm renders the html form for redirect to payment page.
func (d *Driver) RenderRedirectForm(invoice *invoice.Invoice) (string, error) {
	return helpers.RenderRedirectTemplate(d.PayMethod(), d.PayURL(invoice), nil)
}
