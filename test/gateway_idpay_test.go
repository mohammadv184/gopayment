package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mohammadv184/gopayment"
	"github.com/mohammadv184/gopayment/errors"
	"github.com/mohammadv184/gopayment/gateway/idpay"
	"github.com/mohammadv184/gopayment/helpers"
	"github.com/mohammadv184/gopayment/receipt"
	"github.com/mohammadv184/gopayment/test/mock"
	"github.com/stretchr/testify/suite"
)

type GatewayIDPayTestSuite struct {
	suite.Suite
	Driver     *idpay.Driver
	HTTPClient *mock.HTTPClient
}

const (
	idpayMerchant    = "xxxx-xxxx-xxxx-xxxx-xxxx"
	idpayDescription = "test"
	idpayCallbackURL = "http://localhost:8080/callback"
)

func (s *GatewayIDPayTestSuite) SetupTest() {
	s.Driver = &idpay.Driver{
		MerchantID: idpayMerchant,
		Callback:   idpayCallbackURL,
	}
	s.HTTPClient = new(mock.HTTPClient)
	s.Driver.SetClient(s.HTTPClient)
}
func (s *GatewayIDPayTestSuite) TestPurchaseSuccess() {
	payment := gopayment.NewPayment(s.Driver)
	payment.Amount(1000)
	payment.Description(idpayDescription)
	payment.Detail("phone", "+989120000000")
	payment.Detail("email", "example@example.com")
	payment.Detail("name", "example")

	reqBody := map[string]interface{}{
		"desc":     payment.GetInvoice().GetDescription(),
		"amount":   payment.GetInvoice().GetAmount(),
		"order_id": payment.GetInvoice().GetUUID(),
		"phone":    "+989120000000",
		"mail":     "example@example.com",
		"name":     "example",
		"callback": idpayCallbackURL,
	}
	reqHeader := map[string]string{
		"X-API-KEY": idpayMerchant,
		"X-SANDBOX": "false",
	}
	respBody := map[string]interface{}{
		"id": "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIPurchaseURL, reqBody, reqHeader).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 201,
			StatusProperty:     "201 created",
			BodyProperty:       respJSON,
		}, nil).Once()

	err := payment.Purchase()
	s.Nil(err)
	s.HTTPClient.AssertExpectations(s.T())
	s.Equal(payment.GetTransactionID(), fmt.Sprint(respBody["id"]))
	s.Equal(idpay.APIPaymentURL+fmt.Sprint(respBody["id"]), payment.PayURL())
	s.Driver.Sandbox = true
	s.Equal(idpay.APISandBoxPaymentURL+fmt.Sprint(respBody["id"]), payment.PayURL())
	s.Driver.Sandbox = false
	s.Equal("GET", payment.PayMethod())
	actualRedirectTMPL, err := payment.RenderRedirectForm()
	s.NoError(err)
	expectedRedirectTMPL, err := helpers.RenderRedirectTemplate("GET", idpay.APIPaymentURL+fmt.Sprint(respBody["id"]), nil)
	s.NoError(err)
	s.Equal(expectedRedirectTMPL, actualRedirectTMPL)

}

func (s *GatewayIDPayTestSuite) TestVerifySuccess() {
	reqBody := idpay.VerifyRequest{
		RefID: "xxxx-xxxx-xxxx-xxxx-xxxx",
		ID:    "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	respBody := map[string]interface{}{
		"status":   100,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	rec, err := s.Driver.Verify(&reqBody)
	s.Nil(err)
	s.HTTPClient.AssertExpectations(s.T())
	s.IsType(&receipt.Receipt{}, rec)
	s.Equal(respBody["track_id"], rec.GetReferenceID())
	s.Equal(s.Driver.GetDriverName(), rec.GetDriver())
	s.True(rec.Has("cardNumber"))
	s.Equal(respBody["payment"].(map[string]interface{})["card_no"].(string), rec.GetDetail("cardNumber"))
	s.True(rec.Has("HashedCardNumber"))
	s.Equal(respBody["payment"].(map[string]interface{})["hashed_card_no"].(string), rec.GetDetail("HashedCardNumber"))

}

func (s *GatewayIDPayTestSuite) TestPurchaseFailed() {
	payment := gopayment.NewPayment(s.Driver)
	payment.Amount(1000)
	payment.Description(idpayDescription)
	payment.Detail("phone", "+989120000000")
	payment.Detail("email", "example@example.com")
	payment.Detail("name", "example")

	reqBody := map[string]interface{}{
		"desc":     payment.GetInvoice().GetDescription(),
		"amount":   payment.GetInvoice().GetAmount(),
		"order_id": payment.GetInvoice().GetUUID(),
		"phone":    "+989120000000",
		"mail":     "example@example.com",
		"name":     "example",
		"callback": idpayCallbackURL,
	}
	reqHeader := map[string]string{
		"X-API-KEY": idpayMerchant,
		"X-SANDBOX": "false",
	}
	respBody := map[string]interface{}{
		"id": "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIPurchaseURL, reqBody, reqHeader).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 internal server error",
			BodyProperty:       respJSON,
		}, nil).Once()

	err := payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrPurchaseFailed{Message: "500 internal server error"})
	s.Equal("500 internal server error", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", idpay.APIPurchaseURL, reqBody, reqHeader).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 500,
			StatusProperty:     "500 internal server error",
			BodyProperty:       respJSON,
		}, fmt.Errorf("example")).Once()

	err = payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.Equal("example", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	s.HTTPClient.
		On("Post", idpay.APIPurchaseURL, reqBody, reqHeader).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 201,
			StatusProperty:     "201 created",
			BodyProperty:       []byte("example"),
		}, nil).Once()

	err = payment.Purchase()
	s.NotNil(err)
	s.Error(err)
	s.HTTPClient.AssertExpectations(s.T())
}

func (s *GatewayIDPayTestSuite) TestVerifyFailed() {
	reqBody := idpay.VerifyRequest{
		RefID: "xxxx-xxxx-xxxx-xxxx-xxxx",
		ID:    "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	respBody := map[string]interface{}{
		"status":   1,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ := json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err := s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "پرداخت انجام نشده است"})
	s.Equal("پرداخت انجام نشده است", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
	//
	respBody = map[string]interface{}{
		"status":   2,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "پرداخت ناموفق بوده است"})
	s.Equal("پرداخت ناموفق بوده است", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
	//
	respBody = map[string]interface{}{
		"status":   3,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "خطا رخ داده است"})
	s.Equal("خطا رخ داده است", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
	//
	respBody = map[string]interface{}{
		"status":   4,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "بلوکه شده"})
	s.Equal("بلوکه شده", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
	//
	respBody = map[string]interface{}{
		"status":   5,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "برگشت به پرداخت کننده"})
	s.Equal("برگشت به پرداخت کننده", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
	//
	respBody = map[string]interface{}{
		"status":   6,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "برگشت خورده سیستمی"})
	s.Equal("برگشت خورده سیستمی", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
	//
	respBody = map[string]interface{}{
		"status":   7,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "انصراف از پرداخت"})
	s.Equal("انصراف از پرداخت", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
	//
	respBody = map[string]interface{}{
		"status":   8,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "به درگاه پرداخت منتقل شد"})
	s.Equal("به درگاه پرداخت منتقل شد", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
	//
	respBody = map[string]interface{}{
		"status":   10,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "در انتظار تایید پرداخت"})
	s.Equal("در انتظار تایید پرداخت", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
	//
	respBody = map[string]interface{}{
		"status":   101,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "پرداخت قبلا تایید شده است"})
	s.Equal("پرداخت قبلا تایید شده است", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
	//
	respBody = map[string]interface{}{
		"status":   200,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "به دریافت کننده واریز شد"})
	s.Equal("به دریافت کننده واریز شد", err.Error())
	s.HTTPClient.AssertExpectations(s.T())
	//
	respBody = map[string]interface{}{
		"status":   0,
		"track_id": "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)
	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
		}, nil).Once()

	_, err = s.Driver.Verify(&reqBody)
	s.NotNil(err)
	s.Error(err)
	s.ErrorIs(err, errors.ErrInvalidPayment{Message: "وضعیت نامشخص"})
	s.Equal("وضعیت نامشخص", err.Error())
	s.HTTPClient.AssertExpectations(s.T())

	respBody = map[string]interface{}{
		"error_message": "example",
		"status":        3,
		"track_id":      "string",
		"payment": map[string]interface{}{
			"card_no":        "string",
			"hashed_card_no": "string",
		},
	}
	respJSON, _ = json.Marshal(respBody)

	s.HTTPClient.
		On("Post", idpay.APIVerifyURL, &reqBody, map[string]string{
			"X-API-KEY": idpayMerchant,
			"X-SANDBOX": "false",
		}).
		Return(&mock.FakeResponse{
			StatusCodeProperty: 200,
			StatusProperty:     "200 ok",
			BodyProperty:       respJSON,
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
func TestGatewayIDPayTestSuite(t *testing.T) {
	suite.Run(t, new(GatewayIDPayTestSuite))
}
