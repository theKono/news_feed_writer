# news feed writer

It is a microservice Kono used to create news feed to users. The microservice keeps polling AWS SQS, write the news feeds generated from other app/service to database.

## Installation
- The project is tested against Go 1.6
- `GOPATH` should be set properly
- The dependencies is commited into `vendor/`, there should be no worry to dependency problem.

### Configuration
- Required
    - `MYSQL_SHARD_0_URI`: The shard 0 database source name, e.g. `user:pwd@/db1?charset=utf8&parseTime=True&loc=Local`
    - `MYSQL_SHARD_1_URI`: The shard 1 database source name, e.g. `user:pwd@/db1?charset=utf8&parseTime=True&loc=Local`
    - `SQS_QUEUE_REGION`: The SQS queue region, e.g. `ap-southeast-1`.
    - `SQS_QUEUE_URL`: The SQS queue URL, e.g. `https://sqs.ap-southeast-1.amazonaws.com/xxx/test-sqs-queue`.
- Optional
    - `PARALLEL`: How many workers will be created to poll SQS queue. Default is `5`.
    - `DEBUG`: If it is set to `1`, then MySQL statement will be logged. Default is `0`.

### Tests
```
cd main
MYSQL_SHARD_0_URI="user:pwd@/db1?charset=utf8&parseTime=True&loc=Local" \
MYSQL_SHARD_1_URI="user:pwd@/db2?charset=utf8&parseTime=True&loc=Local" \
DEBUG=1 \
go test -v
```

### Run
If the executable is `main`, set required environment and run it,
```
MYSQL_SHARD_0_URI="user:pwd@/db1?charset=utf8&parseTime=True&loc=Local" \
MYSQL_SHARD_1_URI="user:pwd@/db2?charset=utf8&parseTime=True&loc=Local" \
SQS_QUEUE_REGION="ap-southeast-1" \
SQS_QUEUE_URL="https://sqs.ap-southeast-1.amazonaws.com/xxx/test-sqs-queue" \
PARALLEL=10 \
main
```

### Stop
Send `ctrl-c` signal to the executable, it will stop gracefully.
