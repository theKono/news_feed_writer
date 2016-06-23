// +build notification

package dynamo

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/theKono/orchid/cfg"
	"github.com/theKono/orchid/model/messagejson"
)

var timelineTableName string

func init() {
	timelineTableName = cfg.DynamoDBTimelineTableName
	if timelineTableName == "" {
		log.Fatalln("Timeline dynamodb table name is required")
	}
}

// NewTimeline creates a *dynamodb.PutItemInput from a
// messagejson.Timeline instance.
var NewTimeline = func(n *messagejson.Timeline) (p *dynamodb.PutItemInput, err error) {
	p = &dynamodb.PutItemInput{
		TableName: &timelineTableName,
		Item: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprint(n.ID)),
			},
			"user_id": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprint(n.UserID)),
			},
			"summary": &dynamodb.AttributeValue{
				S: aws.String(n.Summary),
			},
		},
	}
	return
}
