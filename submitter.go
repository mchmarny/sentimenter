package sentimenter

import (
	"net/http"

	"cloud.google.com/go/logging"
)

// SubmitFunction represents the request submit functionality
func SubmitFunction(w http.ResponseWriter, r *http.Request) {

	once.Do(func() {
		configInitializer("submit")
	})

	defer logger.Flush()

	logger.Log(logging.Entry{
		HTTPRequest: &logging.HTTPRequest{
			Request: r,
		},
		Payload:  "SubmitFunction invoked",
		Severity: logging.Info,
	})

	term := r.URL.Query().Get("term")
	if term == "" {
		errLogger.Println("Nil term query parameter")
		http.Error(w, "The 'term' parameter is required", http.StatusInternalServerError)
		return
	}

	// create request for term
	infoLogger.Printf("Term: %s", term)
	job := newRequest(term)

	// save request
	if err := saveJob(job); err != nil {
		errLogger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*") //TODO: domain
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")

	w.Write(serializeOrFail(job))

}
