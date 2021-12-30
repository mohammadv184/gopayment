package idpay

import (
	"encoding/json"
	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/receipt"
	"strconv"
)

// VerifyRequest is the request struct for verify
type VerifyRequest struct {
	RefID string `json:"order_id"`
	ID    string `json:"id"`
}

// Verify is the function to verify payment
func (d *Driver) Verify(vReq interface{}) (*receipt.Receipt, error) {
	verifyReq := vReq.(*VerifyRequest)
	resp, _ := client.Post(ApiVerifyUrl, verifyReq, map[string]string{
		"X-API-KEY": d.MerchantID,
		"X-SANDBOX": strconv.FormatBool(d.Sandbox),
	})

	var res map[string]interface{}
	_ = json.Unmarshal(resp.Body(), &res)
	if _, ok := res["error_message"]; ok {
		return nil, e.ErrInvalidPayment{
			Message: res["error_message"].(string),
		}
	}
	if res["status"].(float64) != 100 {
		return nil, e.ErrInvalidPayment{
			Message: convertResponseStatusToMessage(res["status"].(float64)),
		}
	}
	rec := receipt.NewReceipt(res["track_id"].(string), d.GetDriverName())
	rec.Detail("cardNumber", res["payment"].(map[string]interface{})["card_no"].(string))
	rec.Detail("HashedCardNumber", res["payment"].(map[string]interface{})["hashed_card_no"].(string))

	return rec, nil
}
func convertResponseStatusToMessage(status float64) string {
	switch status {
	case 1:
		return "پرداخت انجام نشده است"
	case 2:
		return "پرداخت ناموفق بوده است"
	case 3:
		return "خطا رخ داده است"
	case 4:
		return "بلوکه شده"
	case 5:
		return "برگشت به پرداخت کننده"
	case 6:
		return "برگشت خورده سیستمی"
	case 7:
		return "انصراف از پرداخت"
	case 8:
		return "به درگاه پرداخت منتقل شد"
	case 10:
		return "در انتظار تایید پرداخت"
	case 100:
		return "پرداخت تایید شده است"
	case 101:
		return "پرداخت قبلا تایید شده است"
	case 200:
		return "به دریافت کننده واریز شد"
	default:
		return "وضعیت نامشخص"
	}

}
