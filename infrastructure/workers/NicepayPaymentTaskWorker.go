package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"worker-nicepay/application/dto"
	"worker-nicepay/domain/entities"
	"worker-nicepay/infrastructure/common"
	"worker-nicepay/infrastructure/dependencies"

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

		// Convert payload map to DTO
		jsonBody, _ := json.Marshal(payload)
		var req dto.CreatePaymentRequest
		json.Unmarshal(jsonBody, &req)

		_, result, err := uc.Execute(ctx, req, entities.Incoming{
			IP: payload["ip"].(string),
		})

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
				Message: "Success", // Or extract something meaningful from result if it's not empty
				Data: map[string]string{
					"payment_request_id": result.PaymentRequestID,
				},
			}
		}
	}
}

// PaymentHandler handles payment requests
func PaymentHandler(c *fiber.Ctx) error {

	incoming, ok := c.Locals("incoming").(*entities.Incoming)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Incoming context missing"})
	}

	// Parse payload from request
	var req dto.CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return common.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request payload", err, req, incoming.TransactionID)
	}

	// Execute the payment use case directly
	uc := dependencies.WireCreatePaymentService()

	// Use context from the request
	_, result, err := uc.Execute(c.Context(), req, *incoming)
	if err != nil {
		return common.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), err, req, "")
	}

	return common.SuccessResponse(c, fiber.StatusOK, "Success", result, "")
}

// EnqueueHandler handles asynchronous job requests
func EnqueueHandler(c *fiber.Ctx) error {
	// Parse payload from request
	// Parse payload from request
	var req dto.CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	payload := req.ToPayloadMap() // Used for flexibility in queue or could use DTO in queue

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
