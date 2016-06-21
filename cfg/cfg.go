package cfg

import (
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ORCHID")

	Parallel = viper.GetInt("PARALLEL")

	MysqlMain = viper.GetString("MYSQL_MAIN")
	MysqlShard = viper.GetString("MYSQL_SHARD")
	MysqlDebug = viper.GetInt("MYSQL_DEBUG") == 1

	SqsRegion = viper.GetString("SQS_REGION")
	SqsQueueURL = viper.GetString("SQS_QUEUE_URL")

	DynamoDBRegion = viper.GetString("DYNAMODB_REGION")
	DynamoDBNewsFeedTableName = viper.GetString("DYNAMODB_NEWS_FEED_TABLE")
}

// Parallel specifies the number of consumers at the same time.
var Parallel int

// MysqlMain specifies the main database source name.
// Ex: "user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
var MysqlMain string

// MysqlShard specifies the second database source name.
// Ex: "user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
var MysqlShard string

// MysqlDebug specifies whether or not to print more information about
// SQL.
var MysqlDebug bool

// SqsRegion specifies the region of the SQS service.
var SqsRegion string

// SqsQueueURL specifies the SQS Queue URL.
var SqsQueueURL string

// DynamoDBRegion specifies the region of Dynamodb service.
var DynamoDBRegion string

// DynamoDBNewsFeedTableName specifies the table name of NewsFeed document.
var DynamoDBNewsFeedTableName string
