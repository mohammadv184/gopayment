package invoice

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mohammadv184/gopayment/traits"
)

type Invoice struct {
	UUID          string `json:"uuid"`
	Amount        int32  `json:"amount"`
	TransactionID string `json:"transaction_id"`
	traits.HasDetail
}

func NewInvoice() *Invoice {
	return &Invoice{
		UUID: uuid.New().String(),
	}
}

func (i *Invoice) SetAmount(amount int32) error {
	if amount > 50000000 {
		return fmt.Errorf("amount must be less than 50,000,000")
	}
	i.Amount = amount
	return nil
}

func (i *Invoice) SetUUID(uid ...string) {
	if len(uid) > 0 {
		i.UUID = uid[0]
	}
	if i.UUID == "" {
		i.UUID = uuid.New().String()
	}

}

func (i *Invoice) GetUUID() string {
	if i.UUID == "" {
		i.SetUUID()
	}

	return i.UUID
}
func (i *Invoice) SetTransactionID(transactionID string) {
	i.TransactionID = transactionID
}
func (i *Invoice) GetTransactionID() string {
	return i.TransactionID
}
