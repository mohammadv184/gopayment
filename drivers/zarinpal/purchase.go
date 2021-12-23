package zarinpal

import (
	"encoding/json"
	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/invoice"
)

func (d *Driver) Purchase(invoice *invoice.Invoice) (string, error) {
	var reqBody = map[string]interface{}{
		"merchant_id":  d.MerchantID,
		"callback_url": d.Callback,
		"description":  d.Description,
		"amount":       invoice.Amount,
		"metadata":     invoice.GetDetails(),
	}
	resp, _ := client.Post(ApiPurchaseUrl, reqBody, nil)
	if resp.StatusCode() != 100 {
		return "", e.ErrPurchaseFailed{}
	}
	var res map[string]interface{}
	err := json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return "", err
	}
	return res["data"].(map[string]interface{})["authority"].(string), nil
}
func (d *Driver) PayUrl(invoice *invoice.Invoice) string {
	return ApiPaymentUrl + invoice.TransactionID
}
