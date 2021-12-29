// Package gopayment multi gateway payment package for Golang
package gopayment

import (
	"github.com/mohammadv184/gopayment/drivers"
	"github.com/mohammadv184/gopayment/invoice"
	_ "github.com/mohammadv184/gopayment/pkg/http"
)

// Version is the version of gopayment
const Version = "v1.1.0"

// Payment is the payment main struct of gopayment
type Payment struct {
	driver  drivers.Driver
	invoice *invoice.Invoice
}

// Amount set the amount of payment invoice
func (p *Payment) Amount(amount int) error {
	err := p.invoice.SetAmount(uint32(amount))
	if err != nil {
		return err
	}
	return nil
}

// Purchase send the purchase request to the payment gateway
func (p *Payment) Purchase() error {
	transID, err := p.driver.Purchase(p.invoice)
	if err != nil {
		return err
	}
	p.invoice.SetTransactionID(transID)
	return nil
}

// PayURL return the payment URL
func (p *Payment) PayURL() string {
	return p.driver.PayUrl(p.invoice)
}

// PayMethod returns the Request Method to be used to pay the invoice.
func (p *Payment) PayMethod() string {
	return p.driver.PayMethod()
}

// GetInvoice return the payment invoice
func (p *Payment) GetInvoice() *invoice.Invoice {
	return p.invoice
}

// GetTransactionID return the payment transaction id
func (p *Payment) GetTransactionID() string {
	return p.invoice.GetTransactionID()
}

// NewPayment create a new payment
func NewPayment(driver drivers.Driver) *Payment {
	return &Payment{
		driver:  driver,
		invoice: invoice.NewInvoice(),
	}
}
