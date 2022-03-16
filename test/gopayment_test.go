package test

import (
	"testing"

	httpClient "github.com/mohammadv184/gopayment/pkg/http"

	"github.com/mohammadv184/gopayment"
	"github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/gateway"
	"github.com/mohammadv184/gopayment/invoice"
	"github.com/mohammadv184/gopayment/receipt"
	"github.com/stretchr/testify/suite"
)

type GoPaymentTestSuite struct {
	suite.Suite
	Gateway gateway.Driver
}

func (s *GoPaymentTestSuite) SetupTest() {
	s.Gateway = &Gateway{}
}
func (s *GoPaymentTestSuite) TestCreatePayment() {
	payment := gopayment.NewPayment(s.Gateway)
	payment.Amount(100)

	err := payment.Purchase()
	s.Nil(err)

	payment.Amount(99)
	payment.Client(httpClient.NewHTTP())
	err = payment.Purchase()
	s.NotNil(err)
	s.EqualError(err, "amount is less than 100")

	s.Equal(payment.GetTransactionID(), payment.GetInvoice().GetTransactionID())
	s.Equal("GET", payment.PayMethod())
	s.Equal("Gateway.com/"+payment.GetTransactionID(), payment.PayURL())

	_, err = payment.RenderRedirectForm()
	s.NoError(err)

}
func TestGoPaymentTestSuite(t *testing.T) {
	suite.Run(t, new(GoPaymentTestSuite))
}

// Gateway driver mock
type Gateway struct {
}

func (g Gateway) GetDriverName() string {
	return "MockGateway"
}
func (g *Gateway) Purchase(invoice *invoice.Invoice) (string, error) {
	if invoice.GetAmount() < 100 {
		return "", errors.ErrPurchaseFailed{
			Message: "amount is less than 100",
		}
	}
	return invoice.GetUUID(), nil
}
func (g *Gateway) PayURL(invoice *invoice.Invoice) string {
	return "Gateway.com/" + invoice.GetTransactionID()
}
func (g *Gateway) Verify(interface{}) (*receipt.Receipt, error) {
	return receipt.NewReceipt("test", g.GetDriverName()), nil
}
func (g *Gateway) PayMethod() string {
	return "GET"
}

// SetClient sets the http client
func (g *Gateway) SetClient(_ httpClient.Client) {}
func (g *Gateway) RenderRedirectForm(_ *invoice.Invoice) (string, error) {
	return "", nil
}
