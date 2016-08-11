# Orchid

It is a microservice Kono used to create news feed, notifications, timeline units to users. The microservice keeps polling AWS SQS, write the feeds generated from other app/service to database.

## Installation
- The project is tested against Go 1.6
- `GOPATH` should be set properly
- Use [glide](https://github.com/Masterminds/glide) to manage vendor packages. Please see the glide's installation guide.
    - After installing glide, `glide install` should be good to go.

## Configuration
Orchid uses configuration from environment variables. All environment variables start with `ORCHID_` prefix. We list the required and optional environment variable for each executable.

### bin/newsfeed-writer
- `ORCHID_DYNAMODB_REGION`
    - It specifies the AWS DynamoDB resource region. See `ORCHID_DYNAMODB_NEWS_FEED_TABLE`.
- `ORCHID_DYNAMODB_NEWS_FEED_TABLE`
    - The table name into which the news feed document will be inserted.
- `ORCHID_SQS_REGION`
    - It specifies the AWS SQS resouce region.
- `ORCHID_SQS_QUEUE_URL`
    - The URL of the AWS SQS queue. The workers will compete news feed data from the AWS SQS queue.
- `ORCHID_PARALLEL`
    - It specifies the number of worker to activate. The workers will compete message from AWS SQS queue.

### bin/notification-writer
- `ORCHID_DYNAMODB_REGION`
    - It specifies the AWS DynamoDB resource region. See `ORCHID_DYNAMODB_NOTIFICATION_TABLE`.
- `ORCHID_DYNAMODB_NOTIFICATION_TABLE`
    - The table name into which the notification document will be inserted.
- `ORCHID_SQS_REGION`
    - It specifies the AWS SQS resouce region.
- `ORCHID_SQS_QUEUE_URL`
    - The URL of the AWS SQS queue. The workers will compete notification data from the AWS SQS queue.
- `ORCHID_PARALLEL`
    - It specifies the number of worker to activate. The workers will compete message from AWS SQS queue.

### bin/timeline-writer
- `ORCHID_DYNAMODB_REGION`
    - It specifies the AWS DynamoDB resource region. See `ORCHID_DYNAMODB_TIMELINE_TABLE`.
- `ORCHID_DYNAMODB_TIMELINE_TABLE`
    - The table name into which the timeline document will be inserted.
- `ORCHID_SQS_REGION`
    - It specifies the AWS SQS resouce region.
- `ORCHID_SQS_QUEUE_URL`
    - The URL of the AWS SQS queue. The workers will compete timeline data from the AWS SQS queue.
- `ORCHID_PARALLEL`
    - It specifies the number of worker to activate. The workers will compete message from AWS SQS queue.

## Unit Test
```bash
ORCHID_DYNAMODB_REGION=... \
ORCHID_DYNAMODB_NEWS_FEED_TABLE=... \
ORCHID_DYNAMODB_NOTIFICATION_TABLE=... \
ORCHID_SQS_REGION=... \
ORCHID_SQS_QUEUE_URL=... \
go test -v -cover -tags "unit newsfeed notification timeline"
```

## Integration Test
```bash
ORCHID_DYNAMODB_REGION=... \
ORCHID_DYNAMODB_NEWS_FEED_TABLE=... \
ORCHID_DYNAMODB_NOTIFICATION_TABLE=... \
ORCHID_SQS_REGION=... \
ORCHID_SQS_QUEUE_URL=... \
go test -v -cover -tags "integration newsfeed notification timeline"
```

## Stop
Send `ctrl-c` signal to the executable, it will stop until all received messages are consumed.

## Deploy
### Deploy NewsFeed writer
```
cd deploy/newsfeed
./deploy_news_feed_writer <deployment-group> <app-version> <s3-bucket>
```
- `deployment-group`: AWS CodeDeploy deployment-group
- `app-version`: git tag version
- `s3-bucket`: The deploying instance region

### Deploy Notification writer
```
cd deploy/notification
./deploy_notification_writer <deployment-group> <app-version> <s3-bucket>
```
- `deployment-group`: AWS CodeDeploy deployment-group
- `app-version`: git tag version
- `s3-bucket`: The deploying instance region

### Deploy Timeline writer
```
cd deploy/timeline
./deploy_timeline_writer <deployment-group> <app-version> <s3-bucket>
```
- `deployment-group`: AWS CodeDeploy deployment-group
- `app-version`: git tag version
- `s3-bucket`: The deploying instance region
