package zibal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mohammadv184/gopayment/v2"
	gphttp "github.com/mohammadv184/gopayment/v2/internal/transport/http"
)

// Gateway implements the Zibal gateway
type Gateway struct {
	baseURL  string
	merchant string
	callback string
	client   *http.Client
}

const (
	DefaultBaseURL = "https://gateway.zibal.ir"
)

const (
	EndpointPurchase = "v1/request"
	EndpointVerify   = "v1/verify"
	EndpointPayment  = "start/%s"
)

// New creates a new Zibal gateway
func New(merchant, callback string) *Gateway {
	return &Gateway{
		baseURL:  DefaultBaseURL,
		merchant: merchant,
		callback: callback,
		client:   gphttp.NewClient(),
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

func (g *Gateway) Name() string { return "Zibal" }

func (g *Gateway) Pay(ctx context.Context, p *gopayment.Payment) (*gopayment.Transaction, error) {
	irr, _ := gopayment.GetCurrency(gopayment.IRR)
	amount, err := p.Amount.ConvertTo(irr)
	if err != nil {
		return nil, err
	}

	reqBody := map[string]interface{}{
		"merchant":    g.merchant,
		"amount":      amount.Value(),
		"callbackUrl": g.callback,
		"description": p.Description,
		"orderId":     p.ID,
	}

	if p.Customer != nil && p.Customer.PhoneNumber != "" {
		reqBody["mobile"] = p.Customer.PhoneNumber
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/"+EndpointPurchase, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)

	if result, ok := res["result"].(float64); !ok || result != 100 {
		return nil, fmt.Errorf("zibal error: %v", res["message"])
	}

	trackID := strconv.Itoa(int(res["trackId"].(float64)))
	tx := gopayment.NewTransaction(p, trackID)
	for k, v := range res {
		tx.Data.Set(k, fmt.Sprint(v))
	}

	return tx, nil
}

func (g *Gateway) RedirectURL(_ context.Context, tx *gopayment.Transaction) (string, error) {
	return fmt.Sprintf("%s/%s/%s", g.baseURL, EndpointPayment, tx.RefID), nil
}

func (g *Gateway) Confirm(ctx context.Context, tx *gopayment.Transaction) (*gopayment.Receipt, error) {
	irr, _ := gopayment.GetCurrency(gopayment.IRR)
	amount, _ := tx.Payment.Amount.ConvertTo(irr)

	reqBody := map[string]interface{}{
		"merchant": g.merchant,
		"trackId":  tx.RefID,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/"+EndpointVerify, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)

	if result, ok := res["result"].(float64); !ok || (result != 100 && result != 201) {
		return nil, fmt.Errorf("zibal error: %v", res["message"])
	}

	if resAmount, ok := res["amount"].(float64); ok && int64(resAmount) != amount.Value() {
		return nil, fmt.Errorf("zibal error: amount mismatch")
	}

	receipt := gopayment.NewReceipt(tx, g.Name())
	for k, v := range res {
		receipt.Data.Set(k, fmt.Sprint(v))
	}

	return receipt, nil
}
