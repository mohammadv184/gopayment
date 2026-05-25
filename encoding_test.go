package gopayment

import (
	"reflect"
	"testing"
)

func TestEncoders(t *testing.T) {
	currency, _ := GetCurrency(USD)
	amount := NewAmount(1000, currency)
	payment := NewPayment("test_gateway", amount,
		WithDescription("Test Payment"),
		WithCustomer(Customer{
			ID:    "cust_123",
			Name:  "John Doe",
			Email: "john@example.com",
		}),
		WithData("key1", "value1"),
	)
	transaction := NewTransaction(payment, "ref_123")
	transaction.Data.Set("tkey", "tvalue")

	encoders := []struct {
		name    string
		encoder Encoder
	}{
		{"GobEncoder", NewGobEncoder()},
		{"JSONEncoder", NewJSONEncoder()},
	}

	for _, tc := range encoders {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := tc.encoder.Encode(transaction)
			if err != nil {
				t.Fatalf("failed to encode: %v", err)
			}

			decoded, err := tc.encoder.Decode(encoded)
			if err != nil {
				t.Fatalf("failed to decode: %v", err)
			}

			if !reflect.DeepEqual(transaction.Payment.ID, decoded.Payment.ID) {
				t.Errorf("expected ID %v, got %v", transaction.Payment.ID, decoded.Payment.ID)
			}
			if !reflect.DeepEqual(transaction.Payment.Gateway, decoded.Payment.Gateway) {
				t.Errorf("expected Gateway %v, got %v", transaction.Payment.Gateway, decoded.Payment.Gateway)
			}
			if transaction.Payment.Amount.Value() != decoded.Payment.Amount.Value() {
				t.Errorf("expected Amount Value %v, got %v", transaction.Payment.Amount.Value(), decoded.Payment.Amount.Value())
			}
			if transaction.Payment.Amount.Currency().Code != decoded.Payment.Amount.Currency().Code {
				t.Errorf("expected Currency Code %v, got %v", transaction.Payment.Amount.Currency().Code, decoded.Payment.Amount.Currency().Code)
			}
			if transaction.Payment.Description != decoded.Payment.Description {
				t.Errorf("expected Description %v, got %v", transaction.Payment.Description, decoded.Payment.Description)
			}
			if transaction.Payment.IdempotencyKey != decoded.Payment.IdempotencyKey {
				t.Errorf("expected IdempotencyKey %v, got %v", transaction.Payment.IdempotencyKey, decoded.Payment.IdempotencyKey)
			}
			if !reflect.DeepEqual(transaction.Payment.Data, decoded.Payment.Data) {
				t.Errorf("expected Payment Data %v, got %v", transaction.Payment.Data, decoded.Payment.Data)
			}
			if !reflect.DeepEqual(transaction.Data, decoded.Data) {
				t.Errorf("expected Transaction Data %v, got %v", transaction.Data, decoded.Data)
			}
			if !reflect.DeepEqual(transaction.Payment.Customer, decoded.Payment.Customer) {
				t.Errorf("expected Customer %v, got %v", transaction.Payment.Customer, decoded.Payment.Customer)
			}
			if transaction.RefID != decoded.RefID {
				t.Errorf("expected RefID %v, got %v", transaction.RefID, decoded.RefID)
			}
		})
	}
}
