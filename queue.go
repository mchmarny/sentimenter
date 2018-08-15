package sentimenter

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/pubsub"
)

func publishJob(ctx context.Context, data []byte) error {

	if data == nil {
		log.Println("Nil data")
		return errors.New("Job data required")
	}

	result := config.topic.Publish(ctx, &pubsub.Message{Data: data})
	id, err := result.Get(ctx)
	if err != nil {
		log.Printf("Error while publishing message: %v:%v", err, id)
		return err
	}

	log.Printf("Published: %v", id)

	return nil
}
