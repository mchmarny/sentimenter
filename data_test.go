package sentimenter

import "testing"

func TestJobData(t *testing.T) {

	configFunc()

	if config.err != nil {
		t.Errorf("Error on config: %v", config.err)
	}

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

	err = updateJobStatus(req.ID, jobStatusProcessing)

	if err != nil {
		t.Errorf("Error updating job status: %v", req)
	}
}