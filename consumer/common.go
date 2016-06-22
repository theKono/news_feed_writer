package consumer

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/theKono/orchid/model/dynamo"
	"github.com/theKono/orchid/model/mysql"
)

// insertIntoMysql inserts the record into database.
var insertIntoMysql = func(record mysql.Sharder) (err error) {
	if err = mysql.DBSessions[record.Shard()].Create(record).Error; err != nil {
		log.Println("Cannot insert record into database\n", err)
		return
	}

	return
}

// insertIntoDynamoDB puts the item into DynamoDB.
var insertIntoDynamoDB = func(doc *dynamodb.PutItemInput) (err error) {
	if _, err = dynamo.DynamoDBService.PutItem(doc); err != nil {
		log.Println("Cannot insert doc into dynamodb\n", err)
		return
	}

	return
}
