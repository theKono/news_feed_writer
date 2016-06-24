// +build notification

package dynamo

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/theKono/orchid/cfg"
	"github.com/theKono/orchid/model/messagejson"
	"github.com/theKono/orchid/model/mysql"
)

var notificationTableName string

func init() {
	notificationTableName = cfg.DynamoDBNotificationTableName
	if notificationTableName == "" {
		log.Fatalln("Notification dynamodb table name is required")
	}
}

// NewNotification creates a *dynamodb.PutItemInput from a
// messagejson.Notification instance.
var NewNotification = func(n *messagejson.Notification) (p *dynamodb.PutItemInput, err error) {
	p = &dynamodb.PutItemInput{
		TableName: &notificationTableName,
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
			"state": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprint(mysql.UnseenAndUnread)),
			},
		},
	}
	return
}
