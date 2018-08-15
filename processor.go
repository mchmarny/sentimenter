package sentimenter

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	lang "cloud.google.com/go/language/apiv1"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	langpb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

var (
	langClient *lang.Client
)

// ProcessorFunction processes pubsub messages
func ProcessorFunction(ctx context.Context, m PubSubMessage) error {

	config.once.Do(func() { configFunc() })
	config.once.Do(func() {
		{
			client, err := lang.NewClient(ctx)
			if err != nil {
				log.Panicf("Failed to create client: %v", err)
			}
			langClient = client
		}
	})

	if config.Error() != nil {
		log.Println(config.Error())
		return config.Error()
	}

	job, err := pubSubPayloadToJob(&m)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Processing job: %s", job.ID)

	//TODO: implement

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

func getTweets(query string) error {

	consumerKey := os.Getenv("T_CONSUMER_KEY")
	consumerSecret := os.Getenv("T_CONSUMER_SECRET")
	accessToken := os.Getenv("T_ACCESS_TOKEN")
	accessSecret := os.Getenv("T_ACCESS_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		return errors.New("Both, consumer key/secret and access token/secret are required")
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
		Count:   100,
		Lang:    "en",
		SinceID: maxTweetID,
		//MaxID:      maxTweetID,
		ResultType: "recent",
	}

	log.Printf("Search: %v\n", query)
	search, _, err := client.Search.Tweets(searchArgs)
	if err != nil {
		return err
	}

	// counter stuff
	log.Printf("Found: %d", len(search.Statuses))

	for _, tweet := range search.Statuses {

		log.Printf("ID:%v", tweet.ID)

		text := strings.TrimSuffix(tweet.FullText, "\n")
		log.Printf("Text:%s", text)

		sentiment, err := scoreSentiment(text)

		if err != nil {
			log.Printf("Error while scoring: %v", err)
			return err
		}

		//.Score, result.DocumentSentiment.Magnitude, nil
		log.Printf("Sentiment: %v", sentiment)

	}

	return nil

}

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
