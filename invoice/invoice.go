package invoice

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mohammadv184/gopayment/traits"
)

type Invoice struct {
	uUID          string
	amount        uint32
	transactionID string
	traits.HasDetail
}

func NewInvoice() *Invoice {
	return &Invoice{
		uUID: uuid.New().String(),
	}
}

func (i *Invoice) SetAmount(amount uint32) error {
	if amount > 50000000 {
		return fmt.Errorf("amount must be less than 50,000,000")
	}
	i.amount = amount
	return nil
}
func (i *Invoice) GetAmount() uint32 {
	return i.amount
}

func (i *Invoice) SetUUID(uid ...string) {
	if len(uid) > 0 {
		i.uUID = uid[0]
	}
	if i.uUID == "" {
		i.uUID = uuid.New().String()
	}

}

func (i *Invoice) GetUUID() string {
	if i.uUID == "" {
		i.SetUUID()
	}

	return i.uUID
}
func (i *Invoice) SetTransactionID(transactionID string) {
	i.transactionID = transactionID
}
func (i *Invoice) GetTransactionID() string {
	return i.transactionID
}
