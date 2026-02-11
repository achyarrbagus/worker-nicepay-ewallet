package database

import (
	"log"

	"worker-nicepay/infrastructure/configuration"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

var ElasticsearchClient *elasticsearch.TypedClient

func InitializeElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			configuration.AppConfig.ElasticsearchAddress,
		},
		Username: configuration.AppConfig.ElasticsearchUsername,
		Password: configuration.AppConfig.ElasticsearchPassword,
	}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Verify connection
	// In a real scenario, you might want to ping or check info,
	// but NewTypedClient doesn't strictly connect immediately until a request is made.
	// We can try to get info.
	res, err := es.Info().Do(nil)
	if err != nil {
		log.Printf("Error getting response from Elasticsearch: %s", err)
		// We might not want to fatal here if ES is optional, or we do if it's critical.
		// For now, let's log fatal to ensure we know it moved.
		// log.Fatal(err)
		// Actually, let's just log error to avoid crashing if ES is down during dev
		return
	}

	log.Printf("Elasticsearch initialized: %v", res.Version)

	ElasticsearchClient = es
}
