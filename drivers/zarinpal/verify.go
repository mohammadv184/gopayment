package zarinpal

import (
	"encoding/json"
	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/receipt"
)

func (d *Driver) Verify(amount string, transID string) (*receipt.Receipt, error) {
	var reqBody = map[string]interface{}{
		"authority":   transID,
		"merchant_id": d.MerchantID,
		"amount":      amount,
	}
	resp, _ := client.Post(ApiVerifyUrl, reqBody, nil)
	if resp.StatusCode() != 100 {
		return nil, e.ErrInvalidPayment{}
	}

	var res map[string]interface{}
	err := json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}
	rec := receipt.NewReceipt(res["data"].(map[string]interface{})["ref_id"].(string), d.GetDriverName())

	return rec, nil
}
