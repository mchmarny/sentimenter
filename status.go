package sentimenter

import (
	"net/http"
)

// StatusFunction represents the job status checker functionality
func StatusFunction(w http.ResponseWriter, r *http.Request) {

	once.Do(func() {
		configInitializer()
	})

	defer logger.Flush()

	id := r.URL.Query().Get("id")
	if id == "" {
		errLogger.Println("Nil job ID query parameter")
		http.Error(w, "Job `ID` parameter is required", http.StatusInternalServerError)
		return
	}

	// get jib by ID
	job, err := getJob(id)

	if err != nil {
		errLogger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// handle job not found for this ID
	if job.ID == "" {
		errLogger.Printf("Job not found: %s", id)
		http.Error(w, "Job not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializeOrFail(job))

}
