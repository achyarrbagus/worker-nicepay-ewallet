package entities

import (
	"encoding/json"
	"strconv"
)

type Payment struct {
	BusinessID       string                 `json:"business_id"`
	ReferenceID      string                 `json:"reference_id"`
	PaymentRequestID string                 `json:"payment_request_id"`
	Type             string                 `json:"type"`
	Country          string                 `json:"country"`
	Currency         string                 `json:"currency"`
	RequestAmount    float64                `json:"request_amount"`
	CaptureMethod    string                 `json:"capture_method"`
	ChannelCode      string                 `json:"channel_code"`
	ChannelProps     map[string]interface{} `json:"channel_properties"`
	Actions          []PaymentAction        `json:"actions"`
	Status           PaymentStatus          `json:"status"`
	Description      string                 `json:"description"`
	Metadata         map[string]interface{} `json:"metadata"`
	Created          string                 `json:"created"`
	Updated          string                 `json:"updated"`
}

type PaymentAction struct {
	Type       string `json:"type"`
	Value      string `json:"value"`
	Descriptor string `json:"descriptor"`
}

type PaymentStatus string

func (s *PaymentStatus) UnmarshalJSON(b []byte) error {
	// Accept both JSON string ("REQUIRES_ACTION") and JSON number (200)
	var asString string
	if err := json.Unmarshal(b, &asString); err == nil {
		*s = PaymentStatus(asString)
		return nil
	}

	var asNumber float64
	if err := json.Unmarshal(b, &asNumber); err == nil {
		*s = PaymentStatus(strconv.FormatInt(int64(asNumber), 10))
		return nil
	}

	// Fallback: keep raw
	*s = PaymentStatus(string(b))
	return nil
}
