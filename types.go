package sentimenter

import (
	"fmt"
	"time"
)

const (
	//https://us-central1-s9-demo.cloudfunctions.net/sentimenter-status?id=c83eb7ff-a27b-4eb9-bd7c-70232bb5421d
	statusURLFormat = "https://%s-%s.cloudfunctions.net/sentimenter-status?id=%s"
)

// RequestStatus represents the sentiment request job status
type RequestStatus int

const (
	jobStatusDefault    string = "Default"
	jobStatusReceived   string = "Received"
	jobStatusProcessing string = "Processing"
	jobStatusProcessed  string = "Processed"
	jobStatusFailed     string = "Failed"
)

func newRequest(term string) *SentimentRequest {
	r := &SentimentRequest{
		ID:     getNewID(),
		On:     time.Now(),
		Term:   term,
		Status: jobStatusReceived,
	}
	return r.setStatus()
}

func (r *SentimentRequest) setStatus() *SentimentRequest {
	r.URL = fmt.Sprintf(statusURLFormat, config.region, config.projectID, r.ID)
	return r
}

// SentimentRequest represents the sentiment request job
type SentimentRequest struct {
	ID     string           `json:"id"`
	On     time.Time        `json:"created_on"`
	Term   string           `json:"search_term"`
	Status string           `json:"status"`
	URL    string           `json:"status_url"`
	Result *SentimentResult `json:"result"`
}

type SentimentResult struct {
	Processed time.Time `json:"processed_on"`
	Tweets    int64     `json:"tweets"`
	Positive  int64     `json:"positive"`
	Negative  int64     `json:"negative"`
	Score     float64   `json:"score"`
}

// PubSubMessage represents PubSub payload
type PubSubMessage struct {
	Data string `json:"data"`
}
