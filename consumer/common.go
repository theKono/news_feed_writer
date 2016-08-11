package consumer

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/theKono/orchid/model/dynamo"
	"github.com/theKono/orchid/util"
)

// insertIntoDynamoDB puts the item into DynamoDB.
var insertIntoDynamoDB = func(doc *dynamodb.PutItemInput) (err error) {
	defer util.MeasureExecTime(time.Now(), "insertIntoDynamoDB")

	if _, err = dynamo.DynamoDBService.PutItem(doc); err != nil {
		log.Println("Cannot insert doc into dynamodb\n", err)
		return
	}

	return
}
