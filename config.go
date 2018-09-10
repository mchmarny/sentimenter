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
	defaultJobsCollectionName = "jobs"
	processName               = "sentimenter"
)

var (
	ctx               context.Context
	logger            *logging.Logger
	skipRemoteLogging bool
	db                *firestore.Client
	once              sync.Once
	region            string
	projectID         string
	configValid       bool
	configInitializer = defaultConfigInitializer
)

func logAll(v string, args ...interface{}) {
	logStringAll(fmt.Sprintf(v, args...))
}

func logStringAll(v string) {
	log.Println(v)
	if logger == nil || !skipRemoteLogging {
		logger.StandardLogger(logging.Info).Printf(v)
	}
}

func logErrorAll(err error) error {
	log.Println(err)
	if logger == nil || !skipRemoteLogging {
		logger.StandardLogger(logging.Error).Println(err)
	}

	return err
}

func defaultConfigInitializer() {

	ctx = context.Background()

	projectID = os.Getenv("GCP_PROJECT")
	if projectID == "" {
		log.Fatalln("GCP_PROJECT environment variable not set")
	}

	region = os.Getenv("FUNCTION_REGION")
	if region == "" {
		log.Printf("FUNCTION_REGION not set, using default: us-central1")
		region = "us-central1"
	}

	dbClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error while creating Firestore client: %v", err)
	}
	db = dbClient

	logClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error while configuring logger: %v", err)
	}
	logger = logClient.Logger(processName)

	// on the end
	configValid = true

}
