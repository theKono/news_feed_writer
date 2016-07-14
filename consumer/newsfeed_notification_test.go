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
	"messagejson.NewNewsFeed":     messagejson.NewNewsFeed,
	"mysql.NewNewsFeed":           mysql.NewNewsFeed,
	"dynamo.NewNewsFeed":          dynamo.NewNewsFeed,
	"messagejson.NewNotification": messagejson.NewNotification,
	"mysql.NewNotification":       mysql.NewNotification,
	"dynamo.NewNotification":      dynamo.NewNotification,
	"messagejson.NewTimeline":     messagejson.NewTimeline,
	"mysql.NewTimeline":           mysql.NewTimeline,
	"dynamo.NewTimeline":          dynamo.NewTimeline,
	"insertIntoMysql":             insertIntoMysql,
	"insertIntoDynamoDB":          insertIntoDynamoDB,
}

func restoreFunctions() {
	messagejson.NewNewsFeed = originalFunctions["messagejson.NewNewsFeed"].(func(*string) (*messagejson.NewsFeed, error))
	mysql.NewNewsFeed = originalFunctions["mysql.NewNewsFeed"].(func(*messagejson.NewsFeed) (*mysql.NewsFeed, error))
	dynamo.NewNewsFeed = originalFunctions["dynamo.NewNewsFeed"].(func(*messagejson.NewsFeed) (*dynamodb.PutItemInput, error))

	messagejson.NewNotification = originalFunctions["messagejson.NewNotification"].(func(*string) (*messagejson.Notification, error))
	mysql.NewNotification = originalFunctions["mysql.NewNotification"].(func(*messagejson.Notification) (*mysql.Notification, error))
	dynamo.NewNotification = originalFunctions["dynamo.NewNotification"].(func(*messagejson.Notification) (*dynamodb.PutItemInput, error))

	messagejson.NewTimeline = originalFunctions["messagejson.NewTimeline"].(func(*string) (*messagejson.Timeline, error))
	mysql.NewTimeline = originalFunctions["mysql.NewTimeline"].(func(*messagejson.Timeline) (*mysql.Timeline, error))
	dynamo.NewTimeline = originalFunctions["dynamo.NewTimeline"].(func(*messagejson.Timeline) (*dynamodb.PutItemInput, error))

	insertIntoMysql = originalFunctions["insertIntoMysql"].(func(mysql.Sharder) error)
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

func mockMessagejsonNewNotification(ret *messagejson.Notification) {
	messagejson.NewNotification = func(*string) (*messagejson.Notification, error) {
		return ret, nil
	}
}

func mockMysqlNewNotification(ret *mysql.Notification) {
	mysql.NewNotification = func(*messagejson.Notification) (*mysql.Notification, error) {
		return ret, nil
	}
}

func mockDynamoNewNotification(ret *dynamodb.PutItemInput) {
	dynamo.NewNotification = func(*messagejson.Notification) (*dynamodb.PutItemInput, error) {
		return ret, nil
	}
}

func mockMessagejsonNewTimeline(ret *messagejson.Timeline) {
	messagejson.NewTimeline = func(*string) (*messagejson.Timeline, error) {
		return ret, nil
	}
}

func mockMysqlNewTimeline(ret *mysql.Timeline) {
	mysql.NewTimeline = func(*messagejson.Timeline) (*mysql.Timeline, error) {
		return ret, nil
	}
}

func mockDynamoNewTimeline(ret *dynamodb.PutItemInput) {
	dynamo.NewTimeline = func(*messagejson.Timeline) (*dynamodb.PutItemInput, error) {
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

	// When insertIntoDynamoDB fails
	func() {
		var arg *dynamodb.PutItemInput
		var mjNewsFeed = &messagejson.NewsFeed{}
		var msNewsFeed = &mysql.NewsFeed{}
		var dyPutItemInput = &dynamodb.PutItemInput{}

		mockMessagejsonNewNewsFeed(mjNewsFeed)
		mockMysqlNewNewsFeed(msNewsFeed)
		mockDynamoNewNewsFeed(dyPutItemInput)
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

	// When insertIntoMysql fails
	func() {
		var arg mysql.Sharder
		var mjNewsFeed = &messagejson.NewsFeed{}
		var msNewsFeed = &mysql.NewsFeed{}
		var dyPutItemInput = &dynamodb.PutItemInput{}

		mockMessagejsonNewNewsFeed(mjNewsFeed)
		mockMysqlNewNewsFeed(msNewsFeed)
		mockDynamoNewNewsFeed(dyPutItemInput)
		insertIntoDynamoDB = func(pii *dynamodb.PutItemInput) error {
			return nil
		}
		insertIntoMysql = func(nf mysql.Sharder) error {
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
}

func TestConsumeNotificationMessage(t *testing.T) {
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
		messagejson.NewNotification = func(s *string) (*messagejson.Notification, error) {
			arg = s
			return nil, errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		message.Body = aws.String("Body")

		if consumeNotificationMessage(&message) == nil {
			t.Fatal("Expect consumeNotificationMessage() to return error")
		}

		if *arg != "Body" {
			t.Fatal("messagejson.NewNotification recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When mysql.NewNotification fails
	func() {
		var arg *messagejson.Notification
		var newsFeed = &messagejson.Notification{}

		mockMessagejsonNewNotification(newsFeed)
		mysql.NewNotification = func(nf *messagejson.Notification) (*mysql.Notification, error) {
			arg = nf
			return nil, errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeNotificationMessage(&message) == nil {
			t.Fatal("Expect consumeNotificationMessage() to return error")
		}

		if arg != newsFeed {
			t.Fatal("mysql.NewNotification recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When dynamo.NewNotification fails
	func() {
		var arg *messagejson.Notification
		var mjNotification = &messagejson.Notification{}
		var msNotification = &mysql.Notification{}

		mockMessagejsonNewNotification(mjNotification)
		mockMysqlNewNotification(msNotification)
		dynamo.NewNotification = func(nf *messagejson.Notification) (*dynamodb.PutItemInput, error) {
			arg = nf
			return nil, errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeNotificationMessage(&message) == nil {
			t.Fatal("Expect consumeNotificationMessage() to return error")
		}

		if arg != mjNotification {
			t.Fatal("dynamo.NewNotification recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When insertIntoDynamoDB fails
	func() {
		var arg *dynamodb.PutItemInput
		var mjNotification = &messagejson.Notification{}
		var msNotification = &mysql.Notification{}
		var dyPutItemInput = &dynamodb.PutItemInput{}

		mockMessagejsonNewNotification(mjNotification)
		mockMysqlNewNotification(msNotification)
		mockDynamoNewNotification(dyPutItemInput)
		insertIntoDynamoDB = func(pii *dynamodb.PutItemInput) error {
			arg = pii
			return errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeNotificationMessage(&message) == nil {
			t.Fatal("Expect consumeNotificationMessage() to return error")
		}

		if arg != dyPutItemInput {
			t.Fatal("insertIntoDynamoDB recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When insertIntoMysql fails
	func() {
		var arg mysql.Sharder
		var mjNotification = &messagejson.Notification{}
		var msNotification = &mysql.Notification{}
		var dyPutItemInput = &dynamodb.PutItemInput{}

		mockMessagejsonNewNotification(mjNotification)
		mockMysqlNewNotification(msNotification)
		mockDynamoNewNotification(dyPutItemInput)
		insertIntoDynamoDB = func(pii *dynamodb.PutItemInput) error {
			return nil
		}
		insertIntoMysql = func(nf mysql.Sharder) error {
			arg = nf
			return errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeNotificationMessage(&message) == nil {
			t.Fatal("Expect consumeNotificationMessage() to return error")
		}

		if arg != msNotification {
			t.Fatal("insertIntoMysql recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()
}

func TestConsumeTimelineMessage(t *testing.T) {
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
		messagejson.NewTimeline = func(s *string) (*messagejson.Timeline, error) {
			arg = s
			return nil, errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		message.Body = aws.String("Body")

		if consumeTimelineMessage(&message) == nil {
			t.Fatal("Expect consumeTimelineMessage() to return error")
		}

		if *arg != "Body" {
			t.Fatal("messagejson.NewTimeline recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When mysql.NewTimeline fails
	func() {
		var arg *messagejson.Timeline
		var newsFeed = &messagejson.Timeline{}

		mockMessagejsonNewTimeline(newsFeed)
		mysql.NewTimeline = func(nf *messagejson.Timeline) (*mysql.Timeline, error) {
			arg = nf
			return nil, errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeTimelineMessage(&message) == nil {
			t.Fatal("Expect consumeTimelineMessage() to return error")
		}

		if arg != newsFeed {
			t.Fatal("mysql.NewTimeline recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When dynamo.NewTimeline fails
	func() {
		var arg *messagejson.Timeline
		var mjTimeline = &messagejson.Timeline{}
		var msTimeline = &mysql.Timeline{}

		mockMessagejsonNewTimeline(mjTimeline)
		mockMysqlNewTimeline(msTimeline)
		dynamo.NewTimeline = func(nf *messagejson.Timeline) (*dynamodb.PutItemInput, error) {
			arg = nf
			return nil, errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeTimelineMessage(&message) == nil {
			t.Fatal("Expect consumeTimelineMessage() to return error")
		}

		if arg != mjTimeline {
			t.Fatal("dynamo.NewTimeline recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When insertIntoDynamoDB fails
	func() {
		var arg *dynamodb.PutItemInput
		var mjTimeline = &messagejson.Timeline{}
		var msTimeline = &mysql.Timeline{}
		var dyPutItemInput = &dynamodb.PutItemInput{}

		mockMessagejsonNewTimeline(mjTimeline)
		mockMysqlNewTimeline(msTimeline)
		mockDynamoNewTimeline(dyPutItemInput)
		insertIntoDynamoDB = func(pii *dynamodb.PutItemInput) error {
			arg = pii
			return errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeTimelineMessage(&message) == nil {
			t.Fatal("Expect consumeTimelineMessage() to return error")
		}

		if arg != dyPutItemInput {
			t.Fatal("insertIntoDynamoDB recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()

	// When insertIntoMysql fails
	func() {
		var arg mysql.Sharder
		var mjTimeline = &messagejson.Timeline{}
		var msTimeline = &mysql.Timeline{}
		var dyPutItemInput = &dynamodb.PutItemInput{}

		mockMessagejsonNewTimeline(mjTimeline)
		mockMysqlNewTimeline(msTimeline)
		mockDynamoNewTimeline(dyPutItemInput)
		insertIntoDynamoDB = func(pii *dynamodb.PutItemInput) error {
			return nil
		}
		insertIntoMysql = func(nf mysql.Sharder) error {
			arg = nf
			return errors.New("")
		}
		defer restoreFunctions()

		callDeleteMessage = false
		if consumeTimelineMessage(&message) == nil {
			t.Fatal("Expect consumeTimelineMessage() to return error")
		}

		if arg != msTimeline {
			t.Fatal("insertIntoMysql recevied wrong parameter")
		}

		if !callDeleteMessage {
			t.Fatal("Expect sqs.DeleteMessage to be called")
		}
	}()
}
