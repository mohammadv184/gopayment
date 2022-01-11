package asanpardakht

import (
	"encoding/json"

	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/receipt"
)

// VerifyRequest is the request struct for verify
type VerifyRequest struct {
	InvoiceID string `json:"invoice_id"`
	ID        string `json:"id"`
}

// Verify is the function to verify payment
func (d *Driver) Verify(vReq interface{}) (*receipt.Receipt, error) {
	verifyReq, ok := vReq.(*VerifyRequest)
	if !ok {
		return nil, e.ErrInternal{
			Message: "vReq is not of type VerifyRequest",
		}
	}
	// step 1: get transaction result
	tranRes, err := d.getTranResult(verifyReq)
	if err != nil {
		return nil, err
	}
	payGateTranID, ok := tranRes["PayGateTranID"].(string)
	if !ok {
		return nil, e.ErrInternal{
			Message: "PayGateTranID is not of type string",
		}
	}
	// step 2: verify payment
	err = d.verifyRequest(payGateTranID)
	if err != nil {
		return nil, err
	}
	// step 3: settlement payment
	err = d.settlementRequest(payGateTranID)
	if err != nil {
		return nil, err
	}

	rec := receipt.NewReceipt(tranRes["rrn"].(string), d.GetDriverName())
	rec.Detail("cardNumber", tranRes["payment"].(map[string]interface{})["cardNumber"].(string))
	rec.Detail("transactionId", tranRes["payment"].(map[string]interface{})["refID"].(string))

	return rec, nil
}
func (d *Driver) getTranResult(vReq *VerifyRequest) (map[string]interface{}, error) {
	resp, err := client.Get(APITranResultURL, map[string]string{
		"merchantConfigurationId": d.MerchantConfigID,
		"localInvoiceId":          vReq.InvoiceID,
	}, map[string]string{
		"usr": d.Username,
		"pwd": d.Password,
	})
	if err != nil {
		return nil, e.ErrInternal{
			Message: err.Error(),
		}
	}

	if resp.StatusCode() != 200 {
		return nil, e.ErrInvalidPayment{
			Message: resp.Status(),
		}
	}

	var res map[string]interface{}
	_ = json.Unmarshal(resp.Body(), &res)
	return res, nil
}
func (d *Driver) verifyRequest(payGateTranID string) error {
	resp, err := client.Post(APIVerifyURL, map[string]string{
		"merchantConfigurationId": d.MerchantConfigID,
		"payGateTranID":           payGateTranID,
	}, map[string]string{
		"usr": d.Username,
		"pwd": d.Password,
	})

	if err != nil {
		return e.ErrInternal{
			Message: err.Error(),
		}
	}

	if resp.StatusCode() != 200 {
		return e.ErrInvalidPayment{
			Message: resp.Status(),
		}
	}
	return nil

}
func (d *Driver) settlementRequest(payGateTranID string) error {
	resp, err := client.Post(APISettlementURL, map[string]string{
		"merchantConfigurationId": d.MerchantConfigID,
		"payGateTranID":           payGateTranID,
	}, map[string]string{
		"usr": d.Username,
		"pwd": d.Password,
	})

	if err != nil {
		return e.ErrInternal{
			Message: err.Error(),
		}
	}

	if resp.StatusCode() != 200 {
		return e.ErrInvalidPayment{
			Message: resp.Status(),
		}
	}
	return nil
}
