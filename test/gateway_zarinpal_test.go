package test

import (
	"encoding/json"
	"testing"

	"github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/gateway/zarinpal"
	"github.com/mohammadv184/gopayment/receipt"

	"github.com/mohammadv184/gopayment/helpers"

	"github.com/mohammadv184/gopayment"
	"github.com/mohammadv184/gopayment/gateway"
	"github.com/mohammadv184/gopayment/test/mock"
	"github.com/stretchr/testify/suite"
)

type GatewayZarinPalTestSuite struct {
	suite.Suite
	Driver     gateway.Driver
	HTTPClient *mock.HTTPClient
}

const (
	zarinpalToken       = "xxxx-xxxx-xxxx-xxxx-xxxx"
	zarinpalDescription = "test"
	zarinpalCallbackURL = "http://localhost:8080/callback"
)

func (s *GatewayZarinPalTestSuite) SetupTest() {
	s.Driver = &zarinpal.Driver{
		MerchantID: zarinpalToken,
		Callback:   zarinpalCallbackURL,
	}
	s.HTTPClient = new(mock.HTTPClient)
	s.Driver.SetClient(s.HTTPClient)
}
func (s *GatewayZarinPalTestSuite) TestPurchaseSuccess() {
	payment := gopayment.NewPayment(s.Driver)
	payment.Amount(1000)
	payment.Description(zarinpalDescription)
	payment.Detail("phone", "+989120000000")
	payment.Detail("email", "example@test.com")

	reqBody := map[string]interface{}{
		"callback_url": zarinpalCallbackURL,
		"description":  payment.GetInvoice().GetDescription(),
		"amount":       payment.GetInvoice().GetAmount(),
		"merchant_id":  zarinpalToken,
		"metadata": map[string]string{
			"phone": "+989120000000",
			"email": "example@test.com",
		},
	}
	respBody := map[string]interface{}{
		"data": map[string]interface{}{
			"authority": "xxxx-xxxx-xxxx-xxxx-xxxx",
		},
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", zarinpal.APIPurchaseURL, reqBody, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 100,
			StatusProperty:     "100 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	err := payment.Purchase()
	s.Nil(err)
	s.HTTPClient.AssertExpectations(s.T())
	s.Equal(payment.GetTransactionID(), respBody["data"].(map[string]interface{})["authority"].(string))
	s.Equal(zarinpal.APIPaymentURL+respBody["data"].(map[string]interface{})["authority"].(string), payment.PayURL())

	s.Equal("GET", payment.PayMethod())
	actualRedirectTMPL, err := payment.RenderRedirectForm()
	s.NoError(err)
	expectedRedirectTMPL, err := helpers.RenderRedirectTemplate("GET", zarinpal.APIPaymentURL+respBody["data"].(map[string]interface{})["authority"].(string), nil)
	s.NoError(err)
	s.Equal(expectedRedirectTMPL, actualRedirectTMPL)

}

func (s *GatewayZarinPalTestSuite) TestVerifySuccess() {
	reqBody := zarinpal.VerifyRequest{
		Authority: zarinpalToken,
		Amount:    "1000",
	}
	respBody := map[string]interface{}{
		"data": map[string]interface{}{
			"ref_id": "xxxx-xxxx-xxxx-xxxx-xxxx",
		},
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", zarinpal.APIVerifyURL, &reqBody, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 100,
			StatusProperty:     "100 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	rec, err := s.Driver.Verify(&reqBody)
	s.Nil(err)
	s.HTTPClient.AssertExpectations(s.T())
	s.IsType(&receipt.Receipt{}, rec)
	s.Equal(zarinpalToken, rec.GetReferenceID())
	s.Equal(s.Driver.GetDriverName(), rec.GetDriver())

}

func (s *GatewayZarinPalTestSuite) TestPurchaseFailed() {
	payment := gopayment.NewPayment(s.Driver)
	payment.Amount(1000)
	payment.Description(zarinpalDescription)
	payment.Detail("phone", "+989120000000")
	payment.Detail("email", "example@test.com")

	reqBody := map[string]interface{}{
		"callback_url": zarinpalCallbackURL,
		"description":  payment.GetInvoice().GetDescription(),
		"amount":       payment.GetInvoice().GetAmount(),
		"merchant_id":  zarinpalToken,
		"metadata": map[string]string{
			"phone": "+989120000000",
			"email": "example@test.com",
		},
	}

	respBody := map[string]interface{}{
		"data": map[string]interface{}{
			"authority": "xxxx-xxxx-xxxx-xxxx-xxxx",
		},
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", zarinpal.APIPurchaseURL, reqBody, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 ok",
			BodyProperty:       respJSON,
		}, errors.ErrInvalidPayment{Message: "example"}).Once()

	err := payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrPurchaseFailed{Message: "500 ok purchase failed"})
	s.Equal("500 ok purchase failed", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", zarinpal.APIPurchaseURL, reqBody, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 error",
			BodyProperty:       respJSON,
		}, nil).Once()

	err = payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrPurchaseFailed{Message: "500 error purchase failed"})
	s.Equal("500 error purchase failed", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", zarinpal.APIPurchaseURL, reqBody, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 100,
			StatusProperty:     "100 ok",
			BodyProperty:       []byte("example"),
		}, nil).Once()

	err = payment.Purchase()
	s.NotNil(err)
	s.Error(err)

	s.HTTPClient.AssertExpectations(s.T())
}
func (s *GatewayZarinPalTestSuite) TestVerifyFailed() {
	reqBody := zarinpal.VerifyRequest{
		Authority: zarinpalToken,
		Amount:    "1000",
	}
	respBody := map[string]interface{}{
		"data": map[string]interface{}{
			"ref_id": "xxxx-xxxx-xxxx-xxxx-xxxx",
		},
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", zarinpal.APIVerifyURL, &reqBody, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 400,
			StatusProperty:     "400 error",
			BodyProperty:       respJSON,
		}, errors.ErrInvalidPayment{Message: "example"}).Once()

	_, err := s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "example"})
	s.Equal("example", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", zarinpal.APIVerifyURL, &reqBody, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 400,
			StatusProperty:     "400 error",
			BodyProperty:       []byte("example"),
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "400 error Invalid payment"})
	s.Equal("400 error Invalid payment", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", zarinpal.APIVerifyURL, &reqBody, map[string]string(nil)).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 100,
			StatusProperty:     "100 error",
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
func TestGatewayZarinPalTestSuite(t *testing.T) {
	suite.Run(t, new(GatewayZarinPalTestSuite))
}
