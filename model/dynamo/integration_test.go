// +build integration

package dynamo

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/theKono/orchid/model/messagejson"
)

func TestNewNewsFeed(t *testing.T) {
	nf := &messagejson.NewsFeed{
		messagejson.SocialFeed{UserID: rand.Int31(), Summary: "{}"},
	}
	nf.GenerateID()

	pii, err := NewNewsFeed(nf)
	if err != nil {
		t.Fatal("Expect NewNewsFeed() not to return error\n", err)
	}

	_, err = DynamoDBService.PutItem(pii)
	if err != nil {
		t.Fatal("Expect PutItem() not to return error\n", err)
	}

	_, err = DynamoDBService.DeleteItem(
		&dynamodb.DeleteItemInput{
			TableName: &newsFeedTableName,
			Key: map[string]*dynamodb.AttributeValue{
				"user_id": &dynamodb.AttributeValue{N: aws.String(fmt.Sprint(nf.UserID))},
				"id":      &dynamodb.AttributeValue{N: aws.String(fmt.Sprint(nf.ID))},
			},
		},
	)
	if err != nil {
		t.Fatal("Expect DeleteItem not to return error\n", err)
	}
}

func TestNewNotification(t *testing.T) {
	n := &messagejson.Notification{
		messagejson.SocialFeed{UserID: rand.Int31(), Summary: "{}"},
	}
	n.GenerateID()

	pii, err := NewNotification(n)
	if err != nil {
		t.Fatal("Expect NewNotification() not to return error\n", err)
	}

	_, err = DynamoDBService.PutItem(pii)
	if err != nil {
		t.Fatal("Expect PutItem() not to return error\n", err)
	}

	_, err = DynamoDBService.DeleteItem(
		&dynamodb.DeleteItemInput{
			TableName: &notificationTableName,
			Key: map[string]*dynamodb.AttributeValue{
				"user_id": &dynamodb.AttributeValue{N: aws.String(fmt.Sprint(n.UserID))},
				"id":      &dynamodb.AttributeValue{N: aws.String(fmt.Sprint(n.ID))},
			},
		},
	)
	if err != nil {
		t.Fatal("Expect DeleteItem not to return error\n", err)
	}
}
