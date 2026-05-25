package gopayment

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
)

// Encoder defines the interface for encoding and decoding transactions
type Encoder interface {
	Encode(t *Transaction) ([]byte, error)
	Decode(data []byte) (*Transaction, error)
}

// GobEncoder implements Encoder using gob and base64
type GobEncoder struct{}

// NewGobEncoder creates a new GobEncoder
func NewGobEncoder() *GobEncoder {
	return &GobEncoder{}
}

// Encode returns the encoded transaction as a byte slice
func (e *GobEncoder) Encode(t *Transaction) ([]byte, error) {
	gout := &bytes.Buffer{}
	enc := gob.NewEncoder(gout)
	if err := enc.Encode(mapTransactionToView(t)); err != nil {
		return nil, err
	}

	out := make([]byte, base64.StdEncoding.EncodedLen(gout.Len()))
	base64.StdEncoding.Encode(out, gout.Bytes())

	return out, nil
}

// Decode decodes a byte slice back into a Transaction struct
func (e *GobEncoder) Decode(data []byte) (*Transaction, error) {
	bin := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(bin, data)
	if err != nil {
		return nil, err
	}

	gin := bytes.NewBuffer(bin[:n])
	dec := gob.NewDecoder(gin)

	var tv transactionView
	if err := dec.Decode(&tv); err != nil {
		return nil, err
	}

	return mapViewToTransaction(&tv)
}

// JSONEncoder implements Encoder using json
type JSONEncoder struct{}

// NewJSONEncoder creates a new JSONEncoder
func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

// Encode returns the encoded transaction as a byte slice
func (e *JSONEncoder) Encode(t *Transaction) ([]byte, error) {
	return json.Marshal(mapTransactionToView(t))
}

// Decode decodes a byte slice back into a Transaction struct
func (e *JSONEncoder) Decode(data []byte) (*Transaction, error) {
	var tv transactionView
	if err := json.Unmarshal(data, &tv); err != nil {
		return nil, err
	}

	return mapViewToTransaction(&tv)
}

var defaultEncoder Encoder = NewJSONEncoder()

// SetDefaultEncoder sets the default encoder used by Transaction.Encode and DecodeTransaction
func SetDefaultEncoder(e Encoder) {
	defaultEncoder = e
}

type transactionView struct {
	PID     string `json:"pid"`
	Gateway string `json:"gateway"`
	Amount  struct {
		Val  int64  `json:"val"`
		Curr string `json:"curr"`
	} `json:"amount"`
	Description    string            `json:"description"`
	IdempotencyKey string            `json:"idempotency_key"`
	Data           map[string]string `json:"data"`
	TData          map[string]string `json:"t_data"`
	Customer       *struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		IP          string `json:"ip"`
	} `json:"customer,omitempty"`
	TRefID string `json:"tref_id"`
}

func mapTransactionToView(t *Transaction) *transactionView {
	tv := &transactionView{
		PID:     t.Payment.ID,
		Gateway: t.Payment.Gateway,
		Amount: struct {
			Val  int64  `json:"val"`
			Curr string `json:"curr"`
		}{
			Val:  t.Payment.Amount.Value(),
			Curr: t.Payment.Amount.Currency().Code.String(),
		},
		Description:    t.Payment.Description,
		IdempotencyKey: t.Payment.IdempotencyKey,
		Data:           t.Payment.Data,
		TData:          t.Data,
		TRefID:         t.RefID,
	}

	if t.Payment.Customer != nil {
		tv.Customer = &struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number"`
			IP          string `json:"ip"`
		}{
			ID:          t.Payment.Customer.ID,
			Name:        t.Payment.Customer.Name,
			Email:       t.Payment.Customer.Email,
			IP:          t.Payment.Customer.IP,
			PhoneNumber: t.Payment.Customer.PhoneNumber,
		}
	}

	return tv
}

func mapViewToTransaction(t *transactionView) (*Transaction, error) {
	curr, err := GetCurrency(CurrencyCode(t.Amount.Curr))
	if err != nil {
		return nil, err
	}

	pAmount := NewAmount(t.Amount.Val, curr)
	p := NewPayment(t.Gateway, pAmount).WithID(t.PID).WithIdempotencyKey(t.IdempotencyKey)
	p.Data = t.Data
	p.Description = t.Description
	if t.Customer != nil {
		p.Customer = &Customer{
			ID:          t.Customer.ID,
			Name:        t.Customer.Name,
			Email:       t.Customer.Email,
			PhoneNumber: t.Customer.PhoneNumber,
			IP:          t.Customer.IP,
		}
	}

	transaction := NewTransaction(p, t.TRefID)
	transaction.Data = t.TData
	return transaction, nil
}
