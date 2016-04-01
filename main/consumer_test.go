package main

import (
	"testing"
	"time"
)

func TestIsConsuming(t *testing.T) {
	sc := SqsConsumer{isRunning: true}
	if !sc.IsConsuming() {
		t.Error("Expect sc.IsConsuming() to be true")
	}

	sc = SqsConsumer{isRunning: false}
	if sc.IsConsuming() {
		t.Error("Expect sc.IsConsuming() to be false")
	}
}

func TestGetName(t *testing.T) {
	sc := SqsConsumer{number: 1}

	if *sc.GetName() != "Worker-1" {
		t.Errorf(
			"Expect sc.GetName() to %q, got %q", "Worker-1", sc.GetName(),
		)
	}
}

type MockNewsFeedWriter struct {
	called bool
	m      Message
}

func (mnfw *MockNewsFeedWriter) Run(m Message) {
	mnfw.called = true
	mnfw.m = m
}

func TestConsume(t *testing.T) {
	channel := make(chan Message)
	mnfw := MockNewsFeedWriter{}
	sc := SqsConsumer{NewsFeedWriter: &mnfw}

	go sc.Consume(channel)

	b, r, q := "B", "R", "Q"
	channel <- SqsMessage{Body: &b, ReceiptHandle: &r, QueueURL: &q}
	close(channel)
	// Let sc has time to cleanup in go routine
	time.Sleep(time.Millisecond)

	if sc.IsConsuming() {
		t.Errorf("Expect sc.IsConsuming() to false")
	}

	if !mnfw.called {
		t.Errorf("Expect mnfw.called to true")
	}

	specs := [][]string{
		[]string{*mnfw.m.GetBody(), b, "GetBody"},
		[]string{*mnfw.m.GetReceiptHandle(), r, "GetReceiptHandle"},
		[]string{*mnfw.m.GetQueueUrl(), q, "GetQueueUrl"},
	}
	for _, spec := range specs {
		if spec[0] != spec[1] {
			t.Errorf(
				"Expect mnfw.m.%v() to %q, got %q", spec[2], spec[1], spec[0],
			)
		}
	}
}
