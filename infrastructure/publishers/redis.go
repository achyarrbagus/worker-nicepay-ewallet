package publishers

import (
	"fmt"
	"log"
	"worker-nicepay/infrastructure/configuration"
	"strconv"
)

// LogRedisClient simulates a Redis client for demonstration purposes
type LogRedisClient struct{}

var RDS *LogRedisClient

func InitializeRedis() {
	// Create a simulated Redis client
	RDS = &LogRedisClient{}

	redisAddr := configuration.AppConfig.RedisHost + ":" + strconv.Itoa(configuration.AppConfig.RedisPort)

	fmt.Println("=== Load Cache, Pub/Sub Redis ===")
	log.Printf("Simulated Redis connection established to %s (DB: %d)",
		redisAddr, configuration.AppConfig.RedisDatabase)
	fmt.Println("=== Load Cache, Pub/Sub Redis ===")
}
