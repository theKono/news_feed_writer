package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/theKono/news_feed_writer/config"
)

var (
	sqsClient = InitSqsClient()
)

func InitSqsClient() *sqs.SQS {
	return sqs.New(
		session.New(),
		&aws.Config{Region: aws.String(config.GetSqsQueueRegion())},
	)
}

func InitSqsReceiveMessageParams() *sqs.ReceiveMessageInput {
	return &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(config.GetSqsQueueUrl()),
		MaxNumberOfMessages: aws.Int64(int64(config.GetMaxNumberOfMessages())),
		WaitTimeSeconds:     aws.Int64(int64(config.GetWaitTimeSeconds())),
	}
}

func InitSqsDeleteMessageParams(message Message) *sqs.DeleteMessageInput {
	return &sqs.DeleteMessageInput{
		QueueUrl:      message.GetQueueUrl(),
		ReceiptHandle: message.GetReceiptHandle(),
	}
}

var pollSqs = func(params *sqs.ReceiveMessageInput) []Message {
	resp, err := sqsClient.ReceiveMessage(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			log.Println("Error:", err)
		}

		log.Fatal(err)
	}

	ret := make([]Message, len(resp.Messages))
	for ix, m := range resp.Messages {
		ret[ix] = SqsMessage{
			Body:          m.Body,
			ReceiptHandle: m.ReceiptHandle,
			QueueURL:      params.QueueUrl,
		}
	}

	return ret
}

var deleteMessage = func(message Message) {
	_, err := sqsClient.DeleteMessage(InitSqsDeleteMessageParams(message))

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			log.Println("Error:", err)
		}
	}
}
