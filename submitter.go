package sentimenter

import (
	"log"
	"net/http"

	"cloud.google.com/go/logging"
)

// SubmitFunction represents the request submit functionality
func SubmitFunction(w http.ResponseWriter, r *http.Request) {

	config.once.Do(func() { configFunc() })

	defer logger.Flush()

	if config.Error() != nil {
		log.Println(config.Error())
		logger.StandardLogger(logging.Error).Println(config.Error())
		http.Error(w, config.Error().Error(), http.StatusInternalServerError)
		return
	}

	term := r.URL.Query().Get("term")
	if term == "" {
		log.Println("Nil term")
		logger.StandardLogger(logging.Error).Println("Term parameter required")
		http.Error(w, "The 'term' parameter is required", http.StatusInternalServerError)
		return
	}

	// create request for term
	logger.StandardLogger(logging.Info).Printf("Term: %s", term)
	job := newRequest(term)

	// save request
	if err := saveJob(job); err != nil {
		log.Println(err)
		logger.StandardLogger(logging.Error).Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// publish job
	if err := publishJob(ctx, serializeOrFail(job)); err != nil {
		log.Println(err)
		logger.StandardLogger(logging.Error).Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializeOrFail(job))

}
