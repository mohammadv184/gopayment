package zibal

import (
	"encoding/json"

	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/receipt"
)

// VerifyRequest is the request struct for verify
type VerifyRequest struct {
	TrackID string `json:"track_id"`
}

// Verify is the function to verify a payment
func (d *Driver) Verify(vReq interface{}) (*receipt.Receipt, error) {
	verifyReq, ok := vReq.(*VerifyRequest)
	if !ok {
		return nil, e.ErrInternal{
			Message: "vReq is not of type VerifyRequest",
		}
	}
	resp, _ := client.Post(APIVerifyURL, map[string]string{
		"trackId":  verifyReq.TrackID,
		"merchant": d.Merchant,
	}, nil)

	var res map[string]interface{}
	err := json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 || res["result"].(float64) != 100 {
		return nil, e.ErrPurchaseFailed{
			Message: res["message"].(string),
		}
	}

	rec := receipt.NewReceipt(res["refNumber"].(string), d.GetDriverName())
	rec.Detail("cardNumber", res["cardNumber"].(string))
	return rec, nil
}
