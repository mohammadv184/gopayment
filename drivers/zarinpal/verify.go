package zarinpal

import (
	"encoding/json"
	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/receipt"
)

// VerifyRequest is the request struct for verify
type VerifyRequest struct {
	Amount    string `json:"Amount"`
	Authority string `json:"Authority"`
}

// Verify is the function to verify a payment
func (d *Driver) Verify(vReq interface{}) (*receipt.Receipt, error) {
	verifyReq := vReq.(*VerifyRequest)
	resp, _ := client.Post(ApiVerifyUrl, verifyReq, nil)
	if resp.StatusCode() != 100 {
		return nil, e.ErrInvalidPayment{
			Message: resp.Status() + " Invalid payment",
		}
	}

	var res map[string]interface{}
	err := json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}
	rec := receipt.NewReceipt(res["data"].(map[string]interface{})["ref_id"].(string), d.GetDriverName())
	return rec, nil
}
