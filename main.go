package main

import (
	"log"
	"strconv"

	"payment-airpay/infrastructure/configuration"
	"payment-airpay/infrastructure/database"
	"payment-airpay/infrastructure/publishers"
	"payment-airpay/infrastructure/queue"
	"payment-airpay/infrastructure/workers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	defer func() {
		queue.CloseRabbitMQ()
	}()

	log.Println("Xendit Worker is starting...")

	// Initialize configurations
	log.Println("Initializing configuration...")
	configuration.InitializeAppConfig()
	log.Println("Configuration initialized")

	log.Println("Initializing YugabyteDB...")
	database.InitializeYugabyteDB()
	log.Println("YugabyteDB initialized")

	// Initialize RabbitMQ
	log.Println("Initializing RabbitMQ...")
	queue.InitializeRabbitMQ()
	log.Println("RabbitMQ initialized")

	// Initialize Redis for publishing
	log.Println("Initializing publishers...")
	publishers.InitializeRedis()
	log.Println("Publishers initialized")

	// Initialize worker
	log.Println("Initializing worker...")
	workers.InitializePaymentXenditTaskWorker()
	log.Println("Worker initialized")

	// Initialize fiber app
	app := fiber.New()

	// Register routes
	app.Post("/payment/xendit", workers.PaymentHandler)
	app.Post("/payment/xendit/async", workers.EnqueueHandler)
	app.Get("/jobs/status", workers.StatusHandler)

	// Start server
	port := strconv.Itoa(configuration.AppConfig.ApplicationPort)
	if port == "0" || port == "" {
		port = "8080" // default port
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
