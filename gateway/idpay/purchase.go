package idpay

import (
	"encoding/json"
	"strconv"

	"github.com/mohammadv184/gopayment/helpers"

	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/invoice"
)

// Purchase send purchase request to idpay gateway
func (d *Driver) Purchase(invoice *invoice.Invoice) (string, error) {
	var reqBody = map[string]interface{}{
		"callback": d.Callback,
		"desc":     invoice.GetDescription(),
		"amount":   invoice.GetAmount(),
		"order_id": invoice.GetUUID(),
	}
	if d := invoice.GetDetail("phone"); d != "" {
		reqBody["phone"] = d
	}
	if d := invoice.GetDetail("email"); d != "" {
		reqBody["mail"] = d
	}
	if d := invoice.GetDetail("name"); d != "" {
		reqBody["name"] = d
	}
	resp, err := client.Post(APIPurchaseURL, reqBody, map[string]string{
		"X-API-KEY": d.MerchantID,
		"X-SANDBOX": strconv.FormatBool(d.Sandbox),
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 201 {
		return "", e.ErrPurchaseFailed{
			Message: resp.Status(),
		}
	}
	var res map[string]interface{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return "", err
	}
	return res["id"].(string), nil
}

// PayURL return pay url
func (d *Driver) PayURL(invoice *invoice.Invoice) string {
	if d.Sandbox {
		return APISandBoxPaymentURL + invoice.GetTransactionID()
	}
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
