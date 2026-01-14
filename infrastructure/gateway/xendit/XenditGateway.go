package xendit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"payment-airpay/domain/entities"
)

type APIError struct {
	StatusCode int
	Body       []byte
}

func (e *APIError) Error() string {
	if e == nil {
		return "xendit api error"
	}
	return "xendit api error"
}

type XenditGateway struct {
	URL    string
	APIKey string
	Client *http.Client
}

func NewXenditGateway(url string, apiKey string, timeout time.Duration) *XenditGateway {
	return &XenditGateway{
		URL:    url,
		APIKey: apiKey,
		Client: &http.Client{Timeout: timeout},
	}
}

func (g *XenditGateway) Create(ctx context.Context, payload map[string]interface{}) (entities.Payment, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return entities.Payment{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.URL, bytes.NewBuffer(body))
	if err != nil {
		return entities.Payment{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-version", "2024-11-11")
	if g.APIKey != "" {
		req.SetBasicAuth(g.APIKey, "")
	}
	resp, err := g.Client.Do(req)
	if err != nil {
		return entities.Payment{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return entities.Payment{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Preserve the exact JSON error response from Xendit
		return entities.Payment{}, &APIError{StatusCode: resp.StatusCode, Body: respBytes}
	}

	var okResp entities.Payment
	if err := json.Unmarshal(respBytes, &okResp); err != nil {
		return entities.Payment{}, errors.New("failed to decode xendit response")
	}
	return okResp, nil
}
