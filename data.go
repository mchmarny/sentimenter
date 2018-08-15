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
		log.Fatalf("Error on DB write: %v", err)
	}

	log.Printf("Saved Job:%s Term:%s ", req.ID, req.Term)

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

	return result, err

}
