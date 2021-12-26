package test

import (
	"github.com/mohammadv184/gopayment"
	"github.com/mohammadv184/gopayment/drivers"
	"github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/invoice"
	"github.com/mohammadv184/gopayment/receipt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GoPaymentTestSuite struct {
	suite.Suite
	Gateway drivers.Driver
}

func (s *GoPaymentTestSuite) SetupTest() {
	s.Gateway = &gateway{}
}
func (s *GoPaymentTestSuite) TestCreatePayment() {
	payment := gopayment.NewPayment(s.Gateway)
	err := payment.Amount(100)
	s.Nil(err)

	err = payment.Amount(50000001)
	s.NotNil(err)

	err = payment.Purchase()
	s.Nil(err)

	err = payment.Amount(99)
	s.Nil(err)
	err = payment.Purchase()
	s.NotNil(err)
	s.EqualError(err, "amount is less than 100")

	s.Equal(payment.GetTransactionID(), payment.GetInvoice().GetTransactionID())
	s.Equal("gateway.com/"+payment.GetTransactionID(), payment.PayUrl())

}
func TestGoPaymentTestSuite(t *testing.T) {
	suite.Run(t, new(GoPaymentTestSuite))
}

// gateway driver mock
type gateway struct {
}

func (g gateway) GetDriverName() string {
	return "MockGateway"
}
func (g *gateway) Purchase(invoice *invoice.Invoice) (string, error) {
	if invoice.GetAmount() < 100 {
		return "", errors.ErrPurchaseFailed{
			Message: "amount is less than 100",
		}
	}
	return invoice.GetUUID(), nil
}
func (g *gateway) PayUrl(invoice *invoice.Invoice) string {
	return "gateway.com/" + invoice.GetTransactionID()
}
func (g *gateway) Verify(interface{}) (*receipt.Receipt, error) {
	return receipt.NewReceipt("test", g.GetDriverName()), nil
}
