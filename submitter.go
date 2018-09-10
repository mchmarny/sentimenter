package sentimenter

import (
	"net/http"
)

// SubmitFunction represents the request submit functionality
func SubmitFunction(w http.ResponseWriter, r *http.Request) {

	once.Do(func() {
		configInitializer()
	})

	defer logger.Flush()

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

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializeOrFail(job))

}
