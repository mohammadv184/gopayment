package invoice

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mohammadv184/gopayment/traits"
)

// Invoice is a struct that holds the invoice data
type Invoice struct {
	uUID          string
	amount        uint32
	transactionID string
	traits.HasDetail
}

// NewInvoice creates a new invoice
func NewInvoice() *Invoice {
	return &Invoice{
		uUID: uuid.New().String(),
	}
}

// SetAmount sets the amount of the invoice
func (i *Invoice) SetAmount(amount uint32) error {
	if amount > 50000000 {
		return fmt.Errorf("amount must be less than 50,000,000")
	}
	i.amount = amount
	return nil
}

// GetAmount returns the amount of the invoice
func (i *Invoice) GetAmount() uint32 {
	return i.amount
}

// SetUUID sets the UUID of the invoice
func (i *Invoice) SetUUID(uid ...string) {
	if len(uid) > 0 {
		i.uUID = uid[0]
	}
	if i.uUID == "" {
		i.uUID = uuid.New().String()
	}

}

// GetUUID returns the UUID of the invoice
func (i *Invoice) GetUUID() string {
	if i.uUID == "" {
		i.SetUUID()
	}

	return i.uUID
}

// SetTransactionID sets the transaction ID of the invoice
func (i *Invoice) SetTransactionID(transactionID string) {
	i.transactionID = transactionID
}

// GetTransactionID returns the transaction ID of the invoice
func (i *Invoice) GetTransactionID() string {
	return i.transactionID
}
