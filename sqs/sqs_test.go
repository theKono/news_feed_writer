// +build unit

package sqs

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsSqs "github.com/aws/aws-sdk-go/service/sqs"
)

func TestStart(t *testing.T) {
	var pollArg chan<- *awsSqs.Message

	q := make(chan *awsSqs.Message)

	// mock poll function
	origPoll := poll
	poll = func(c chan<- *awsSqs.Message) {
		pollArg = c
	}
	defer func() { poll = origPoll }()

	err := Start(q)

	// hold on until poll is finished
	time.Sleep(time.Millisecond * 10)

	if err != nil {
		t.Fatal("Expect `err` to be nil\nerr:", err.Error())
	}

	if !keepPolling {
		t.Fatal("Expect `keepPolling` to be true")
	}

	if pollArg != q {
		t.Fatal("Expect poll is called with correct argument")
	}
}

func TestStop(t *testing.T) {
	// mock wait function
	called := false
	origWait := wait
	wait = func() { called = true }
	defer func() { wait = origWait }()

	err := Stop()

	if err != nil {
		t.Fatal("Expect `err` to be nil\nerr:", err.Error())
	}

	if keepPolling {
		t.Fatal("Expect `keepPolling` to be false")
	}

	if !called {
		t.Fatal("Expect `wait` function to be called")
	}
}

func TestMakeReceiveMessageInput(t *testing.T) {
	queueURL = "qq"
	output := makeReceiveMessageInput()

	if *output.QueueUrl != "qq" {
		t.Fatalf("Expect `queueURL` to be `%v`, but got `%v`\n", queueURL, *output.QueueUrl)
	}

	if *output.MaxNumberOfMessages != DefaultNumberOfMessages {
		t.Fatalf("Expect `MaxNumberOfMessages` to be `%v`, but got `%v`\n", DefaultWaitTimeSeconds, *output.MaxNumberOfMessages)
	}

	if *output.WaitTimeSeconds != DefaultWaitTimeSeconds {
		t.Fatalf("Expect `WaitTimeSeconds` to be `%v`, but got `%v`\n", DefaultWaitTimeSeconds, *output.WaitTimeSeconds)
	}
}

func TestPopulate(t *testing.T) {
	messageCount := 10
	c := make(chan *awsSqs.Message, messageCount)
	messages := make([]*awsSqs.Message, messageCount)

	for i := 0; i < messageCount; i++ {
		messages[i] = &awsSqs.Message{MessageId: aws.String(fmt.Sprint(i))}
	}

	populate(messages, c)

	for i := 0; i < messageCount; i++ {
		m := <-c
		if *m.MessageId != fmt.Sprint(i) {
			t.Fatalf("Expect `MessageId` to be `%v`, but got `%v`", i, *m.MessageId)
		}
	}
}

func TestWaitDone(t *testing.T) {
	var worker string
	var supervisor string

	go func() {
		worker = "sleeping"
		time.Sleep(10 * time.Millisecond)
		worker = "done"
		done()
	}()

	go func() {
		supervisor = "waiting"
		wait()
		supervisor = "done"

		if worker != "done" {
			t.Fatalf("Expect `worker` to be `done`, but got `%v`", worker)
		}
	}()

	time.Sleep(20 * time.Millisecond)

	if supervisor != "done" {
		t.Fatalf("Expect `supervisor` to be `done`, but got `%v`", supervisor)
	}
}
