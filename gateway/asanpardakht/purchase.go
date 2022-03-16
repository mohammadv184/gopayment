package asanpardakht

import (
	"time"

	"github.com/mohammadv184/gopayment/helpers"

	e "github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/invoice"
)

// Purchase send purchase request to asanpardakht gateway
func (d *Driver) Purchase(invoice *invoice.Invoice) (string, error) {
	var reqBody = map[string]interface{}{
		"callbackURL":             d.Callback,
		"additionalData":          invoice.GetDescription(),
		"amountInRials":           invoice.GetAmount() * 10,
		"localInvoiceId":          invoice.GetUUID(),
		"serviceTypeId":           1,
		"localDate":               time.Now().Format("20060102 150405"),
		"merchantConfigurationId": d.MerchantConfigID,
		"paymentId":               0,
	}

	resp, err := client.Post(APIPurchaseURL, reqBody, map[string]string{
		"usr": d.Username,
		"pwd": d.Password,
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", e.ErrPurchaseFailed{
			Message: resp.Status(),
		}
	}
	return string(resp.Body()), nil
}

// PayURL return pay url
func (*Driver) PayURL(_ *invoice.Invoice) string {
	return APIPaymentURL
}

// PayMethod returns the Request Method to be used to pay the invoice.
func (*Driver) PayMethod() string {
	return "POST"
}

// RenderRedirectForm renders the html form for redirect to payment page.
func (d *Driver) RenderRedirectForm(invoice *invoice.Invoice) (string, error) {
	return helpers.RenderRedirectTemplate(d.PayMethod(), d.PayURL(invoice), map[string]string{
		"RefId": invoice.GetUUID(),
	})
}
