package sentimenter

import (
	"log"
	"net/http"
)

// ReportFunction represents the job status checker functionality
func ReportFunction(w http.ResponseWriter, r *http.Request) {

	config.once.Do(func() { configFunc() })

	if config.Error() != nil {
		log.Println(config.Error())
		http.Error(w, config.Error().Error(), http.StatusInternalServerError)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("Nil id")
		http.Error(w, "Job `ID` parameter is required", http.StatusInternalServerError)
		return
	}

	// get jib by ID
	job, err := getJob(id)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// save request
	if job.ID == "" {
		log.Printf("Job not found: %s", id)
		http.Error(w, "Job not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializeOrFail(job))

}
