package payping

import (
	"encoding/json"

	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/receipt"
)

// VerifyRequest is the request struct for verify
type VerifyRequest struct {
	RefID  string `json:"refId"`
	Amount string `json:"amount"`
}

// Verify is the function to verify payment
func (d *Driver) Verify(vReq interface{}) (*receipt.Receipt, error) {
	verifyReq, ok := vReq.(*VerifyRequest)
	if ok {
		return nil, e.ErrInternal{
			Message: "vReq is not of type VerifyRequest",
		}
	}
	resp, _ := client.Post(APIVerifyURL, verifyReq, map[string]string{
		"Authorization": "Bearer " + d.Token,
	})
	var res map[string]interface{}
	_ = json.Unmarshal(resp.Body(), &res)
	if resp.StatusCode() != 200 {
		if res == nil {
			return nil, e.ErrInvalidPayment{
				Message: "error in verify payment",
			}
		}

		for _, k := range res {
			return nil, e.ErrInvalidPayment{
				Message: k.(string),
			}
		}
		return nil, e.ErrInvalidPayment{}
	}
	rec := receipt.NewReceipt(verifyReq.RefID, d.GetDriverName())
	rec.Detail("cardNumber", res["cardNumber"].(string))

	return rec, nil
}
