package gateway

import (
	"github.com/mohammadv184/gopayment/invoice"
	"github.com/mohammadv184/gopayment/receipt"
)

// Driver is the interface that must be implemented by all drivers.
type Driver interface {
	// Purchase sends a purchase request to the driver's gateway.
	Purchase(invoice *invoice.Invoice) (string, error)
	// PayUrl returns the url to redirect the user to for payment.
	PayUrl(invoice *invoice.Invoice) string
	// GetDriverName returns the name of the driver.
	GetDriverName() string
	// Verify checks the payment status of the invoice.
	Verify(interface{}) (*receipt.Receipt, error)
	// PayMethod returns the payment request method.
	PayMethod() string
}
