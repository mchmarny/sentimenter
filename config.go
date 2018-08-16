package sentimenter

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/spanner"
)

func init() {
	config = &configuration{}
	if config.Error() != nil {
		log.Printf("Error on init: %v", config.err)
	}
}

var (
	config *configuration
	ctx    context.Context
)

// configFunc sets the global configuration; it's overridden in tests.
var configFunc = getDefaultConfig

type configuration struct {
	topic  *pubsub.Topic
	client *pubsub.Client
	db     *spanner.Client
	once   sync.Once
	err    error
}

func (c *configuration) Error() error {
	return c.err
}

type envError struct {
	name string
}

func (e *envError) Error() string {
	return fmt.Sprintf("%s environment variable unset or missing", e.name)
}

func getDefaultConfig() error {

	ctx = context.Background()

	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		config.err = &envError{"GCP_PROJECT"}
		return config.err
	}

	topicName := os.Getenv("TOPIC_NAME")
	if topicName == "" {
		config.err = &envError{"TOPIC_NAME"}
		return config.err
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		config.err = &envError{"DB_PATH"}
		return config.err
	}

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		config.err = err
		return config.err
	}

	config.client = client
	config.topic = client.Topic(topicName)

	db, err := spanner.NewClient(ctx, dbPath)
	if err != nil {
		config.err = err
		return config.err
	}

	config.db = db

	return nil

}
