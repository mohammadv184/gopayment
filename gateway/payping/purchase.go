package payping

import (
	"encoding/json"

	"github.com/mohammadv184/gopayment/helpers"

	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/invoice"
)

// Purchase send purchase request to payping
func (d *Driver) Purchase(invoice *invoice.Invoice) (string, error) {
	var reqBody = map[string]interface{}{
		"returnUrl":   d.Callback,
		"description": invoice.GetDescription(),
		"amount":      invoice.GetAmount(),
		"clientRefId": invoice.GetUUID(),
	}
	if d := invoice.GetDetail("phone"); d != "" {
		reqBody["payerIdentity"] = d
	} else if d := invoice.GetDetail("email"); d != "" {
		reqBody["payerIdentity"] = d
	}
	if d := invoice.GetDetail("name"); d != "" {
		reqBody["payerName"] = d
	}
	resp, err := client.Post(APIPurchaseURL, reqBody, map[string]string{
		"Authorization": "Bearer " + d.Token,
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", e.ErrPurchaseFailed{
			Message: "Purchase failed",
		}
	}
	var res map[string]interface{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return "", err
	}
	return res["code"].(string), nil
}

// PayURL return pay url
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
