# sentimenter

Example of long running process leveraging GCF and multiple back-end services (Spanner, PubSub, NLP, etc.)


## Submit Job

```
HTTPS_TRIGGER_URL=$(gcloud alpha functions describe sentimenter-submitter \
  --format='value(httpsTrigger.url)')
curl https://$(HTTPS_TRIGGER_URL}/sentimenter-submitter?term=google
```

Result

```
{
    "id": "c24774a1-89df-4ec0-a962-121a36d6966c",
    "created_on": "2018-08-15T21:19:06.869021913Z",
    "search_term": "google",
    "status": "Received",
    "result": null
}
```

## Check Job Status

```
HTTPS_TRIGGER_URL=$(gcloud alpha functions describe sentimenter-status \
  --format='value(httpsTrigger.url)')
curl https://$(HTTPS_TRIGGER_URL}/sentimenter-status?id=c24774a1-89df-4ec0-a962-121a36d6966c
```

Result

```
{
    "id": "c24774a1-89df-4ec0-a962-121a36d6966c",
    "created_on": "2018-08-15T21:19:06.869021913Z",
    "search_term": "google",
    "status": "Processing",
    "result": null
}
```



## Get Job Results

```
HTTPS_TRIGGER_URL=$(gcloud alpha functions describe sentimenter-result \
  --format='value(httpsTrigger.url)')
curl https://$(HTTPS_TRIGGER_URL}/sentimenter-result?id=c24774a1-89df-4ec0-a962-121a36d6966c
```

Result

```
{
    "id": "c24774a1-89df-4ec0-a962-121a36d6966c",
    "created_on": "2018-08-15T21:19:06.869021913Z",
    "search_term": "google",
    "status": "Processed",
    "result": {
        "processed_on": "2018-08-15T21:19:06.869021913Z",
        "tweets": 100,
        "positive": 34,
        "negative": 4,
        "score": 23.91
    }
}
```
