package invoice

import (
	"github.com/google/uuid"
	"github.com/mohammadv184/gopayment/trait"
)

// Invoice is a struct that holds the invoice data
type Invoice struct {
	uUID          string
	amount        uint32
	transactionID string
	description   string
	trait.HasDetail
}

// NewInvoice creates a new invoice
func NewInvoice() *Invoice {
	return &Invoice{
		uUID: uuid.New().String(),
	}
}

// SetAmount sets the amount of the invoice
func (i *Invoice) SetAmount(amount uint32) *Invoice {
	i.amount = amount
	return i.returnThis()
}

// GetAmount returns the amount of the invoice
func (i *Invoice) GetAmount() uint32 {
	return i.amount
}

// SetUUID sets the UUID of the invoice
func (i *Invoice) SetUUID(uid ...string) *Invoice {
	if len(uid) > 0 {
		i.uUID = uid[0]
	}
	if i.uUID == "" {
		i.uUID = uuid.New().String()
	}
	return i.returnThis()
}

// GetUUID returns the UUID of the invoice
func (i *Invoice) GetUUID() string {
	if i.uUID == "" {
		i.SetUUID()
	}

	return i.uUID
}

// SetTransactionID sets the transaction ID of the invoice
func (i *Invoice) SetTransactionID(transactionID string) *Invoice {
	i.transactionID = transactionID
	return i.returnThis()
}

// GetTransactionID returns the transaction ID of the invoice
func (i *Invoice) GetTransactionID() string {
	return i.transactionID
}

// SetDescription sets the description of the invoice
func (i *Invoice) SetDescription(description string) *Invoice {
	i.description = description
	return i.returnThis()
}

// GetDescription returns the description of the invoice
func (i *Invoice) GetDescription() string {
	return i.description
}
func (i *Invoice) returnThis() *Invoice {
	return i
}
