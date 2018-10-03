# sentimenter

The `sentimenter` app allows the user to query for public sentiment from recent tweets for specific topic.

Example of multi-step process leveraging GCF and multiple back-end services:

* [Firestore](https://cloud.google.com/firestore/) - To store app data in a Cloud-native way at global scale
* [Natural Language API](https://cloud.google.com/natural-language/) - ML API to derive insights from unstructured text
* [Stackdriver](https://cloud.google.com/stackdriver/) - To monitor app, services, and functions as well as the underlining infrastructure

> To experiment with the new go 1.11 support in Google Cloud Functions sign up for the [private alpha](https://goo.gl/forms/rwRxKsajWXmdwwPt1)

## Usage

### 1. Term Submission

The `submit` function which the user can invoke over HTTPS with their search `term` will create a `job`, save it with `Received` state in Spanner DB, and queue that job for processing in Pub/Sub topic.

```shell
HTTPS_TRIGGER_URL=$(gcloud alpha functions describe sentimenter-submit \
    --format='value(httpsTrigger.url)')
curl "${HTTPS_TRIGGER_URL}?term=serverless"
```

Returns

```json
{
  "id": "tid-d50ff5b2-2120-4587-a99e-c4aea5c3f592",
  "created_on": "2018-08-16T14:52:20.195459344Z",
  "search_term": "serverless",
  "status": "Received",
  "result": null
}
```

### 2. Job Processing (Background)

The `process` function will be automatically triggered by GCF when a new job is written into Firestore. GCF uses the `providers/cloud.firestore/eventTypes/document.create` event to fire when new record is created. The processor will change the state of that `job` to `Processing`, retrieve tweets using Twitter API, score each tweet's sentiment using Google's Natural Language API, and then derive total score using `sentimenter` algorithm. When done, the total score of that job is saved and the job status updated to `Processed`.

### 3. Job Status

Throughout the entire process, the user can invoke the `status` function over HTTPS and get the current state of the submitted job. If the job status is `Processed`, the status of the job will also include its score.

```shell
HTTPS_TRIGGER_URL=$(gcloud alpha functions describe sentimenter-status \
    --format='value(httpsTrigger.url)')
curl "${HTTPS_TRIGGER_URL}?id=tid-c24774a1-89df-4ec0-a962-121a36d6966c"
```

Result

```json
{
  "id": "tid-c24774a1-89df-4ec0-a962-121a36d6966c",
  "created_on": "2018-08-16T14:54:02.822679302Z",
  "search_term": "serverless",
  "status": "Processed",
  "status_url": "https://us-central1-s9-demo.cloudfunctions.net/sentimenter-status?id=6c211819-30ef-4bdb-a723-a5be4979c101",
  "result": {
    "processed_on": "2018-08-16T14:54:06.636848268Z",
    "tweets": 100,
    "positive": 9,
    "negative": 0,
    "score": 13.880000105127692
  }
}
```

> Note, while the positive or negative classification of each tweet is reliable, the overall score of the sentiment for all tweets is derived by combining sentiment with the magnitude of each tweet which tends to favour longer tweets. As a result, the score is only a relative indicator of the overall strength of the sentiment and probably meaningless in itself.


## Setup


### Firestore DB

//TODO: add content

### Functions

First define the necessary environment variables

```shell
FS_COLL_NAME="projects/s9-demo/databases/(default)/documents/jobs/{id}"
T_VARS="T_CONSUMER_KEY=${T_CONSUMER_KEY},T_CONSUMER_SECRET=${T_CONSUMER_SECRET}"
T_VARS="${T_VARS},T_ACCESS_TOKEN=${T_ACCESS_TOKEN},T_ACCESS_SECRET=${T_ACCESS_SECRET}"
```

> Note, I'm obfuscating the Twitter API variables by pulling them form my local variables.
> You can just type these keys here if you need to. See [this instructions](https://developer.twitter.com/en/docs/basics/authentication/guides/access-tokens.html)
> on how to create Twitter API credentials see

Then deploy the three functions using the GCP `gcloud` command, specifying the entry point as well as environment variables and few other parameters for each.


```shell
gcloud alpha functions deploy sentimenter-submit \
  --entry-point SubmitFunction \
  --memory 128MB \
  --region us-central1 \
  --runtime go111 \
  --trigger-http

gcloud alpha functions deploy sentimenter-status \
  --entry-point StatusFunction \
  --memory 128MB \
  --region us-central1 \
  --runtime go111 \
  --trigger-http

gcloud alpha functions deploy sentimenter-process \
  --entry-point ProcessorFunction \
  --set-env-vars $T_VARS \
  --memory 256MB \
  --region us-central1 \
  --runtime go111 \
  --trigger-event providers/cloud.firestore/eventTypes/document.create \
  --trigger-resource $FS_COLL_NAME \
  --timeout=540s
```

If everything goes well, you should see this kind of response to every one of these above functions

```shell
Deploying function (may take a while - up to 2 minutes)...done.
availableMemoryMb: 128
entryPoint: StatusFunction
httpsTrigger:
  url: https://us-central1-s9-demo.cloudfunctions.net/sentimenter-status
labels:
  deployment-tool: cli-gcloud
name: projects/s9-demo/locations/us-central1/functions/sentimenter-status
runtime: go111
serviceAccountEmail: s9-demo@appspot.gserviceaccount.com
sourceUploadUrl: ...
status: ACTIVE
timeout: 60s
updateTime: '2018-08-16T00:38:33Z'
versionId: '4'
```

## Learnings

### IDs

Using `UUIDv4` seemed like a good idea until on ocassion I would get error on save:

```shell
Error on job save: rpc error: code = InvalidArgument desc = Document name "projects/.../databases/(default)/documents/jobs/" has invalid trailing "/".
```

The resons for this is that Firestore IDs can't start with a number. To avoid this in my `getNewID` utility I prepend `tid-` to generated IDs.

## TODO

* (WIP) sort out duplicate logging
*  Add local logging
*
