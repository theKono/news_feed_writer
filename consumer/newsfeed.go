package consumer

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	awsSqs "github.com/aws/aws-sdk-go/service/sqs"

	"github.com/theKono/orchid/model/dynamo"
	"github.com/theKono/orchid/model/messagejson"
	"github.com/theKono/orchid/model/mysql"
	"github.com/theKono/orchid/sqs"
)

// insertIntoMysql inserts the record into database.
var insertIntoMysql = func(record *mysql.NewsFeed) (err error) {
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

// consumeNewsFeedMessage is the a NewsFeed consumer.
//
// It deletes the SQS message regardless of its validity. If it is a
// valid message, then it will insert the NewsFeed both into MySQL and
// DynamoDB.
var consumeNewsFeedMessage = func(message *awsSqs.Message) error {
	var (
		newsFeed *messagejson.NewsFeed
		record   *mysql.NewsFeed
		doc      *dynamodb.PutItemInput
		err      error
	)

	defer sqs.DeleteMessage(message)

	if newsFeed, err = messagejson.NewNewsFeed(message.Body); err != nil {
		log.Println("Cannot parse news feed message\n", err)
		return err
	}

	if record, err = mysql.NewNewsFeed(newsFeed); err != nil {
		log.Println("Cannot create news feed model\n", err)
		return err
	}

	if doc, err = dynamo.NewNewsFeed(newsFeed); err != nil {
		log.Println("Cannot create news feed dynamodb model\n", err)
		return err
	}

	if err = insertIntoMysql(record); err != nil {
		return err
	}

	if err = insertIntoDynamoDB(doc); err != nil {
		return err
	}

	return nil
}

// ConsumeNewsFeed is a function that implements MessageConsumer
// interface.
//
// It is derived from consumeNewsFeedMessage decorated by
// DecorateConsumeFn.
var ConsumeNewsFeed = DecorateConsumeFn(consumeNewsFeedMessage)