package receipt

import (
	"github.com/mohammadv184/gopayment/traits"
	"time"
)

type Receipt struct {
	referenceID string
	date        time.Time
	driver      string
	traits.HasDetail
}

func NewReceipt(reference string, driver string) *Receipt {
	return &Receipt{
		referenceID: reference,
		driver:      driver,
		date:        time.Now(),
	}
}
func (r *Receipt) GetReferenceID() string {
	return r.referenceID
}
func (r *Receipt) GetDriver() string {
	return r.driver
}
func (r *Receipt) GetDate() time.Time {
	return r.date
}
