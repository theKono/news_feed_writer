package dynamo

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/theKono/orchid/cfg"
)

// DynamoDBService is an AWS DynamoDB client.
var DynamoDBService *dynamodb.DynamoDB

func init() {
	if cfg.DynamoDBRegion == "" {
		log.Fatalln("DynamoDBRegion is required")
	}

	DynamoDBService = dynamodb.New(
		session.New(),
		aws.NewConfig().WithRegion(cfg.DynamoDBRegion),
	)
}
