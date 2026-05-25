package gopayment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Metadata represents a collection of additional information
type Metadata map[string]string

func (m Metadata) Set(key, value string) { m[key] = value }
func (m Metadata) Get(key string) string { return m[key] }
func (m Metadata) Delete(key string)     { delete(m, key) }
func (m Metadata) Has(key string) bool   { _, ok := m[key]; return ok }

// Payment represents a payment request.
type Payment struct {
	ID             string
	Gateway        string
	Amount         Amount
	Description    string
	IdempotencyKey string
	Data           Metadata
	Customer       *Customer
}

// Customer contains payer-related information
type Customer struct {
	ID          string
	Name        string
	Email       string
	PhoneNumber string
	IP          string
}

// NewPayment initializes a transaction with defaults
func NewPayment(gateway string, amount Amount) *Payment {
	return &Payment{
		ID:             uuid.New().String(),
		Amount:         amount,
		Gateway:        gateway,
		IdempotencyKey: uuid.New().String(),
		Data:           make(Metadata),
	}
}

// WithID overrides the default auto-generated transaction ID
func (p *Payment) WithID(id string) *Payment {
	p.ID = id
	return p
}

// WithDescription sets a human-readable Description for the payment
func (p *Payment) WithDescription(desc string) *Payment {
	p.Description = desc
	return p
}

// WithIdempotencyKey sets a custom key to prevent double charges
func (p *Payment) WithIdempotencyKey(key string) *Payment {
	p.IdempotencyKey = key
	return p
}

// WithData adds custom key-value pairs to the transaction
func (p *Payment) WithData(keyValue ...string) *Payment {
	for i := 0; i < len(keyValue)-1; i += 2 {
		p.Data.Set(keyValue[i], keyValue[i+1])
	}
	return p
}

// WithCustomer attaches Customer information to the transaction
func (p *Payment) WithCustomer(c Customer) *Payment {
	p.Customer = &c
	return p
}

func (p *Payment) Pay(ctx context.Context) (*Transaction, error) {
	if g, err := getGateway(p.Gateway); err == nil {
		return g.Pay(ctx, p)
	} else {
		return nil, err
	}
}

// Transaction wraps a Payment with its gateway-specific reference and status
type Transaction struct {
	Payment *Payment
	RefID   string
	Data    Metadata
}

func NewTransaction(payment *Payment, refID string) *Transaction {
	return &Transaction{
		Payment: payment,
		RefID:   refID,
		Data:    make(Metadata),
	}
}

func (t *Transaction) WithData(keyValue ...string) *Transaction {
	for i := 0; i < len(keyValue)-1; i += 2 {
		t.Data.Set(keyValue[i], keyValue[i+1])
	}
	return t
}

func (t *Transaction) RedirectURL(ctx context.Context) (string, error) {
	if g, err := getGateway(t.Payment.Gateway); err == nil {
		if gr, ok := g.(GatewayRedirector); ok {
			return gr.RedirectURL(ctx, t)
		}
		return "", ErrGatewayNotRedirector
	} else {
		return "", err
	}
}

// Encode returns the encoded transaction as a byte slice
func (t *Transaction) Encode() ([]byte, error) {
	return defaultEncoder.Encode(t)
}

// DecodeTransaction decodes a byte slice back into a Transaction struct
func DecodeTransaction(data []byte) (*Transaction, error) {
	return defaultEncoder.Decode(data)
}

// Receipt is the confirmation of a successful payment verification
type Receipt struct {
	PID     string
	RefID   string
	Gateway string
	Amount  Amount
	Date    time.Time
	Data    Metadata // Additional information returned by the gateway
}

func (r *Receipt) WithData(keyValue ...string) *Receipt {
	for i := 0; i < len(keyValue)-1; i += 2 {
		r.Data.Set(keyValue[i], keyValue[i+1])
	}
	return r
}

func (r *Receipt) GetData(key string) string { return r.Data.Get(key) }

func NewReceipt(t *Transaction, gateway string) *Receipt {
	return &Receipt{
		PID:     t.Payment.ID,
		RefID:   t.RefID,
		Gateway: gateway,
		Amount:  t.Payment.Amount,
		Date:    time.Now(),
		Data:    make(Metadata),
	}
}

var gatewayRegistry = make(map[string]Gateway)

func RegisterGateway(g Gateway, alias ...string) error {
	gk := g.Name()
	if len(alias) > 0 {
		gk = alias[0]
	}

	if _, ok := gatewayRegistry[gk]; ok {
		return ErrGatewayAlreadyRegistered
	}

	gatewayRegistry[gk] = g
	return nil
}

func getGateway(name string) (Gateway, error) {
	if g, ok := gatewayRegistry[name]; ok {
		return g, nil
	}
	return nil, ErrGatewayNotFound
}
