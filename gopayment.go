package gopayment

import (
	"github.com/mohammadv184/gopayment/drivers"
	"github.com/mohammadv184/gopayment/invoice"
	"log"
)

type Payment struct {
	driver  drivers.Driver
	invoice *invoice.Invoice
}

func (p *Payment) Amount(amount int) error {
	err := p.invoice.SetAmount(int32(amount))
	if err != nil {
		return err
	}
	return nil
}
func (p *Payment) Purchase() {
	transID, err := p.driver.Purchase(p.invoice)
	if err != nil {
		log.Println(err)
	}
	p.invoice.SetTransactionID(transID)
}
func (p *Payment) PayUrl() string {
	return p.driver.PayUrl(p.invoice)
}
func NewPayment(driver drivers.Driver) *Payment {
	return &Payment{
		driver:  driver,
		invoice: invoice.NewInvoice(),
	}
}
