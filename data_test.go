package sentimenter

import (
	"testing"
)

func TestJobData(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestJobData")
	}

	configInitializer("test-data")

	termReq := newRequest("test")

	err := saveJob(termReq)

	if err != nil {
		t.Errorf("Error on job save: %v", err)
	}

	req, err := getJob(termReq.ID)

	if err != nil {
		t.Errorf("Error on job read: %v", err)
	}

	if req.ID != termReq.ID {
		t.Errorf("Got invalid job: %v", req)
	}

}
