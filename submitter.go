package sentimenter

import (
	"net/http"

	"cloud.google.com/go/logging"
)

// SubmitFunction represents the request submit functionality
func SubmitFunction(w http.ResponseWriter, r *http.Request) {

	once.Do(func() {
		configInitializer()
	})

	defer logger.Flush()

	term := r.URL.Query().Get("term")
	if term == "" {
		logStringAll("Nil term query parameter")
		http.Error(w, "The 'term' parameter is required", http.StatusInternalServerError)
		return
	}

	// create request for term
	logger.StandardLogger(logging.Info).Printf("Term: %s", term)
	job := newRequest(term)

	// save request
	if err := saveJob(job); err != nil {
		logErrorAll(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializeOrFail(job))

}
