package dynamo

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/theKono/orchid/cfg"
	"github.com/theKono/orchid/model/messagejson"
)

var tableName string

func init() {
	tableName = cfg.DynamoDBNewsFeedTableName
	if tableName == "" {
		log.Fatalln("NewsFeed dynamodb table name is required")
	}
}

// NewNewsFeed creates a *dynamodb.PutItemInput from a
// messagejson.NewsFeed instance.
var NewNewsFeed = func(nf *messagejson.NewsFeed) (p *dynamodb.PutItemInput, err error) {
	p = &dynamodb.PutItemInput{
		TableName: &tableName,
		Item: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprint(nf.ID)),
			},
			"user_id": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprint(nf.UserID)),
			},
			"summary": &dynamodb.AttributeValue{
				S: aws.String(nf.Summary),
			},
		},
	}
	return
}
