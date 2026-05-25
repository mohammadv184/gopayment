package sepehr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mohammadv184/gopayment/v2"
	"github.com/mohammadv184/gopayment/v2/gateway"
	ihttp "github.com/mohammadv184/gopayment/v2/internal/transport/http"
)

// Gateway implements the Sepehr gateway.
type Gateway struct {
	baseURL    string
	terminalID string
	callback   string
	client     *http.Client
}

const (
	DefaultBaseURL = "https://sepehr.shaparak.ir"
)
const (
	EndpointGetToken = "Rest/V1/PeymentApi/GetToken"
	EndpointPay      = "Pay"
	EndpointVerify   = "Rest/V1/PeymentApi/Advice"
)

// New creates a new Sepehr driver
func New(terminalID, callback string) *Gateway {
	return &Gateway{
		terminalID: terminalID,
		callback:   callback,
		client:     ihttp.NewClient(),
		baseURL:    DefaultBaseURL,
	}
}
func (g *Gateway) WithBaseURL(url string) *Gateway {
	g.baseURL = url
	return g
}
func (g *Gateway) WithHTTPClient(client *http.Client) *Gateway {
	g.client = client
	return g
}

func (g *Gateway) Name() string { return gateway.Sepehr }

func (g *Gateway) Pay(ctx context.Context, p *gopayment.Payment) (*gopayment.Transaction, error) {
	irr, _ := gopayment.GetCurrency(gopayment.IRR)
	amount, err := p.Amount.ConvertTo(irr)
	if err != nil {
		return nil, err
	}

	reqBody := url.Values{
		"terminalID":  {g.terminalID},
		"callbackURL": {g.callback},
		"Amount":      {fmt.Sprintf("%d", amount.Value())},
		"InvoiceID":   {p.ID},
	}

	resp, err := g.client.PostForm(g.baseURL+"/"+EndpointGetToken, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Status      int    `json:"Status"`
		AccessToken string `json:"AccessToken"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if res.Status != 0 {
		return nil, fmt.Errorf("sepehr: status %d:%s", res.Status, statusToMessageMap[res.Status])
	}

	tx := gopayment.NewTransaction(p, res.AccessToken)
	tx.Data.Set("accessToken", res.AccessToken)

	return tx, nil
}

func (g *Gateway) RedirectURL(_ context.Context, tx *gopayment.Transaction) (string, error) {
	return fmt.Sprintf("%s?terminalID=%s&token=%s", g.baseURL+"/"+EndpointPay, g.terminalID, tx.RefID), nil
}

func (g *Gateway) Confirm(ctx context.Context, tx *gopayment.Transaction) (*gopayment.Receipt, error) {
	digitalReceipt := tx.Data.Get("digitalreceipt")
	if digitalReceipt == "" {
		return nil, fmt.Errorf("sepehr: missing digital receipt in transaction data")
	}

	reqBody := url.Values{
		"digitalreceipt": {digitalReceipt},
		"Tid":            {g.terminalID},
		"InvoiceID":      {tx.Payment.ID},
	}

	resp, err := g.client.PostForm(g.baseURL+"/"+EndpointVerify, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Status   string `json:"Status"`
		ReturnId string `json:"ReturnId"`
		Message  string `json:"Message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	// Status "OK" means verification is successful
	if res.Status != "OK" {
		return nil, fmt.Errorf("sepehr verify error: %s (status: %s)", res.Message, res.Status)
	}

	receipt := gopayment.NewReceipt(tx, g.Name())
	receipt.WithData("status", res.Status, "returnId", res.ReturnId)

	return receipt, nil
}

var statusToMessageMap = map[int]string{
	0:  "OK",
	-1: "Transaction not found",
	-2: "Unauthorized ip access",
	-3: "Total error",
	-4: "CallbackURL mismatch",
	-5: "Unauthorized ip access",
	-6: "Refund service is disabled",
	-7: "InvoiceID is not valid",
}
