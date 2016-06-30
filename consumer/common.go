package consumer

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/theKono/orchid/model/dynamo"
	"github.com/theKono/orchid/model/mysql"
	"github.com/theKono/orchid/util"
)

// insertIntoMysql inserts the record into database.
var insertIntoMysql = func(record mysql.Sharder) (err error) {
	defer util.MeasureExecTime(time.Now(), "insertIntoMysql")

	if err = mysql.DBSessions[record.Shard()].Create(record).Error; err != nil {
		log.Println("Cannot insert record into database\n", err)
		return
	}

	return
}

// insertIntoDynamoDB puts the item into DynamoDB.
var insertIntoDynamoDB = func(doc *dynamodb.PutItemInput) (err error) {
	defer util.MeasureExecTime(time.Now(), "insertIntoDynamoDB")

	if _, err = dynamo.DynamoDBService.PutItem(doc); err != nil {
		log.Println("Cannot insert doc into dynamodb\n", err)
		return
	}

	return
}
