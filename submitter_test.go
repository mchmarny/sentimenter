package sentimenter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitFunction(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestSubmitFunction")
	}

	req, err := http.NewRequest("GET", "/?term=test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SubmitFunction)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: %v", rr)
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

}
