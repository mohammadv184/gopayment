package payping

import (
	"encoding/json"
	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/receipt"
)

func (d *Driver) Verify(amount string, refID string) (*receipt.Receipt, error) {
	var reqBody = map[string]interface{}{
		"refId":  refID,
		"amount": amount,
	}
	resp, _ := client.Post(ApiVerifyUrl, reqBody, map[string]string{
		"Authorization": "Bearer " + d.Token,
	})
	if resp.StatusCode() != 200 {
		return nil, e.ErrInvalidPayment{}
	}

	var res map[string]interface{}
	err := json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}
	rec := receipt.NewReceipt(refID, d.GetDriverName())
	rec.Detail("cardNumber", res["cardNumber"].(string))

	return rec, nil
}
