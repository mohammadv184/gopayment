package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/helpers"
	"github.com/mohammadv184/gopayment/receipt"

	"github.com/mohammadv184/gopayment"
	"github.com/mohammadv184/gopayment/gateway"
	"github.com/mohammadv184/gopayment/gateway/zibal"
	"github.com/mohammadv184/gopayment/test/mock"
	"github.com/stretchr/testify/suite"
)

type GatewayZibalTestSuite struct {
	suite.Suite
	Driver     gateway.Driver
	HTTPClient *mock.HTTPClient
}

const (
	zibalMerchant    = "xxxx-xxxx-xxxx-xxxx-xxxx"
	zibalDescription = "test"
	zibalCallbackURL = "http://localhost:8080/callback"
)

func (s *GatewayZibalTestSuite) SetupTest() {
	s.Driver = &zibal.Driver{
		Merchant: zibalMerchant,
		Callback: zibalCallbackURL,
	}
	s.HTTPClient = new(mock.HTTPClient)
	s.Driver.SetClient(s.HTTPClient)
}
func (s *GatewayZibalTestSuite) TestPurchaseSuccess() {
	payment := gopayment.NewPayment(s.Driver)
	payment.Amount(1000)
	payment.Description(zibalDescription)
	payment.Detail("phone", "+989120000000")

	reqBody := map[string]interface{}{
		"merchant":    zibalMerchant,
		"description": payment.GetInvoice().GetDescription(),
		"amount":      payment.GetInvoice().GetAmount(),
		"orderId":     payment.GetInvoice().GetUUID(),
		"mobile":      "+989120000000",
		"callbackUrl": zibalCallbackURL,
	}

	respBody := map[string]interface{}{
		"trackId": 1515143154135,
		"result":  100,
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", zibal.APIPurchaseURL, reqBody, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	err := payment.Purchase()
	s.Nil(err)
	s.HTTPClient.AssertExpectations(s.T())
	s.Equal(payment.GetTransactionID(), fmt.Sprint(respBody["trackId"]))
	s.Equal(zibal.APIPaymentURL+fmt.Sprint(respBody["trackId"]), payment.PayURL())

	s.Equal("GET", payment.PayMethod())
	actualRedirectTMPL, err := payment.RenderRedirectForm()
	s.NoError(err)
	expectedRedirectTMPL, err := helpers.RenderRedirectTemplate("GET", zibal.APIPaymentURL+fmt.Sprint(respBody["trackId"]), nil)
	s.NoError(err)
	s.Equal(expectedRedirectTMPL, actualRedirectTMPL)

}

func (s *GatewayZibalTestSuite) TestVerifySuccess() {
	reqBody := zibal.VerifyRequest{
		TrackID: uuid.New().String(),
	}
	respBody := map[string]interface{}{
		"result":     100,
		"refNumber":  "string",
		"cardNumber": "string",
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", zibal.APIVerifyURL, map[string]string{
			"merchant": zibalMerchant,
			"trackId":  reqBody.TrackID,
		}, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	rec, err := s.Driver.Verify(&reqBody)
	s.Nil(err)
	s.HTTPClient.AssertExpectations(s.T())
	s.IsType(&receipt.Receipt{}, rec)
	s.Equal(respBody["refNumber"], rec.GetReferenceID())
	s.Equal(s.Driver.GetDriverName(), rec.GetDriver())
	s.True(rec.Has("cardNumber"))
	s.Equal(respBody["cardNumber"], rec.GetDetail("cardNumber"))

}

func (s *GatewayZibalTestSuite) TestPurchaseFailed() {
	payment := gopayment.NewPayment(s.Driver)
	payment.Amount(1000)
	payment.Description(zibalDescription)
	payment.Detail("phone", "+989120000000")

	reqBody := map[string]interface{}{
		"merchant":    zibalMerchant,
		"description": payment.GetInvoice().GetDescription(),
		"amount":      payment.GetInvoice().GetAmount(),
		"orderId":     payment.GetInvoice().GetUUID(),
		"mobile":      "+989120000000",
		"callbackUrl": zibalCallbackURL,
	}

	respBody := map[string]interface{}{
		"trackId": 1515143154135,
		"result":  101, // failed
		"message": "example",
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", zibal.APIPurchaseURL, reqBody, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	err := payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrPurchaseFailed{Message: "200 ok example"})
	s.Equal("200 ok example", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", zibal.APIPurchaseURL, reqBody, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 internal server error",
			BodyProperty:       respJSON,
		}, nil).Once()

	err = payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrPurchaseFailed{Message: "500 internal server error example"})
	s.Equal("500 internal server error example", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", zibal.APIPurchaseURL, reqBody, map[string]string(nil)).
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

func (s *GatewayZibalTestSuite) TestVerifyFailed() {
	reqBody := zibal.VerifyRequest{
		TrackID: uuid.New().String(),
	}
	respBody := map[string]interface{}{
		"result":     101, // failed
		"message":    "example",
		"refNumber":  "string",
		"cardNumber": "string",
	}
	respJSON, _ := json.Marshal(respBody)

	s.HTTPClient.
		On("Post", zibal.APIVerifyURL, map[string]string{
			"merchant": zibalMerchant,
			"trackId":  reqBody.TrackID,
		}, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err := s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrPurchaseFailed{Message: "200 ok example"})
	s.Equal("200 ok example", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", zibal.APIVerifyURL, map[string]string{
			"merchant": zibalMerchant,
			"trackId":  reqBody.TrackID,
		}, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte("example"),
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.HTTPClient.AssertExpectations(s.T())

	_, err = s.Driver.Verify("example")
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInternal{Message: "vReq is not of type VerifyRequest"})
	s.Equal("vReq is not of type VerifyRequest", err.Error())
}
func TestGatewayZibalTestSuite(t *testing.T) {
	suite.Run(t, new(GatewayZibalTestSuite))
}
