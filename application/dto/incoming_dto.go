package dto

type IncomingRequest struct {
	Country       string `json:"country"`
	ChannelCode   string `json:"channel_code"`
	CallbackUrl   string `json:"callback_url"`
	Description   string `json:"description"`
	PaymentMethod string `json:"payment_method"`
	Event         string `json:"event"`
	Email         string `json:"email"`
	Currency      string `json:"currency"`
}
