package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
)

func TestShouldStop(t *testing.T) {
	sp := SqsProducer{shouldStop: false}
	if sp.ShouldStop() {
		t.Error("Expect sp.ShouldStop() to be false")
	}

	sp = SqsProducer{shouldStop: true}
	if !sp.ShouldStop() {
		t.Error("Expect sp.ShouldStop() to be true")
	}
}

func TestStop(t *testing.T) {
	sp := SqsProducer{shouldStop: false}
	sp.Stop()

	if !sp.ShouldStop() {
		t.Error("Expect sp.ShouldStop() to be true")
	}
}

func TestNotify(t *testing.T) {
	sp := SqsProducer{}
	sp.Notify()

	if !sp.ShouldStop() {
		t.Error("Expect sp.ShouldStop() to be true")
	}
}

func TestProduce(t *testing.T) {
	sp := SqsProducer{}
	channel := make(chan Message)
	body, receiptHandle, queueURL := "Body", "ReceiptHandle", "QueueUrl"

	pollSqs = func(sp *SqsProducer) func(*sqs.ReceiveMessageInput) []Message {
		return func(_ *sqs.ReceiveMessageInput) []Message {
			sp.Stop()
			return []Message{
				SqsMessage{
					Body:          &body,
					ReceiptHandle: &receiptHandle,
					QueueURL:      &queueURL,
				},
			}
		}
	}(&sp)

	go sp.Produce(channel)

	m := <-channel
	specs := [][]string{
		{*m.GetBody(), body, "GetBody"},
		{*m.GetReceiptHandle(), receiptHandle, "GetReceiptHandle"},
		{*m.GetQueueUrl(), queueURL, "GetQueueUrl"},
	}

	for _, spec := range specs {
		if spec[0] != spec[1] {
			t.Errorf("Expect m.%v() to %v, got %v", spec[2], spec[1], spec[0])
		}
	}

	if _, ok := <-channel; ok {
		t.Error("Expect channel is closed")
	}
}
