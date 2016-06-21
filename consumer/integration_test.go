// +build integration

package consumer

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/theKono/orchid/model/dynamo"
	"github.com/theKono/orchid/model/messagejson"
	"github.com/theKono/orchid/model/mysql"
)

func TestInsertIntoMysql(t *testing.T) {
	var msNewsFeed *mysql.NewsFeed
	var mjNewsFeed *messagejson.NewsFeed
	var record = mysql.NewsFeed{}

	mjNewsFeed = &messagejson.NewsFeed{UserID: 2, Summary: "{}"}
	mjNewsFeed.GenerateID()
	msNewsFeed, _ = mysql.NewNewsFeed(mjNewsFeed)
	if err := insertIntoMysql(msNewsFeed); err != nil {
		t.Fatal("Expect insertIntoMysql not to return error\n", err)
	}

	mysql.DBSessions[0].Where("newsfeed_id = ?", mjNewsFeed.ID).First(&record)
	if record.NewsfeedID == 0 {
		t.Fatal("Expect record exist")
	} else {
		mysql.DBSessions[0].Delete(&record)
	}

	mjNewsFeed = &messagejson.NewsFeed{UserID: 1, Summary: "{}"}
	mjNewsFeed.GenerateID()
	msNewsFeed, _ = mysql.NewNewsFeed(mjNewsFeed)
	if err := insertIntoMysql(msNewsFeed); err != nil {
		t.Fatal("Expect insertIntoMysql not to return error\n", err)
	}

	mysql.DBSessions[1].Where("newsfeed_id = ?", mjNewsFeed.ID).First(&record)
	if record.NewsfeedID == 0 {
		t.Fatal("Expect record exist")
	} else {
		mysql.DBSessions[1].Delete(&record)
	}
}

func TestInsertIntoDynamoDB(t *testing.T) {
	mjNewsFeed := &messagejson.NewsFeed{UserID: rand.Int31(), Summary: "{}"}
	mjNewsFeed.GenerateID()

	pii, err := dynamo.NewNewsFeed(mjNewsFeed)
	if err := insertIntoDynamoDB(pii); err != nil {
		t.Fatal("Expect insertIntoDynamoDB not to return error\n", err)
	}

	resp, err := dynamo.DynamoDBService.GetItem(
		&dynamodb.GetItemInput{
			TableName: pii.TableName,
			Key: map[string]*dynamodb.AttributeValue{
				"user_id": {N: aws.String(fmt.Sprint(mjNewsFeed.UserID))},
				"id":      {N: aws.String(fmt.Sprint(mjNewsFeed.ID))},
			},
		},
	)
	if err != nil {
		t.Fatal("Expect GetItemInput not to return error\n", err)
	}
	if *resp.Item["user_id"].N != fmt.Sprint(mjNewsFeed.UserID) {
		t.Fatal("Expect user_id to equal `%v`, but got `%v`", mjNewsFeed.UserID, *resp.Item["user_id"].N)
	}
	if *resp.Item["id"].N != fmt.Sprint(mjNewsFeed.ID) {
		t.Fatal("Expect id to equal `%v`, but got `%v`", mjNewsFeed.ID, *resp.Item["id"].N)
	}
	if *resp.Item["summary"].S != mjNewsFeed.Summary {
		t.Fatal("Expect summary to equal `%v`, but got `%v`", mjNewsFeed.Summary, *resp.Item["summary"].S)
	}

	_, err = dynamo.DynamoDBService.DeleteItem(
		&dynamodb.DeleteItemInput{
			TableName: pii.TableName,
			Key: map[string]*dynamodb.AttributeValue{
				"user_id": {N: aws.String(fmt.Sprint(mjNewsFeed.UserID))},
				"id":      {N: aws.String(fmt.Sprint(mjNewsFeed.ID))},
			},
		},
	)
	if err != nil {
		t.Fatal("Expect DeleteItem not to return error\n", err)
	}
}
