// +build unit

package consumer

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	awsSqs "github.com/aws/aws-sdk-go/service/sqs"

	"github.com/theKono/orchid/model/dynamo"
	"github.com/theKono/orchid/model/messagejson"
	"github.com/theKono/orchid/model/mysql"
	"github.com/theKono/orchid/sqs"
)

var originalFunctions = map[string]interface{}{
	"messagejson.NewNewsFeed": messagejson.NewNewsFeed,
	"mysql.NewNewsFeed":       mysql.NewNewsFeed,
	"dynamo.NewNewsFeed":      dynamo.NewNewsFeed,
	"insertIntoMysql":         insertIntoMysql,
	"insertIntoDynamoDB":      insertIntoDynamoDB,
}

func restoreFunctions() {
	messagejson.NewNewsFeed = originalFunctions["messagejson.NewNewsFeed"].(func(*string) (*messagejson.NewsFeed, error))
	mysql.NewNewsFeed = originalFunctions["mysql.NewNewsFeed"].(func(*messagejson.NewsFeed) (*mysql.NewsFeed, error))
	dynamo.NewNewsFeed = originalFunctions["dynamo.NewNewsFeed"].(func(*messagejson.NewsFeed) (*dynamodb.PutItemInput, error))
	insertIntoMysql = originalFunctions["insertIntoMysql"].(func(*mysql.NewsFeed) error)
	insertIntoDynamoDB = originalFunctions["insertIntoDynamoDB"].(func(*dynamodb.PutItemInput) error)
}

func mockMessagejsonNewNewsFeed(ret *messagejson.NewsFeed) {
	messagejson.NewNewsFeed = func(*string) (*messagejson.NewsFeed, error) {
		return ret, nil
	}
}

func mockMysqlNewNewsFeed(ret *mysql.NewsFeed) {
	mysql.NewNewsFeed = func(*messagejson.NewsFeed) (*mysql.NewsFeed, error) {
		return ret, nil
	}
}

func mockDynamoNewNewsFeed(ret *dynamodb.PutItemInput) {
	dynamo.NewNewsFeed = func(*messagejson.NewsFeed) (*dynamodb.PutItemInput, error) {
		return ret, nil
	}
}

func TestConsumeNewsFeedMessage(t *testing.T) {
	message := awsSqs.Message{}

	orig := sqs.DeleteMessage
	callDeleteMessage := false
	sqs.DeleteMessage = func(*awsSqs.Message) error {
		callDeleteMessage = true
		return nil
	}
	defer func() { sqs.DeleteMessage = orig }()

	// When messagejson.NewNewsFeed fails
	func() {
		var arg *string
		messagejson.NewNewsFeed = func(s *string) (*messagejson.NewsFeed, error) {
			arg = s
			return nil, errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		message.Body = aws.String("Body")

		if consumeNewsFeedMessage(&message) == nil {
			t.Fatal("Expect consumeNewsFeedMessage() to return error")
		}

		if *arg != "Body" {
			t.Fatal("messagejson.NewNewsFeed recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When mysql.NewNewsFeed fails
	func() {
		var arg *messagejson.NewsFeed
		var newsFeed = &messagejson.NewsFeed{}

		mockMessagejsonNewNewsFeed(newsFeed)
		mysql.NewNewsFeed = func(nf *messagejson.NewsFeed) (*mysql.NewsFeed, error) {
			arg = nf
			return nil, errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeNewsFeedMessage(&message) == nil {
			t.Fatal("Expect consumeNewsFeedMessage() to return error")
		}

		if arg != newsFeed {
			t.Fatal("mysql.NewNewsFeed recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When dynamo.NewNewsFeed fails
	func() {
		var arg *messagejson.NewsFeed
		var mjNewsFeed = &messagejson.NewsFeed{}
		var msNewsFeed = &mysql.NewsFeed{}

		mockMessagejsonNewNewsFeed(mjNewsFeed)
		mockMysqlNewNewsFeed(msNewsFeed)
		dynamo.NewNewsFeed = func(nf *messagejson.NewsFeed) (*dynamodb.PutItemInput, error) {
			arg = nf
			return nil, errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeNewsFeedMessage(&message) == nil {
			t.Fatal("Expect consumeNewsFeedMessage() to return error")
		}

		if arg != mjNewsFeed {
			t.Fatal("dynamo.NewNewsFeed recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When insertIntoMysql fails
	func() {
		var arg *mysql.NewsFeed
		var mjNewsFeed = &messagejson.NewsFeed{}
		var msNewsFeed = &mysql.NewsFeed{}
		var dyPutItemInput = &dynamodb.PutItemInput{}

		mockMessagejsonNewNewsFeed(mjNewsFeed)
		mockMysqlNewNewsFeed(msNewsFeed)
		mockDynamoNewNewsFeed(dyPutItemInput)
		insertIntoMysql = func(nf *mysql.NewsFeed) error {
			arg = nf
			return errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeNewsFeedMessage(&message) == nil {
			t.Fatal("Expect consumeNewsFeedMessage() to return error")
		}

		if arg != msNewsFeed {
			t.Fatal("insertIntoMysql recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When insertIntoDynamoDB fails
	func() {
		var arg *dynamodb.PutItemInput
		var mjNewsFeed = &messagejson.NewsFeed{}
		var msNewsFeed = &mysql.NewsFeed{}
		var dyPutItemInput = &dynamodb.PutItemInput{}

		mockMessagejsonNewNewsFeed(mjNewsFeed)
		mockMysqlNewNewsFeed(msNewsFeed)
		mockDynamoNewNewsFeed(dyPutItemInput)
		insertIntoMysql = func(nf *mysql.NewsFeed) error {
			return nil
		}
		insertIntoDynamoDB = func(pii *dynamodb.PutItemInput) error {
			arg = pii
			return errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeNewsFeedMessage(&message) == nil {
			t.Fatal("Expect consumeNewsFeedMessage() to return error")
		}

		if arg != dyPutItemInput {
			t.Fatal("insertIntoDynamoDB recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()
}
