package gopayment

import (
	"github.com/mohammadv184/gopayment/drivers"
	"github.com/mohammadv184/gopayment/invoice"
	_ "github.com/mohammadv184/gopayment/pkg/http"
)

type Payment struct {
	driver  drivers.Driver
	invoice *invoice.Invoice
}

func (p *Payment) Amount(amount int) error {
	err := p.invoice.SetAmount(uint32(amount))
	if err != nil {
		return err
	}
	return nil
}
func (p *Payment) Purchase() error {
	transID, err := p.driver.Purchase(p.invoice)
	if err != nil {
		return err
	}
	p.invoice.SetTransactionID(transID)
	return nil
}
func (p *Payment) PayUrl() string {
	return p.driver.PayUrl(p.invoice)
}
func (p *Payment) GetInvoice() *invoice.Invoice {
	return p.invoice
}

func NewPayment(driver drivers.Driver) *Payment {
	return &Payment{
		driver:  driver,
		invoice: invoice.NewInvoice(),
	}
}
