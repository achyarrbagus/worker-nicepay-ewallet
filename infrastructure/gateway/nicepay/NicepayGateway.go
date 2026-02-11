package nicepay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type NicepayGateway struct {
	URL    string
	APIKey string
	Client *resty.Client
}

func NewNicepayGateway(url string, apiKey string, timeout time.Duration) *NicepayGateway {
	return &NicepayGateway{
		URL:    url,
		APIKey: apiKey,
		Client: resty.New().SetTimeout(timeout),
	}
}

func (g *NicepayGateway) RequestPaymentLink(ctx context.Context, req RequestPaymentLinkDTO, url string) (ResponsePaymentLinkDTO, error) {

	var response ResponsePaymentLinkDTO
	queries, _ := json.Marshal(req)

	resp, err := g.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(queries).
		Post(url)

	reqHeaders, _ := json.Marshal(resp.Request.Header)
	respHeaders, _ := json.Marshal(resp.Header())

	response.RequestAPICallResult.RequestURL = url
	response.RequestAPICallResult.Method = resp.Request.RawRequest.Method
	response.RequestAPICallResult.RequestLatency = resp.Time().String()
	response.RequestAPICallResult.RequestBody = string(queries)
	response.RequestAPICallResult.ResponseBody = string(resp.Body())
	response.RequestAPICallResult.RequestHeaders = string(reqHeaders)
	response.RequestAPICallResult.ResponseHeaders = string(respHeaders)
	response.RequestAPICallResult.ResponseStatusCode = resp.StatusCode()

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "excedeed") {
			return ResponsePaymentLinkDTO{}, errors.New("timeout")
		} else {
			return ResponsePaymentLinkDTO{}, err
		}
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return ResponsePaymentLinkDTO{}, err
	}

	if resp.StatusCode() >= 400 {
		return ResponsePaymentLinkDTO{}, fmt.Errorf(response.Message)
	}

	return response, nil
}
