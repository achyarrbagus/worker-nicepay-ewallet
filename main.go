package main

import (
	"log"
	"strconv"

	"worker-nicepay/infrastructure/configuration"
	"worker-nicepay/infrastructure/database"
	"worker-nicepay/infrastructure/middleware"
	"worker-nicepay/infrastructure/publishers"
	"worker-nicepay/infrastructure/queue"
	"worker-nicepay/infrastructure/workers"

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

	// Initialize YugabyteDB
	log.Println("Initializing YugabyteDB...")
	database.InitializeYugabyteDB()
	log.Println("YugabyteDB initialized")

	// Initialize Elasticsearch
	log.Println("Initializing Elasticsearch...")
	database.InitializeElasticsearch()
	log.Println("Elasticsearch initialized")

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
	// tambhkan middleware incoming dsini
	m := middleware.Middlewares{}
	app.Use(m.Incoming())

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	// Register routes
	app.Post("/payment/nicepay", workers.PaymentHandler)
	app.Post("/payment/nicepay/async", workers.EnqueueHandler)
	app.Get("/jobs/status", workers.StatusHandler)

	// Start server
	port := strconv.Itoa(configuration.AppConfig.ApplicationPort)
	if port == "0" || port == "" {
		port = "8080" // default port
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
