package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"payment-airpay/domain/entities"
)

type ApidogGateway struct {
	URL    string
	Client *http.Client
}

func (g *ApidogGateway) Create(ctx context.Context, payload map[string]interface{}) (entities.Payment, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return entities.Payment{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.URL, bytes.NewBuffer(body))
	if err != nil {
		return entities.Payment{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.Client.Do(req)
	if err != nil {
		return entities.Payment{}, err
	}
	defer resp.Body.Close()

	var mockResp entities.Payment
	if err := json.NewDecoder(resp.Body).Decode(&mockResp); err != nil {
		return entities.Payment{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return entities.Payment{}, errors.New(resp.Status)
	}
	return mockResp, nil
}
