package gateway

import (
	"github.com/mohammadv184/gopayment/invoice"
	httpClient "github.com/mohammadv184/gopayment/pkg/http"
	"github.com/mohammadv184/gopayment/receipt"
)

// Driver is the interface that must be implemented by all drivers.
type Driver interface {
	// Purchase sends a purchase request to the driver's gateway.
	Purchase(invoice *invoice.Invoice) (string, error)
	// PayURL returns the url to redirect the user to for payment.
	PayURL(invoice *invoice.Invoice) string
	// GetDriverName returns the name of the driver.
	GetDriverName() string
	// Verify checks the payment status of the invoice.
	Verify(vReq interface{}) (*receipt.Receipt, error)
	// PayMethod returns the payment request method.
	PayMethod() string
	// SetClient sets the http client.
	SetClient(client httpClient.Client)
}
