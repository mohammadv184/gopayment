package receipt

import (
	"time"

	"github.com/mohammadv184/gopayment/trait"
)

// Receipt is a struct that holds the information of a receipt
type Receipt struct {
	referenceID string
	date        time.Time
	driver      string
	trait.HasDetail
}

// NewReceipt creates a new receipt
func NewReceipt(reference string, driver string) *Receipt {
	return &Receipt{
		referenceID: reference,
		driver:      driver,
		date:        time.Now(),
	}
}

// GetReferenceID returns the reference id of the receipt
func (r *Receipt) GetReferenceID() string {
	return r.referenceID
}

// GetDriver returns the driver of the payment
func (r *Receipt) GetDriver() string {
	return r.driver
}

// GetDate returns the date of the receipt
func (r *Receipt) GetDate() time.Time {
	return r.date
}
