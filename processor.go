package sentimenter

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	lang "cloud.google.com/go/language/apiv1"
	"cloud.google.com/go/logging"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	langpb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

const (
	magnitudeThreshold     = 0.5
	positiveScoreThreshold = 0.5
	negativeScoreThreshold = -0.5
)

var (
	langClient *lang.Client
	maxTweets  = 100
)

// ProcessorFunction processes pubsub messages
func ProcessorFunction(ctx context.Context, m PubSubMessage) error {

	config.once.Do(func() {
		configFunc()

		client, err := lang.NewClient(ctx)
		if err != nil {
			log.Printf("Failed to create client: %v", err)
			config.err = err
		}
		langClient = client

	})

	defer logger.Flush()

	if config.Error() != nil {
		log.Println(config.Error())
		logger.StandardLogger(logging.Error).Println(config.Error())
		return config.Error()
	}

	job, err := pubSubPayloadToJob(&m)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Processing job: %s", job.ID)
	logger.StandardLogger(logging.Info).Printf("Processing job: %s", job.ID)
	err = updateJobStatus(job.ID, jobStatusProcessing)
	if err != nil {
		log.Printf("Error updating job status: %v", err)
		logger.StandardLogger(logging.Error).Println(err)
		return err
	}

	sent, err := processTerm(job.Term)
	if err != nil {
		log.Printf("Error updating job status: %v", err)
		logger.StandardLogger(logging.Error).Println(err)
		updateJobStatus(job.ID, jobStatusFailed)
		return err
	}

	// update job
	job.Result = sent

	// save results
	err = saveResults(job)
	if err != nil {
		log.Println(err)
		logger.StandardLogger(logging.Error).Println(err)
		updateJobStatus(job.ID, jobStatusFailed)
		return err
	}

	// save job status
	err = updateJobStatus(job.ID, jobStatusProcessed)
	if err != nil {
		log.Println(err)
		logger.StandardLogger(logging.Error).Println(err)
		return err
	}

	return nil

}

func pubSubPayloadToJob(m *PubSubMessage) (job *SentimentRequest, err error) {

	if m == nil {
		log.Println("Nil PubSubMessage")
		return nil, errors.New("PubSubMessage required")
	}

	d, err := base64.StdEncoding.DecodeString(m.Data)
	if err != nil {
		log.Printf("Decoding error: %v", err)
		return nil, err
	}

	j := &SentimentRequest{}

	err = json.Unmarshal([]byte(d), &j)
	if err != nil {
		log.Printf("Error deserislizing: %v", err)
		return nil, err
	}

	return j, nil

}

func processTerm(query string) (r *SentimentResult, err error) {

	consumerKey := os.Getenv("T_CONSUMER_KEY")
	consumerSecret := os.Getenv("T_CONSUMER_SECRET")
	accessToken := os.Getenv("T_ACCESS_TOKEN")
	accessSecret := os.Getenv("T_ACCESS_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		return nil, errors.New("Both, consumer key/secret and access token/secret are required")
	}

	// init convif
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// HTTP Client - will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	maxTweetID := int64(0)

	searchArgs := &twitter.SearchTweetParams{
		Query:   query,
		Count:   maxTweets,
		Lang:    "en",
		SinceID: maxTweetID,
		//MaxID:      maxTweetID,
		ResultType: "recent",
	}

	log.Printf("Search: %v\n", query)
	search, _, err := client.Search.Tweets(searchArgs)
	if err != nil {
		return nil, err
	}

	// results
	result := &SentimentResult{
		Tweets:    int64(len(search.Statuses)),
		Processed: time.Now(),
	}

	log.Printf("Found: %d", result.Tweets)

	for _, tweet := range search.Statuses {

		log.Printf("ID:%v", tweet.ID)

		txt := strings.TrimSuffix(tweet.Text, "\n")
		log.Printf("Text:%s", txt)

		sentiment, err := scoreSentiment(txt)

		if err != nil {
			log.Printf("Error while scoring: %v", err)
			return nil, err
		}

		//.Score, result.DocumentSentiment.Magnitude, nil
		log.Printf("Sentiment: %v", sentiment)

		if sentiment.Score < negativeScoreThreshold && sentiment.Magnitude > magnitudeThreshold {
			result.Negative++
		}

		if sentiment.Score > positiveScoreThreshold && sentiment.Magnitude > magnitudeThreshold {
			result.Positive++
		}

		result.Score += float64(sentiment.Score * sentiment.Magnitude)

	}

	return result, nil

}

/*
Clearly Positive*	"score": 0.8, 	"magnitude": 3.0
Clearly Negative*	"score": -0.6, 	"magnitude": 4.0
Neutral				"score": 0.1, 	"magnitude": 0.0
Mixed				"score": 0.0, 	"magnitude": 4.0
*/
func scoreSentiment(s string) (sentiment *langpb.Sentiment, err error) {

	result, err := langClient.AnalyzeSentiment(ctx, &langpb.AnalyzeSentimentRequest{
		Document: &langpb.Document{
			Source: &langpb.Document_Content{
				Content: s,
			},
			Type: langpb.Document_PLAIN_TEXT,
		},
		EncodingType: langpb.EncodingType_UTF8,
	})
	if err != nil {
		log.Printf("Error while scoring: %v", err)
		return nil, err
	}

	return result.DocumentSentiment, nil

}
