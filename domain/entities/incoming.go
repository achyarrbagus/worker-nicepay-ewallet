package entities

import "time"

type Incoming struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt     time.Time `json:"created_at"`
	Track         string    `json:"track"`
	Service       string    `json:"service"`
	Webtype       string    `json:"webtype"`
	Path          string    `json:"path"`
	Merchant      string    `json:"merchant"`
	IP            string    `json:"ip"`
	Method        string    `json:"method"`
	RequestQuery  string    `json:"request_query"`
	RequestHeader string    `json:"request_header"`
	RequestBody   string    `json:"request_body"`
	ResponseBody  string    `json:"response_body"`
	TransactionID string    `json:"transaction_id"`
	StatusCode    int       `json:"status_code"`
	Latency       string    `json:"latency"`
	UserAgent     string    `json:"user_agent"`
	Device        string    `json:"device"`
	Browser       string    `json:"browser"`
	Callback      *string   `json:"callback"`
	Country       string    `json:"country"`
	ChannelCode   string    `json:"channel_code"`
	CallbackUrl   string    `json:"callback_url"`
	Description   string    `json:"description"`
	PaymentMethod string    `json:"payment_method"`
	Event         string    `json:"event"`
	Email         string    `json:"email"`
	Curency       string    `json:"currency"`
	Save          bool      `json:"save" gorm:"-"` // Ignore in SQL DB if only for logic or Elastic
}
