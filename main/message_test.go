package main

import (
	"testing"
)

var (
	body          = "Body"
	receiptHandle = "ReceiptHandle"
	queueUrl      = "QueueUrl"
	sm            = SqsMessage{
		Body:          &body,
		ReceiptHandle: &receiptHandle,
		QueueURL:      &queueUrl,
	}
)

func TestGetBody(t *testing.T) {
	if *sm.GetBody() != body {
		t.Errorf("Expect sm.GetBody() to %v, got %v", body, sm.GetBody())
	}
}

func TestGetReceiptHandle(t *testing.T) {
	if *sm.GetReceiptHandle() != receiptHandle {
		t.Errorf(
			"Expect sm.GetReceiptHandle() to %v, got %v",
			receiptHandle,
			sm.GetReceiptHandle(),
		)
	}
}

func TestGetQueueUrl(t *testing.T) {
	if *sm.GetQueueUrl() != queueUrl {
		t.Errorf(
			"Expect sm.GetQueueUrl() to %v, got %v",
			queueUrl,
			sm.GetQueueUrl(),
		)
	}
}
