// Package gopayment multi gateway payment package for Golang
package gopayment

import (
	"github.com/mohammadv184/gopayment/gateway"
	"github.com/mohammadv184/gopayment/invoice"
	httpClient "github.com/mohammadv184/gopayment/pkg/http"
)

// Version is the version of gopayment
const Version = "v1.5.0"

// Payment is the payment main struct of gopayment
type Payment struct {
	driver  gateway.Driver
	invoice *invoice.Invoice
}

// Amount set the amount of payment invoice
func (p *Payment) Amount(amount int) *Payment {
	p.invoice.SetAmount(uint32(amount))
	return p.returnThis()
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
	return p.driver.PayURL(p.invoice)
}

// PayMethod returns the Request Method to be used to pay the invoice.
func (p *Payment) PayMethod() string {
	return p.driver.PayMethod()
}

// Client sets the driver http client.
func (p *Payment) Client(client httpClient.Client) *Payment {
	p.driver.SetClient(client)
	return p.returnThis()
}

// GetInvoice return the payment invoice
func (p *Payment) GetInvoice() *invoice.Invoice {
	return p.invoice
}

// GetTransactionID return the payment transaction id
func (p *Payment) GetTransactionID() string {
	return p.invoice.GetTransactionID()
}

// Description set the payment description
func (p *Payment) Description(description string) *Payment {
	p.invoice.SetDescription(description)
	return p.returnThis()
}

// Detail set the payment detail
func (p *Payment) Detail(key string, value string) *Payment {
	p.invoice.Detail(key, value)
	return p.returnThis()
}

func (p *Payment) returnThis() *Payment {
	return p
}

// NewPayment create a new payment
func NewPayment(driver gateway.Driver) *Payment {
	return &Payment{
		driver:  driver,
		invoice: invoice.NewInvoice(),
	}
}
