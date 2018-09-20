package sentimenter

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	lang "cloud.google.com/go/language/apiv1"
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
func ProcessorFunction(ctx context.Context, e FirestoreEvent) error {

	once.Do(func() {
		configInitializer("process")
		client, err := lang.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create lang client: %v", err)
		}
		langClient = client
	})

	defer logger.Flush()

	infoLogger.Printf("Parsing Job: %v", e.Value.Fields)

	// data to map to struc
	m := e.Value.Fields.(map[string]interface{})
	job, err := eventMapToJob(m)
	if err != nil {
		errLogger.Printf("Error parsing job: %v - %v", err, job)
		return errors.New("Error processing job")
	}

	// processing job
	infoLogger.Printf("Processing job: %s", job.ID)
	err = updateStatusAndSaveJob(job, jobStatusProcessing)
	if err != nil {
		return err
	}

	// process term
	infoLogger.Printf("Processing term: %s", job.Term)
	sent, err := processTerm(job.Term)
	if err != nil {
		updateStatusAndSaveJob(job, jobStatusFailed)
		return err
	}

	// update job
	job.Result = sent
	job.Status = jobStatusProcessed
	infoLogger.Printf("Saving results: %v", job)
	err = saveJob(job)
	if err != nil {
		// best effort only, no error capture on purpose
		updateStatusAndSaveJob(job, jobStatusFailed)
		return err
	}

	infoLogger.Printf("Job processed: %s", job.ID)

	return nil

}

func updateStatusAndSaveJob(job *SentimentRequest, status string) error {
	job.Status = status
	infoLogger.Printf("Changing job status: %s (%s)", job.ID, job.Status)
	err := saveJob(job)
	if err != nil {
		errLogger.Println(err)
		job.Status = jobStatusFailed
		saveJob(job)
		return err
	}

	return nil
}

// TODO: This has be easier, map is of proto serialized struc not json
func eventMapToJob(m map[string]interface{}) (job *SentimentRequest, err error) {

	if m == nil {
		return nil, errors.New("Event data required")
	}

	j := &SentimentRequest{}
	id, e := getMapValAsString(m, "id")
	if err != nil {
		return nil, e
	}
	j.ID = id

	term, e := getMapValAsString(m, "search_term")
	if err != nil {
		return nil, e
	}
	j.Term = term

	status, e := getMapValAsString(m, "status")
	if err != nil {
		return nil, e
	}
	j.Status = status

	url, e := getMapValAsString(m, "status_url")
	if err != nil {
		return nil, e
	}
	j.URL = url

	on, err := getMapValAsTimestamp(m, "created_on")
	if err != nil {
		return nil, e
	}
	j.On = on

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

	infoLogger.Printf("Search: %v", query)
	search, _, err := client.Search.Tweets(searchArgs)
	if err != nil {
		return nil, err
	}

	// results
	result := &SentimentResult{
		Tweets:    int64(len(search.Statuses)),
		Processed: time.Now(),
	}

	infoLogger.Printf("Found: %d", result.Tweets)
	for _, tweet := range search.Statuses {

		txt := strings.TrimSuffix(tweet.Text, "\n")

		sentiment, err := scoreSentiment(txt)

		if err != nil {
			errLogger.Printf("Error while scoring: %v", err)
			return result, nil
		}

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
		errLogger.Printf("Error while scoring: %v", err)
		return nil, err
	}

	return result.DocumentSentiment, nil

}
