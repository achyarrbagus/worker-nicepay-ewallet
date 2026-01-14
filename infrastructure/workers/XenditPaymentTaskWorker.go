package workers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"payment-airpay/infrastructure/dependencies"
	"payment-airpay/infrastructure/gateway/xendit"
	"time"

	"github.com/gofiber/fiber/v2"
)

// JobStatus represents the current status of a job
type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusDone       JobStatus = "done"
	StatusError      JobStatus = "error"
)

// JobResult contains the result and status of a job
type JobResult struct {
	ID      string          `json:"id"`
	Status  JobStatus       `json:"status"`
	Message string          `json:"message,omitempty"`
	Data    interface{}     `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
	Raw     json.RawMessage `json:"-"`
	Code    int             `json:"-"`
}

type Metadata struct {
	Merchant       string `json:"merchant"`
	PaymentGateway string `json:"payment_gateway"`
}

type ChannelProperties struct {
	DisplayName string `json:"display_name"`
}

type CreatePaymentVirtualAccount struct {
	ReferenceID   string             `json:"reference_id"`
	Type          string             `json:"type"`
	Country       string             `json:"country"`
	Currency      string             `json:"currency"`
	RequestAmount float64            `json:"request_amount"`
	Metadata      *Metadata          `json:"metadata"`
	ChannelCode   string             `json:"channel_code"`
	ChannelProps  *ChannelProperties `json:"channel_properties"`
}

type CreatePaymentRequest struct {
	ReferenceID   string                 `json:"reference_id"`
	Type          string                 `json:"type"`
	Country       string                 `json:"country"`
	Currency      string                 `json:"currency"`
	RequestAmount float64                `json:"request_amount"`
	ChannelCode   string                 `json:"channel_code"`
	ChannelProps  *ChannelProperties     `json:"channel_properties"`
	Description   string                 `json:"description"`
	Metadata      map[string]interface{} `json:"metadata"`
}

func (r *CreatePaymentRequest) ToPayloadMap() map[string]interface{} {
	payload := map[string]interface{}{}

	// Prefer new schema if present
	if r.ReferenceID != "" {
		payload["reference_id"] = r.ReferenceID
	}
	if r.Type != "" {
		payload["type"] = r.Type
	}
	if r.Country != "" {
		payload["country"] = r.Country
	}
	if r.Currency != "" {
		payload["currency"] = r.Currency
	}
	if r.RequestAmount != 0 {
		payload["request_amount"] = r.RequestAmount
	}
	if r.ChannelCode != "" {
		payload["channel_code"] = r.ChannelCode
	}
	if r.ChannelProps != nil {
		payload["channel_properties"] = map[string]interface{}{
			"display_name": r.ChannelProps.DisplayName,
		}
	}
	if r.Description != "" {
		payload["description"] = r.Description
	}
	if r.Metadata != nil {
		payload["metadata"] = map[string]interface{}{
			"merchant":        r.Metadata["merchant"],
			"payment_gateway": r.Metadata["payment_gateway"],
		}
	}

	return payload
}

// Worker manages job processing
type Worker struct {
	queue   chan map[string]interface{}
	results map[string]*JobResult
}

var workerInstance *Worker

func InitializePaymentXenditTaskWorker() {
	workerInstance = &Worker{
		queue:   make(chan map[string]interface{}, 100),
		results: make(map[string]*JobResult),
	}

	// Start the worker goroutine
	go workerInstance.processQueue()
}

func (w *Worker) processQueue() {
	for payload := range w.queue {
		// Get dependencies
		uc := dependencies.WireCreatePaymentService()

		// Process the payment
		ctx := context.Background()
		result, err := uc.Execute(ctx, payload)

		// Handle the result
		jobID := payload["job_id"].(string)
		if err != nil {
			log.Printf("Error processing job %s: %v", jobID, err)
			w.results[jobID] = &JobResult{
				ID:     jobID,
				Status: StatusError,
				Error:  err.Error(),
			}
		} else {
			log.Printf("Job %s completed: %v", jobID, result)
			w.results[jobID] = &JobResult{
				ID:      jobID,
				Status:  StatusDone,
				Message: string(result.Status),
				Data: map[string]string{
					"payment_request_id": result.PaymentRequestID,
				},
			}
		}
	}
}

// PaymentHandler handles payment requests
func PaymentHandler(c *fiber.Ctx) error {
	// Parse payload from request
	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	payload := req.ToPayloadMap()

	// Execute the payment use case directly
	uc := dependencies.WireCreatePaymentService()

	// Use context from the request
	result, err := uc.Execute(c.Context(), payload)
	if err != nil {
		var apiErr *xendit.APIError
		if errors.As(err, &apiErr) {
			return c.Status(apiErr.StatusCode).Type("json").Send(apiErr.Body)
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  fiber.StatusBadGateway,
			"error":   true,
			"message": err.Error(),
			"data":    fiber.Map{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// EnqueueHandler handles asynchronous job requests
func EnqueueHandler(c *fiber.Ctx) error {
	// Parse payload from request
	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	payload := req.ToPayloadMap()

	// Generate job ID
	jobID := generateJobID()
	payload["job_id"] = jobID

	// Create job result entry
	workerInstance.results[jobID] = &JobResult{
		ID:      jobID,
		Status:  StatusQueued,
		Message: "Job queued",
	}

	// Queue the job
	workerInstance.queue <- payload

	// Return job ID
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"job_id": jobID,
		"status": "queued",
	})
}

// StatusHandler gets job status
func StatusHandler(c *fiber.Ctx) error {
	jobID := c.Query("id")
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Job ID is required",
		})
	}

	result, exists := workerInstance.results[jobID]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Job not found",
		})
	}

	return c.JSON(result)
}

// Helper to generate a job ID
func generateJobID() string {
	return "job-" + time.Now().Format("20060102-150405-999999999")
}
