// +build integration

package sqs

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	awsSqs "github.com/aws/aws-sdk-go/service/sqs"
)

func TestReceiveDeleteMessage(t *testing.T) {
	_, err := service.SendMessage(
		&awsSqs.SendMessageInput{
			QueueUrl:    &queueURL,
			MessageBody: aws.String("qq"),
		},
	)
	if err != nil {
		t.Fatal("Expect SendMessage not to return error\n", err)
	}

	resp, err := receiveMessage(makeReceiveMessageInput())
	if err != nil {
		t.Fatal("Expect receiveMessage not to return error\n", err)
	}
	if len(resp.Messages) != 1 {
		t.Fatal("Expect only one message\n", resp.Messages)
	}

	message := resp.Messages[0]
	if *message.Body != "qq" {
		t.Fatal("Expect message body to equal `qq`, but got `%v`", *message.Body)
	}

	err = DeleteMessage(message)
	if err != nil {
		t.Fatal("Expect DeleteMessage not to return error\n", err)
	}
}
