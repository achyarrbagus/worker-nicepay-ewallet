package nicepay

import "worker-nicepay/infrastructure/gateway"

type RequestPaymentLinkDTO struct {
	CallbackURL string  `json:"url_callback"`
	ReturnURL   string  `json:"url_return"`
	MSISDN      string  `json:"msisdn"`
	Name        string  `json:"name"`
	Number      string  `json:"number"`
	Channel     string  `json:"channel"`
	Amount      float64 `json:"amount"`
	Email       string  `json:"email"`
	Description string  `json:"description"`
	IPAddress   string  `json:"ip_address"`
}

type ResponsePaymentLinkDTO struct {
	Error        bool        `json:"error"`
	StatusCode   int         `json:"status_code"`
	Message      string      `json:"message"`
	ErrorMessage interface{} `json:"error_message"`
	RedirectURL  string      `json:"redirect_url"`

	// APICall Result
	RequestAPICallResult gateway.RequestAPICallResult `json:"-"`
}

func (s *ResponsePaymentLinkDTO) GetAPICall() gateway.RequestAPICallResult {
	return s.RequestAPICallResult
}
