package sentimenter

import (
	"time"
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
	return &SentimentRequest{
		ID:     getNewID(),
		On:     time.Now(),
		Term:   term,
		Status: jobStatusReceived,
	}
}

// SentimentRequest represents the sentiment request job
type SentimentRequest struct {
	ID     string           `json:"id"`
	On     time.Time        `json:"created_on"`
	Term   string           `json:"search_term"`
	Status string           `json:"status"`
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
