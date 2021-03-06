package sentimenter

import (
	"encoding/json"
	"testing"
)

func getTestJobMap() map[string]interface{} {
	// object to json
	j := []byte(`{
		"created_on":{"timestampValue":"2018-09-09T15:02:48.198087Z"},
		"id":{"stringValue":"tid-123"},
		"result":{"nullValue":null},
		"search_term":{"stringValue":"google"},
		"status":{"stringValue":"Received"},
		"status_url":{"stringValue":"https://test.com"}
	}`)

	// json to map
	var m map[string]interface{}
	json.Unmarshal(j, &m)
	return m
}

func TestConversions(t *testing.T) {

	// no config needed here

	m := getTestJobMap()

	job, err := eventMapToJob(m)
	if err != nil {
		t.Errorf("Error parsing: %v", err)
		return
	}

	if job.ID != "tid-123" || job.Term != "google" || job.Status != "Received" || job.URL != "https://test.com" {
		t.Errorf("Parsed job in an invalid state: %v", job)
		return
	}

}

func TestProcessorFunction(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestProcessorFunction")
	}

	configInitializer("test-process")

	m := getTestJobMap()

	job, err := eventMapToJob(m)
	if err != nil {
		t.Errorf("Error parsing job: %v", err)
		return
	}

	err = saveJob(job)
	if err != nil {
		t.Errorf("Error saving job: %v", err)
		return
	}

	// map into firestore event
	f := FirestoreEvent{
		Value: FirestoreValue{
			Fields: m,
		},
	}

	maxTweets = 1
	err = ProcessorFunction(ctx, f)
	if err != nil {
		t.Errorf("Error processing job: %v", err)
		return
	}

}
