package cfg

import (
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ORCHID")

	Parallel = viper.GetInt("PARALLEL")

	SqsRegion = viper.GetString("SQS_REGION")
	SqsQueueURL = viper.GetString("SQS_QUEUE_URL")

	DynamoDBRegion = viper.GetString("DYNAMODB_REGION")
	DynamoDBNewsFeedTableName = viper.GetString("DYNAMODB_NEWS_FEED_TABLE")
	DynamoDBNotificationTableName = viper.GetString("DYNAMODB_NOTIFICATION_TABLE")
	DynamoDBTimelineTableName = viper.GetString("DYNAMODB_TIMELINE_TABLE")
}

// Parallel specifies the number of consumers at the same time.
var Parallel int

// SqsRegion specifies the region of the SQS service.
var SqsRegion string

// SqsQueueURL specifies the SQS Queue URL.
var SqsQueueURL string

// DynamoDBRegion specifies the region of Dynamodb service.
var DynamoDBRegion string

// DynamoDBNewsFeedTableName specifies the table name of NewsFeed document.
var DynamoDBNewsFeedTableName string

// DynamoDBNotificationTableName specifies the table name of Notification document.
var DynamoDBNotificationTableName string

// DynamoDBTimelineTableName specifies the table name of Timeline document.
var DynamoDBTimelineTableName string
