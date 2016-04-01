package config

import (
	"log"
	"os"
	"strconv"
)

const (
	DefaultParallel            = 5
	DefaultMaxNumberOfMessages = 10
	DefaultWaitTimeSeconds     = 20
)

var (
	mysqlShardUris = []string{
		os.Getenv("MYSQL_SHARD_0_URI"),
		os.Getenv("MYSQL_SHARD_1_URI"),
	}
	sqsQueueRegion      = os.Getenv("SQS_QUEUE_REGION")
	sqsQueueURL         = os.Getenv("SQS_QUEUE_URL")
	parallel            = os.Getenv("PARALLEL")
	debug               = (os.Getenv("DEBUG") == "1")
	paralleln           int
	maxNumberOfMessages = DefaultMaxNumberOfMessages
	waitTimeSeconds     = DefaultWaitTimeSeconds
)

func GetSqsQueueRegion() string {
	return sqsQueueRegion
}

func GetSqsQueueUrl() string {
	return sqsQueueURL
}

func GetParallel() int {
	if paralleln != 0 {
		return paralleln
	}

	n, err := strconv.Atoi(parallel)

	if err != nil {
		log.Println("Cannot parse PARALLEL environment variable, use default")
		paralleln = DefaultParallel
	} else {
		paralleln = n
	}

	return paralleln
}

func GetMaxNumberOfMessages() int {
	return maxNumberOfMessages
}

func GetWaitTimeSeconds() int {
	return waitTimeSeconds
}

func GetShardDbUri(shard int) string {
	return mysqlShardUris[shard]
}

func IsDebug() bool {
	return debug
}
