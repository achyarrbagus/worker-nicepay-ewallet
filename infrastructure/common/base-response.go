package common

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HttpResponse struct {
	logger *zap.Logger
}

func NewHttpResponse(logger *zap.Logger) *HttpResponse {
	return &HttpResponse{
		logger: logger,
	}
}

type MetaData struct {
	Page      int `json:"page"`
	TotalPage int `json:"total_pages"`
	TotalRows int `json:"total_rows"`
	Limit     int `json:"limit"`
}

type Response struct {
	Status  int         `json:"status"`
	Error   bool        `json:"error"`
	TrxId   string      `json:"trx_id,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *MetaData   `json:"meta_data,omitempty"`
}

func BuildErrorResponse(message string, status int, err error, trxId string) Response {

	res := Response{
		Error:   true,
		TrxId:   trxId,
		Status:  status,
		Message: message,
		Data:    err.Error(),
	}
	return res
}

func BuildSuccessResponse(message string, status int, data interface{}, trxId string) Response {
	res := Response{
		Error:   false,
		Status:  status,
		Message: message,
		TrxId:   trxId,
		Data:    data,
	}

	return res
}

func SuccessResponse(c *fiber.Ctx, status int, message string, data interface{}, trxId string) error {
	if message == "" {
		message = http.StatusText(http.StatusOK)
	}

	resp := BuildSuccessResponse(message, status, data, trxId)

	c.Locals("response", resp)

	return c.Status(status).JSON(resp)
}

func ErrorResponse(c *fiber.Ctx, status int, message string, err error, request interface{}, trxId string) error {

	// Logger usage would need a global logger or be passed in, or we skip logging in this static helper
	// and rely on middleware or caller.
	// For now, removing logger dependency from this static helper to match user request of simple call.

	resp := BuildErrorResponse(message, status, err, trxId)
	c.Locals("response", resp)

	return c.Status(status).JSON(resp)
}
