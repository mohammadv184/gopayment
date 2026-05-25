package gopayment

import (
	"context"
)

// Gateway is the standard interface for all payment service providers
type Gateway interface {
	// Name returns the provider name (e.g., "Zibal", "Stripe")
	Name() string
	// Pay initiates a payment process
	Pay(ctx context.Context, p *Payment) (*Transaction, error)
	// Confirm confirms the payment after redirection
	Confirm(ctx context.Context, tx *Transaction) (*Receipt, error)
}

type GatewayRedirector interface {
	Gateway
	// RedirectURL returns the URL to redirect the user to for payment completion
	RedirectURL(ctx context.Context, tx *Transaction) (string, error)
}

// GatewayRefunder defines an optional interface for gateways that support refunds
type GatewayRefunder interface {
	Gateway
	// Refund initiates a refund process
	Refund(ctx context.Context, tx *Transaction) (bool, error)
}
