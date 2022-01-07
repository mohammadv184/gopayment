package test

import (
	"encoding/json"
	"github.com/mohammadv184/gopayment/helpers"
	"testing"

	"github.com/google/uuid"
	"github.com/mohammadv184/gopayment"
	"github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/gateway"
	"github.com/mohammadv184/gopayment/gateway/payping"
	"github.com/mohammadv184/gopayment/receipt"
	"github.com/mohammadv184/gopayment/test/mock"
	"github.com/stretchr/testify/suite"
)

type GatewayPayPingTestSuite struct {
	suite.Suite
	Driver     gateway.Driver
	HTTPClient *mock.HTTPClient
}

const (
	token       = "xxxx-xxxx-xxxx-xxxx-xxxx"
	description = "test"
	callbackURL = "http://localhost:8080/callback"
)

func (s *GatewayPayPingTestSuite) SetupTest() {
	s.Driver = &payping.Driver{
		Token:    token,
		Callback: callbackURL,
	}
	s.HTTPClient = new(mock.HTTPClient)
	s.Driver.SetClient(s.HTTPClient)
}
func (s *GatewayPayPingTestSuite) TestPurchaseSuccess() {
	payment := gopayment.NewPayment(s.Driver)
	payment.Amount(1000)
	payment.Description(description)
	payment.Detail("phone", "+989120000000")
	payment.Detail("name", "John Doe")

	reqBody := map[string]interface{}{
		"returnUrl":     callbackURL,
		"description":   payment.GetInvoice().GetDescription(),
		"amount":        payment.GetInvoice().GetAmount(),
		"clientRefId":   payment.GetInvoice().GetUUID(),
		"payerIdentity": "+989120000000",
		"payerName":     "John Doe",
	}
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}
	respBody := map[string]interface{}{
		"code": "cosuh87gf",
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", payping.APIPurchaseURL, reqBody, header).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	err := payment.Purchase()
	s.Nil(err)
	s.HTTPClient.AssertExpectations(s.T())
	s.Equal(payment.GetTransactionID(), respBody["code"])
	s.Equal(payping.APIPaymentURL+respBody["code"].(string), payment.PayURL())

	s.Equal("GET", payment.PayMethod())
	actualRedirectTMPL, err := payment.RenderRedirectForm()
	s.NoError(err)
	expectedRedirectTMPL, err := helpers.RenderRedirectTemplate("GET", payping.APIPaymentURL+respBody["code"].(string), nil)
	s.NoError(err)
	s.Equal(expectedRedirectTMPL, actualRedirectTMPL)

}
func (s *GatewayPayPingTestSuite) TestVerifySuccess() {
	reqBody := payping.VerifyRequest{
		RefID:  uuid.New().String(),
		Amount: "1000",
	}
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}
	respBody := map[string]interface{}{
		"amount":      1000,
		"cardNumber":  "string",
		"cardHashPan": "string",
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", payping.APIVerifyURL, &reqBody, header).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	rec, err := s.Driver.Verify(&reqBody)
	s.Nil(err)
	s.HTTPClient.AssertExpectations(s.T())
	s.IsType(&receipt.Receipt{}, rec)
	s.Equal(reqBody.RefID, rec.GetReferenceID())
	s.Equal(s.Driver.GetDriverName(), rec.GetDriver())
	s.True(rec.Has("cardNumber"))
	s.Equal(respBody["cardNumber"], rec.GetDetail("cardNumber"))

}
func (s *GatewayPayPingTestSuite) TestPurchaseFailed() {
	payment := gopayment.NewPayment(s.Driver)
	payment.Amount(1000)
	payment.Description(description)
	payment.Detail("email", "test@example.com")
	payment.Detail("name", "John Doe")

	reqBody := map[string]interface{}{
		"returnUrl":     callbackURL,
		"description":   payment.GetInvoice().GetDescription(),
		"amount":        payment.GetInvoice().GetAmount(),
		"clientRefId":   payment.GetInvoice().GetUUID(),
		"payerIdentity": "test@example.com",
		"payerName":     "John Doe",
	}
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}
	respBody := map[string]interface{}{
		"code": "cosuh87gf",
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", payping.APIPurchaseURL, reqBody, header).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, errors.ErrInvalidPayment{Message: "example"}).Once()

	err := payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "example"})
	s.Equal("example", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", payping.APIPurchaseURL, reqBody, header).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 error",
			BodyProperty:       respJSON,
		}, nil).Once()

	err = payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrPurchaseFailed{Message: "Purchase failed"})
	s.Equal("Purchase failed", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", payping.APIPurchaseURL, reqBody, header).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte("example"),
		}, nil).Once()

	err = payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.HTTPClient.AssertExpectations(s.T())
}
func (s *GatewayPayPingTestSuite) TestVerifyFailed() {
	reqBody := payping.VerifyRequest{
		RefID:  uuid.New().String(),
		Amount: "1000",
	}
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}
	respBody := map[string]interface{}{
		"1": "error",
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", payping.APIVerifyURL, &reqBody, header).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 400,
			StatusProperty:     "400 error",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err := s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "error"})
	s.Equal("error", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", payping.APIVerifyURL, &reqBody, header).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 400,
			StatusProperty:     "400 error",
			BodyProperty:       []byte("example"),
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "error in verify payment"})
	s.Equal("error in verify payment", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	_, err = s.Driver.Verify("example")
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInternal{Message: "vReq is not of type VerifyRequest"})
	s.Equal("vReq is not of type VerifyRequest", err.Error())
}
func TestGatewayPayPingTestSuite(t *testing.T) {
	suite.Run(t, new(GatewayPayPingTestSuite))
}
