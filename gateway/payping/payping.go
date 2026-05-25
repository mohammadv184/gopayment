package payping

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mohammadv184/gopayment/v2"
	"github.com/mohammadv184/gopayment/v2/gateway"
	ihttp "github.com/mohammadv184/gopayment/v2/internal/transport/http"
)

// Gateway implements the PayPing gateway.
type Gateway struct {
	baseURL  string
	token    string
	callback string
	client   *http.Client
}

const (
	DefaultBaseURL = "https://api.payping.ir"
)

const (
	EndpointPay    = "v3/pay"
	EndpointVerify = "v3/pay/verify"
	EndpointStart  = "v3/pay/start"
	EndpointRefund = "v3/pay/reverse"
)

// New creates a new PayPing gateway.
func New(token, callback string) *Gateway {
	return &Gateway{
		baseURL:  DefaultBaseURL,
		token:    token,
		callback: callback,
		client:   ihttp.NewClient(),
	}
}

// WithBaseURL sets a custom base URL for the gateway.
func (g *Gateway) WithBaseURL(url string) *Gateway {
	g.baseURL = url
	return g
}

// WithHTTPClient sets a custom http client for the gateway.
func (g *Gateway) WithHTTPClient(client *http.Client) *Gateway {
	g.client = client
	return g
}

// Name returns the gateway name.
func (g *Gateway) Name() string { return gateway.PayPing }

// Pay initiates a payment request.
func (g *Gateway) Pay(ctx context.Context, p *gopayment.Payment) (*gopayment.Transaction, error) {
	irt, _ := gopayment.GetCurrency(gopayment.IRT)
	amount, err := p.Amount.ConvertTo(irt)
	if err != nil {
		return nil, err
	}

	reqBody := map[string]interface{}{
		"amount":      amount.Value(),
		"returnUrl":   g.callback,
		"description": p.Description,
		"clientRefId": p.ID,
	}

	if p.Customer != nil {
		if p.Customer.Name != "" {
			reqBody["payerName"] = p.Customer.Name
		}
		if p.Customer.PhoneNumber != "" {
			reqBody["payerIdentity"] = p.Customer.PhoneNumber
		} else if p.Customer.Email != "" {
			reqBody["payerIdentity"] = p.Customer.Email
		}
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/"+EndpointPay, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errRes map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errRes)
		return nil, fmt.Errorf("payping error: %v", errRes)
	}

	var res struct {
		PaymentCode string `json:"paymentCode"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	tx := gopayment.NewTransaction(p, res.PaymentCode)

	return tx, nil
}

// RedirectURL returns the URL to redirect the user to.
func (g *Gateway) RedirectURL(_ context.Context, tx *gopayment.Transaction) (string, error) {
	return fmt.Sprintf("%s/%s/%s", g.baseURL, EndpointStart, tx.RefID), nil
}

// Confirm verifies the payment.
func (g *Gateway) Confirm(ctx context.Context, tx *gopayment.Transaction) (*gopayment.Receipt, error) {
	refID := tx.Data.Get("paymentRefId")
	if refID == "" {
		return nil, fmt.Errorf("payping: missing paymentRefId in transaction data")
	}

	irt, _ := gopayment.GetCurrency(gopayment.IRT)
	amount, _ := tx.Payment.Amount.ConvertTo(irt)

	reqBody := map[string]interface{}{
		"paymentCode":  tx.RefID,
		"paymentRefId": refID,
		"amount":       amount.Value(),
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/"+EndpointVerify, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errRes map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errRes)
		return nil, fmt.Errorf("payping verify error: %v", errRes)
	}

	receipt := gopayment.NewReceipt(tx, g.Name())
	receipt.WithData("refId", refID)

	return receipt, nil
}

// Refund initiates a refund process.
func (g *Gateway) Refund(ctx context.Context, tx *gopayment.Transaction) (bool, error) {
	refID := tx.Data.Get("paymentRefId")
	if refID == "" {
		return false, fmt.Errorf("payping: missing paymentRefId in transaction data")
	}

	reqBody := map[string]interface{}{
		"paymentRefId": refID,
		"paymentCode":  tx.RefID,
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("%s/%s", g.baseURL, EndpointRefund)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errRes map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errRes)
		return false, fmt.Errorf("payping refund error: %v", errRes)
	}

	return true, nil
}
