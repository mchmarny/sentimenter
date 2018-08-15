package sentimenter

import (
	"encoding/base64"
	"testing"
)

func TestProcessorFunction(t *testing.T) {

	// Save a DB
	configFunc()
	if config.err != nil {
		t.Errorf("Error on config: %v", config.err)
	}

	job := newRequest("google")

	jobJSON := serializeOrFail(job)

	jobEncoded := base64.StdEncoding.EncodeToString(jobJSON)

	msg := PubSubMessage{
		Data: jobEncoded,
	}

	maxTweets = 1

	err := ProcessorFunction(ctx, msg)

	if err != nil {
		t.Errorf("Error processing job: %v", err)
	}

}
