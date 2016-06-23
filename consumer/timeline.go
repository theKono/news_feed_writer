// +build timeline

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

// consumeTimelineMessage is a Timeline consumer.
//
// It deletes the SQS message regardless of its validity. If it is a
// valid message, then it will insert the Timeline both into MySQL and
// DynamoDB.
var consumeTimelineMessage = func(message *awsSqs.Message) error {
	var (
		newsFeed *messagejson.Timeline
		record   *mysql.Timeline
		doc      *dynamodb.PutItemInput
		err      error
	)

	defer sqs.DeleteMessage(message)

	if newsFeed, err = messagejson.NewTimeline(message.Body); err != nil {
		log.Println("Cannot parse news feed message\n", err)
		return err
	}

	if record, err = mysql.NewTimeline(newsFeed); err != nil {
		log.Println("Cannot create news feed model\n", err)
		return err
	}

	if doc, err = dynamo.NewTimeline(newsFeed); err != nil {
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

// ConsumeTimeline is a function that implements MessageConsumer
// interface.
//
// It is derived from consumeTimelineMessage decorated by
// DecorateConsumeFn.
var ConsumeTimeline = DecorateConsumeFn(consumeTimelineMessage)
