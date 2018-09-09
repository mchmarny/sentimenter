package sentimenter

import (
	"errors"
	"fmt"
)

func saveJob(job *SentimentRequest) error {

	if config.client == nil {
		return errors.New("Client not configured in saveJobs")
	}

	if job == nil {
		return errors.New("Nil parameter")
	}

	if job.ID == "" {
		return errors.New("Nil job ID")
	}

	_, err := config.client.Collection(jobsCollectionName).Doc(job.ID).Set(ctx, job)
	if err != nil {
		return fmt.Errorf("Error on job save: %v", err)
	}

	return nil

}

func getJob(id string) (req *SentimentRequest, err error) {

	if config.client == nil {
		return nil, errors.New("Client not configured in getJob")
	}

	if id == "" {
		return nil, errors.New("Nil job ID parameter")
	}

	d, err := config.client.Collection(jobsCollectionName).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	if d == nil {
		return nil, fmt.Errorf("No doc for ID: %s", id)
	}

	c := newRequest("")

	if e := d.DataTo(&c); err != nil {
		return nil, fmt.Errorf("Error converting doc to job: %v", e)
	}

	return c, nil

}
