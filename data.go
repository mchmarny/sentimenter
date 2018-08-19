package sentimenter

import (
	"errors"
	"log"

	"cloud.google.com/go/spanner"
)

func saveJob(req *SentimentRequest) error {

	if config.db == nil {
		log.Fatal("DB not configured in saveJobs")
	}

	if req == nil {
		return errors.New("Nil parameter")
	}

	m := spanner.InsertOrUpdate("jobs",
		[]string{"id", "search_term", "created_on", "status"},
		[]interface{}{req.ID, req.Term, req.On, req.Status})

	_, err := config.db.Apply(ctx, []*spanner.Mutation{m})

	if err != nil {
		log.Fatalf("Error on job save: %v", err)
	}

	log.Printf("Saved Job:%s Term:%s ", req.ID, req.Term)

	return err

}

func saveResults(req *SentimentRequest) error {

	if config.db == nil {
		log.Fatal("DB not configured in saveJobs")
	}

	if req == nil {
		return errors.New("Nil parameter")
	}

	m := spanner.InsertOrUpdate("results",
		[]string{"id", "processed_on", "tweets", "positive", "negative", "score"},
		[]interface{}{req.ID, req.Result.Processed, req.Result.Tweets,
			req.Result.Positive, req.Result.Negative, req.Result.Score})

	_, err := config.db.Apply(ctx, []*spanner.Mutation{m})

	if err != nil {
		log.Fatalf("Error on result write: %v", err)
	}

	log.Printf("Saved Job Result:%s Score:%v", req.ID, req.Result.Score)

	return err

}

func updateJobStatus(id, status string) error {

	if config.db == nil {
		log.Fatal("DB not configured in saveJobs")
	}

	m := spanner.InsertOrUpdate("jobs",
		[]string{"id", "status"},
		[]interface{}{id, status})

	_, err := config.db.Apply(ctx, []*spanner.Mutation{m})

	if err != nil {
		log.Fatalf("Error on DB update: %v", err)
	}

	log.Printf("Updated Job:%s Status:%s ", id, status)

	return err

}

func getJob(id string) (req *SentimentRequest, err error) {

	if config.db == nil {
		log.Fatal("DB not configured")
	}

	if id == "" {
		return nil, errors.New("Nil parameter")
	}

	var result = &SentimentRequest{}
	row, err := config.db.Single().ReadRow(ctx, "jobs",
		spanner.Key{id}, []string{"id", "search_term", "created_on", "status"})

	if err != nil {
		log.Println(err)
		return result, err
	}

	err = row.Columns(&result.ID, &result.Term, &result.On, &result.Status)

	return result.setStatus(), err

}

func getResult(id string) (r *SentimentResult, err error) {

	if config.db == nil {
		log.Fatal("DB not configured")
	}

	if id == "" {
		return nil, errors.New("Nil parameter")
	}

	var result = &SentimentResult{}
	row, err := config.db.Single().ReadRow(ctx, "results", spanner.Key{id},
		[]string{"processed_on", "tweets", "positive", "negative", "score"})

	if err != nil {
		log.Println(err)
		return result, err
	}

	err = row.Columns(&result.Processed, &result.Tweets,
		&result.Positive, &result.Negative, &result.Score)

	return result, err

}
