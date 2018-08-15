package sentimenter

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// RequestStatus represents the sentiment request job status
type RequestStatus int

const (
	jobStatusDefault   string = "Default"
	jobStatusReceived  string = "Received"
	jobStatusProcessed string = "Processed"
	jobStatusFailed    string = "Failed"
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
	ID     string    `json:"id"`
	On     time.Time `json:"created_on"`
	Term   string    `json:"search_term"`
	Status string    `json:"status"`
}

func (r *SentimentRequest) String() string {
	return fmt.Sprintf("ID:%s, On:%s, Term:%s. Status:%s", r.ID, r.On, r.Term, r.Status)

}

func serializeOrFail(o interface{}) []byte {
	b, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
