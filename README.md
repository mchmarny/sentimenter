# sentimenter

Example of multi-step process leveraging GCF and multiple back-end services:

* [Cloud Spanner](https://cloud.google.com/spanner/) - horizontally scalable, strongly consistent, relational database service
* [Pub/Sub](https://cloud.google.com/pubsub/) - Ingest event streams at any scale from anywhere for real-time streaming
* [Google Cloud Natural Language API](https://cloud.google.com/natural-language/) - Derive insights from unstructured text using Google ML

The `sentimenter` solutions allows the user get the sentiment report from the last `100` tweets for submitted term. The solution includes:


## Term Submission

The `submitter` function which the user can invoke over HTTPS with their search `term` will create a `job`, save it with `Received` state in Spanner DB, and queue that job for processing in Pub/Sub topic.

```
HTTPS_TRIGGER_URL=$(gcloud alpha functions describe sentimenter-submitter \
    --format='value(httpsTrigger.url)')
curl "${HTTPS_TRIGGER_URL}?term=serverless"
```

Returns

```
{
  "id": "d50ff5b2-2120-4587-a99e-c4aea5c3f592",
  "created_on": "2018-08-16T14:52:20.195459344Z",
  "search_term": "serverless",
  "status": "Received",
  "result": null
}
```

## Job Processing (Background)

The `processor` function will be automatically triggered by GCF when a new job arrives on Pub/Sub topic. The processor will change the state of that `job` to `Processing`, retrieve last `100` tweets using Twitter API, and score each tweet's sentiment using Google's Natural Language API. When done, the score of that job will be saved in the Spanner DB and the job status updated to `Processed`.


## Job Status

Throughout the entire process, the user can invoke the `status` function over HTTPS and get the current state of the submitted job. If the job status is `Processed`, the status of the job will also include its score.

```
HTTPS_TRIGGER_URL=$(gcloud alpha functions describe sentimenter-status \
    --format='value(httpsTrigger.url)')
curl "${HTTPS_TRIGGER_URL}?id=c24774a1-89df-4ec0-a962-121a36d6966c"
```

Result

```
{
  "id": "6c211819-30ef-4bdb-a723-a5be4979c101",
  "created_on": "2018-08-16T14:54:02.822679302Z",
  "search_term": "serverless",
  "status": "Processed",
  "result": {
    "processed_on": "2018-08-16T14:54:06.636848268Z",
    "tweets": 100,
    "positive": 9,
    "negative": 0,
    "score": 13.880000105127692
  }
}
```

> None of the Cloud Functions in this example know about each other. They only interaction point between the, is the state persisted in the Spanner DB and the payloads on the Pub/Sub queue
