package test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/mohammadv184/gopayment"
	"github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/gateway"
	"github.com/mohammadv184/gopayment/gateway/asanpardakht"
	"github.com/mohammadv184/gopayment/helpers"
	"github.com/mohammadv184/gopayment/receipt"
	"github.com/mohammadv184/gopayment/test/mock"
	"github.com/stretchr/testify/suite"
)

type GatewayAsanPardakhtTestSuite struct {
	suite.Suite
	Driver     gateway.Driver
	HTTPClient *mock.HTTPClient
}

const (
	asanpardakhtConfigID    = "xxxx-xxxx-xxxx-xxxx-xxxx"
	asanpardakhtDescription = "test"
	asanpardakhtCallbackURL = "http://localhost:8080/callback"
	asanpardakhtUsername    = "test"
	asanpardakhtPassword    = "test"
)

func (s *GatewayAsanPardakhtTestSuite) SetupTest() {
	s.Driver = &asanpardakht.Driver{
		MerchantConfigID: asanpardakhtConfigID,
		Callback:         asanpardakhtCallbackURL,
		Username:         asanpardakhtUsername,
		Password:         asanpardakhtPassword,
	}
	s.HTTPClient = new(mock.HTTPClient)
	s.Driver.SetClient(s.HTTPClient)
}
func (s *GatewayAsanPardakhtTestSuite) TestPurchaseSuccess() {
	payment := gopayment.NewPayment(s.Driver)
	payment.Amount(1000)
	payment.Description(asanpardakhtDescription)

	reqBody := map[string]interface{}{
		"callbackURL":             asanpardakhtCallbackURL,
		"additionalData":          payment.GetInvoice().GetDescription(),
		"amountInRials":           payment.GetInvoice().GetAmount() * 10,
		"merchantConfigurationId": asanpardakhtConfigID,
		"localInvoiceId":          payment.GetInvoice().GetUUID(),
		"serviceTypeId":           1,
		"localDate":               time.Now().Format("20060102 150405"),
		"paymentId":               0,
	}
	reqHeader := map[string]string{
		"usr": asanpardakhtUsername,
		"pwd": asanpardakhtPassword,
	}
	resBody := []byte(`xxxx-xxxx-xxxx-xxxx-xxxx`)
	s.HTTPClient.
		On("Post", asanpardakht.APIPurchaseURL, reqBody, reqHeader).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       resBody,
		}, nil).Once()

	err := payment.Purchase()
	s.Nil(err)
	s.HTTPClient.AssertExpectations(s.T())
	s.Equal(payment.GetTransactionID(), "xxxx-xxxx-xxxx-xxxx-xxxx")
	s.Equal(asanpardakht.APIPaymentURL, payment.PayURL())

	s.Equal("POST", payment.PayMethod())
	actualRedirectTMPL, err := payment.RenderRedirectForm()
	s.NoError(err)
	expectedRedirectTMPL, err := helpers.RenderRedirectTemplate(
		"POST",
		asanpardakht.APIPaymentURL,
		map[string]string{
			"RefId": payment.GetInvoice().GetUUID(),
		})
	s.NoError(err)
	s.Equal(expectedRedirectTMPL, actualRedirectTMPL)

}

func (s *GatewayAsanPardakhtTestSuite) TestVerifySuccess() {
	reqBody := asanpardakht.VerifyRequest{
		InvoiceID: "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	respBody := map[string]interface{}{
		"PayGateTranID": "xxxx-xxxx-xxxx-xxxx-xxxx",
		"rrn":           "xxxx-xxxx-xxxx-xxxx-xxxx",
		"payment": map[string]interface{}{
			"cardNumber": "xxxx-xxxx-xxxx-xxxx-xxxx",
			"refID":      "xxxx-xxxx-xxxx-xxxx-xxxx",
		},
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Get", asanpardakht.APITranResultURL, map[string]string{
			"localInvoiceId":          "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	s.HTTPClient.
		On("Post", asanpardakht.APIVerifyURL, map[string]string{
			"payGateTranID":           "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte{},
		}, nil).Once()
	s.HTTPClient.
		On("Post", asanpardakht.APISettlementURL, map[string]string{
			"payGateTranID":           "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte{},
		}, nil).Once()

	rec, err := s.Driver.Verify(&reqBody)
	s.Nil(err)
	s.HTTPClient.AssertExpectations(s.T())
	s.IsType(&receipt.Receipt{}, rec)
	s.Equal(asanpardakhtConfigID, rec.GetReferenceID())
	s.Equal(s.Driver.GetDriverName(), rec.GetDriver())

}

func (s *GatewayAsanPardakhtTestSuite) TestPurchaseFailed() {
	payment := gopayment.NewPayment(s.Driver)
	payment.Amount(1000)
	payment.Description(asanpardakhtDescription)

	reqBody := map[string]interface{}{
		"callbackURL":             asanpardakhtCallbackURL,
		"additionalData":          payment.GetInvoice().GetDescription(),
		"amountInRials":           payment.GetInvoice().GetAmount() * 10,
		"merchantConfigurationId": asanpardakhtConfigID,
		"localInvoiceId":          payment.GetInvoice().GetUUID(),
		"serviceTypeId":           1,
		"localDate":               time.Now().Format("20060102 150405"),
		"paymentId":               0,
	}
	reqHeader := map[string]string{
		"usr": asanpardakhtUsername,
		"pwd": asanpardakhtPassword,
	}
	resBody := []byte(`xxxx-xxxx-xxxx-xxxx-xxxx`)
	s.HTTPClient.
		On("Post", asanpardakht.APIPurchaseURL, reqBody, reqHeader).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 internal server error",
			BodyProperty:       resBody,
		}, nil).Once()

	err := payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrPurchaseFailed{Message: "500 internal server error"})
	s.Equal("500 internal server error", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", asanpardakht.APIPurchaseURL, reqBody, reqHeader).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 internal server error",
			BodyProperty:       resBody,
		}, fmt.Errorf("example")).Once()

	err = payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.Equal("example", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
}

func (s *GatewayAsanPardakhtTestSuite) TestVerifyFailed() {
	reqBody := asanpardakht.VerifyRequest{
		InvoiceID: "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	respBody := map[string]interface{}{
		"PayGateTranID": 123,
		"rrn":           "xxxx-xxxx-xxxx-xxxx-xxxx",
		"payment": map[string]interface{}{
			"cardNumber": "xxxx-xxxx-xxxx-xxxx-xxxx",
			"refID":      "xxxx-xxxx-xxxx-xxxx-xxxx",
		},
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Get", asanpardakht.APITranResultURL, map[string]string{
			"localInvoiceId":          "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err := s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInternal{Message: "PayGateTranID is not of type string"})
	s.Equal("PayGateTranID is not of type string", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	respBody = map[string]interface{}{
		"PayGateTranID": "xxxx-xxxx-xxxx-xxxx-xxxx",
		"rrn":           "xxxx-xxxx-xxxx-xxxx-xxxx",
		"payment": map[string]interface{}{
			"cardNumber": "xxxx-xxxx-xxxx-xxxx-xxxx",
			"refID":      "xxxx-xxxx-xxxx-xxxx-xxxx",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Get", asanpardakht.APITranResultURL, map[string]string{
			"localInvoiceId":          "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 internal server error",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "500 internal server error"})
	s.Equal("500 internal server error", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Get", asanpardakht.APITranResultURL, map[string]string{
			"localInvoiceId":          "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 internal server error",
			BodyProperty:       respJSON,
		}, fmt.Errorf("example")).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.Equal("example", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Get", asanpardakht.APITranResultURL, map[string]string{
			"localInvoiceId":          "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	s.HTTPClient.
		On("Post", asanpardakht.APIVerifyURL, map[string]string{
			"payGateTranID":           "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 internal server error",
			BodyProperty:       []byte{},
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.Equal("500 internal server error", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Get", asanpardakht.APITranResultURL, map[string]string{
			"localInvoiceId":          "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	s.HTTPClient.
		On("Post", asanpardakht.APIVerifyURL, map[string]string{
			"payGateTranID":           "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte{},
		}, fmt.Errorf("example")).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.Equal("example", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Get", asanpardakht.APITranResultURL, map[string]string{
			"localInvoiceId":          "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	s.HTTPClient.
		On("Post", asanpardakht.APIVerifyURL, map[string]string{
			"payGateTranID":           "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte{},
		}, nil).Once()
	s.HTTPClient.
		On("Post", asanpardakht.APISettlementURL, map[string]string{
			"payGateTranID":           "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 internal server error",
			BodyProperty:       []byte{},
		}, nil).Once()
	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.Equal("500 internal server error", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Get", asanpardakht.APITranResultURL, map[string]string{
			"localInvoiceId":          "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	s.HTTPClient.
		On("Post", asanpardakht.APIVerifyURL, map[string]string{
			"payGateTranID":           "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte{},
		}, nil).Once()
	s.HTTPClient.
		On("Post", asanpardakht.APISettlementURL, map[string]string{
			"payGateTranID":           "xxxx-xxxx-xxxx-xxxx-xxxx",
			"merchantConfigurationId": asanpardakhtConfigID,
		}, map[string]string{
			"usr": asanpardakhtUsername,
			"pwd": asanpardakhtPassword,
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       []byte{},
		}, fmt.Errorf("example")).Once()
	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.Equal("example", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	_, err = s.Driver.Verify("example")
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInternal{Message: "vReq is not of type VerifyRequest"})
	s.Equal("vReq is not of type VerifyRequest", err.Error())
}
func TestGatewayAsanPardakhtTestSuite(t *testing.T) {
	suite.Run(t, new(GatewayAsanPardakhtTestSuite))
}
