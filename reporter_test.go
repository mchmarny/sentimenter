package sentimenter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReportFunction(t *testing.T) {

	// Save a DB
	configFunc()
	if config.err != nil {
		t.Errorf("Error on config: %v", config.err)
	}

	job := newRequest("test")
	err := saveJob(job)
	if err != nil {
		t.Errorf("Error on job save: %v", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("/?id=%s", job.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ReportFunction)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v", rr)
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	//TODO: Add test for counter incrementing

}
