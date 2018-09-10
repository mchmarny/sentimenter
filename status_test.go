package sentimenter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusFunction(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestStatusFunction")
	}

	configInitializer()

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
	handler := http.HandlerFunc(StatusFunction)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v", rr)
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	//TODO: Add test for counter incrementing

}
