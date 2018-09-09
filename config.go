package sentimenter

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"

	"cloud.google.com/go/logging"
)

const (
	jobsCollectionName = "jobs"
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
	logger *logging.Logger
)

// configFunc sets the global configuration; it's overridden in tests.
var configFunc = getDefaultConfig

type configuration struct {
	client    *firestore.Client
	once      sync.Once
	err       error
	region    string
	projectID string
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
	config.projectID = projectID

	loggingClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		config.err = err
		return config.err
	}

	logger = loggingClient.Logger("sentimenter")

	region := os.Getenv("FUNCTION_REGION")
	if region == "" {
		// hack for testing
		region = "us-central1"
	}
	config.region = region

	// define client
	c, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Cannot create client: %v", err)
	}

	config.client = c

	return nil

}
