package payping

import (
	"encoding/json"
	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/invoice"
)

// Purchase send purchase request to payping
func (d *Driver) Purchase(invoice *invoice.Invoice) (string, error) {
	var reqBody = map[string]interface{}{
		"returnUrl":   d.Callback,
		"description": d.Description,
		"amount":      invoice.GetAmount(),
		"clientRefId": invoice.GetUUID(),
	}
	if d, err := invoice.GetDetail("phone"); err == nil {
		reqBody["payerIdentity"] = d
	} else if d, err := invoice.GetDetail("email"); err == nil {
		reqBody["payerIdentity"] = d
	}
	if d, err := invoice.GetDetail("name"); err == nil {
		reqBody["payerName"] = d
	}
	resp, err := client.Post(ApiPurchaseUrl, reqBody, map[string]string{
		"Authorization": "Bearer " + d.Token,
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", e.ErrPurchaseFailed{}
	}
	var res map[string]interface{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return "", err
	}
	return res["code"].(string), nil
}

// PayUrl return pay url
func (d *Driver) PayUrl(invoice *invoice.Invoice) string {
	return ApiPaymentUrl + invoice.GetTransactionID()
}
