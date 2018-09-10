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
	r.URL = fmt.Sprintf(statusURLFormat, region, projectID, r.ID)
	return r
}

// SentimentRequest represents the sentiment request job
type SentimentRequest struct {
	ID     string           `json:"id" firestore:"id"`
	On     time.Time        `json:"created_on" firestore:"created_on"`
	Term   string           `json:"search_term" firestore:"search_term"`
	Status string           `json:"status" firestore:"status"`
	URL    string           `json:"status_url" firestore:"status_url"`
	Result *SentimentResult `json:"result" firestore:"result"`
}

// SentimentResult represents results of the job
type SentimentResult struct {
	Processed time.Time `json:"processed_on" firestore:"processed_on"`
	Tweets    int64     `json:"tweets" firestore:"tweets"`
	Positive  int64     `json:"positive" firestore:"positive"`
	Negative  int64     `json:"negative" firestore:"negative"`
	Score     float64   `json:"score" firestore:"score"`
}

// FirestoreValue is the payload of a FirestoreEvent event
type FirestoreValue struct {
	Fields interface{} `json:"fields"`
}

// FirestoreEvent is the Firestore document payload
type FirestoreEvent struct {
	OldValue FirestoreValue `json:"oldValue"`
	Value    FirestoreValue `json:"value"`
}
