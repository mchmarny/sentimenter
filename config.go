package sentimenter

import (
	"context"
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
	infoLogger        *log.Logger
	errLogger         *log.Logger
	db                *firestore.Client
	once              sync.Once
	region            string
	projectID         string
	configValid       bool
	configInitializer = defaultConfigInitializer
)

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
	infoLogger = logger.StandardLogger(logging.Info)
	errLogger = logger.StandardLogger(logging.Error)

	// on the end
	configValid = true

}
