package idpay

import (
	"encoding/json"
	"strconv"

	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/invoice"
)

// Purchase send purchase request to idpay gateway
func (d *Driver) Purchase(invoice *invoice.Invoice) (string, error) {
	var reqBody = map[string]interface{}{
		"callback": d.Callback,
		"desc":     d.Description,
		"amount":   invoice.GetAmount(),
		"order_id": invoice.GetUUID(),
	}
	if d, err := invoice.GetDetail("phone"); err == nil {
		reqBody["phone"] = d
	}
	if d, err := invoice.GetDetail("email"); err == nil {
		reqBody["mail"] = d
	}
	if d, err := invoice.GetDetail("name"); err == nil {
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
			Message: resp.Status() + ": " + convertStatusCodeToString(resp.StatusCode()),
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
func convertStatusCodeToString(statusCode int) string {
	switch statusCode {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 202:
		return "Accepted"
	case 204:
		return "No Content"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 406:
		return "Not Acceptable"
	case 409:
		return "Conflict"
	case 410:
		return "Gone"
	case 422:
		return "Unprocessable Entity"
	case 500:
		return "Internal Server Error"
	case 501:
		return "Not Implemented"
	case 503:
		return "Service Unavailable"
	default:
		return "Unknown"
	}
}
