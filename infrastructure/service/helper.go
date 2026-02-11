package service

import (
	"context"
	"fmt"
	"time"

	entity "worker-nicepay/domain/entities"
	"worker-nicepay/infrastructure/database"
	pkg "worker-nicepay/infrastructure/gateway"
)

func formatErrorToString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func SaveAPICall(
	ctx context.Context,
	param pkg.APICall,
	merchant string,
	err error,
	service string,
	track string,
	msisdn string,
	webtype string,
	transactionId string,
) {
	resp := param.GetAPICall()
	data := entity.ApiCall{
		ID:             0,
		CreatedAt:      time.Now(),
		Track:          track,
		Service:        service,
		Webtype:        webtype,
		Merchant:       merchant,
		Msisdn:         msisdn,
		URL:            resp.RequestURL,
		Method:         resp.Method,
		RequestQuery:   resp.RequestQuery,
		RequestBody:    resp.RequestBody,
		ResponseBody:   resp.ResponseBody,
		StatusCode:     resp.ResponseStatusCode,
		RequestHeader:  resp.RequestHeaders,
		ResponseHeader: resp.ResponseHeaders,
		Latency:        resp.RequestLatency,
		Error:          formatErrorToString(err),
		TransactionID:  transactionId,
	}

	if database.ElasticsearchClient != nil {
		_, err := database.ElasticsearchClient.Index("api_call_logs").
			Request(&data).
			Do(context.Background())
		if err != nil {
			fmt.Printf("Failed to save API call to Elasticsearch: %v\n", err)
		}
	} else {
		fmt.Println("ElasticsearchClient is not initialized")
	}
}
