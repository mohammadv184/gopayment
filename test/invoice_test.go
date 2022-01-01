package test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mohammadv184/gopayment/invoice"
	"github.com/stretchr/testify/suite"
)

type InvoiceTestSuite struct {
	suite.Suite
	Invoice *invoice.Invoice
}

func (s *InvoiceTestSuite) SetupTest() {
	s.Invoice = invoice.NewInvoice()
}
func (s *InvoiceTestSuite) TestCreateInvoice() {
	_, err := uuid.Parse(s.Invoice.GetUUID())
	s.Nil(err)
	s.Invoice.SetUUID()
	_, err = uuid.Parse(s.Invoice.GetUUID())
	s.Nil(err)

	s.Invoice.SetUUID("test")
	s.Equal("test", s.Invoice.GetUUID())

	var testInvoice invoice.Invoice
	_, err = uuid.Parse(testInvoice.GetUUID())
	s.Nil(err)

	s.Invoice.SetAmount(100)

	s.Equal(uint32(100), s.Invoice.GetAmount())

	s.Invoice.SetTransactionID("test")
	s.Equal("test", s.Invoice.GetTransactionID())

}

func TestInvoiceTestSuite(t *testing.T) {
	suite.Run(t, new(InvoiceTestSuite))
}
