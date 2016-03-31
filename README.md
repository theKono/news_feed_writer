# news_feed_writer
------------------

It is a microservice Kono used to create news feed to users. The microservice keeps polling AWS SQS, write the news feeds generated from other app/service to database.

## Installation
- The project is tested against Go 1.6
- `GOPATH` should be set properly
- The dependencies is commited into `vendor/`, there should be no worry to dependency problem.

### Configuration
- `SQS_QUEUE_REGION`: The SQS queue region, e.g. `ap-southeast-1`.
- `SQS_QUEUE_URL`: The SQS queue URL, e.g. `https://sqs.ap-southeast-1.amazonaws.com/xxx/test-sqs-queue`.
- `PARALLEL`: How many workers will be created to poll SQS queue. Default is `5`.
- `DEBUG`: If it is set to `1`, then MySQL statement will be logged. Default is `0`.

### Tests
```
cd main
go test -v
```
