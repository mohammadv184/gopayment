package drivers

import (
	"github.com/mohammadv184/gopayment/invoice"
	"github.com/mohammadv184/gopayment/receipt"
)

type Driver interface {
	Purchase(invoice *invoice.Invoice) (string, error)
	PayUrl(invoice *invoice.Invoice) string
	GetDriverName() string
	Verify(interface{}) (*receipt.Receipt, error)
}
