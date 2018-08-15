# sentimenter

Example of long running process leveraging GCF and multiple back-end services (Spanner, PubSub, NLP, etc.)


## Submit Job

```
curl https://$(SERVICE_HOST}/sentimenter-submitter?term=google
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
curl https://$(SERVICE_HOST}/sentimenter-status?id=c24774a1-89df-4ec0-a962-121a36d6966c
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
curl https://$(SERVICE_HOST}/sentimenter-result?id=c24774a1-89df-4ec0-a962-121a36d6966c
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
