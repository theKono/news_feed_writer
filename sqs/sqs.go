package sqs

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsSqs "github.com/aws/aws-sdk-go/service/sqs"

	"github.com/theKono/orchid/cfg"
)

// DefaultNumberOfMessages is the default value for max number of
// messages to return.
const DefaultNumberOfMessages = 1

// DefaultWaitTimeSeconds is the default value for duration for which
// the call will wait for a message to arrive.
const DefaultWaitTimeSeconds = 20

// service is the SQS client
var service *awsSqs.SQS

// queueURL is the URL of SQS queue
var queueURL string

// keepPolling is a loop condition which decides whether or not to
// keep recieving message from SQS
var keepPolling bool

// wg is used by wait(), registerWaitGroup(), done()
var wg *sync.WaitGroup

// New initializes package variables.
func init() {
	if cfg.SqsRegion == "" {
		log.Fatal("SQS region is required")
	}
	if cfg.SqsQueueURL == "" {
		log.Fatal("SQS Queue URL is required")
	}

	service = awsSqs.New(session.New(), aws.NewConfig().WithRegion(cfg.SqsRegion))
	queueURL = cfg.SqsQueueURL
	keepPolling = false
	wg = new(sync.WaitGroup)
}

// Start makes it start to poll SQS.
func Start(output chan<- *awsSqs.Message) error {
	keepPolling = true
	for i := 0; i < cfg.Parallel; i++ {
		go poll(output)
	}
	return nil
}

// Stop stops polling SQS.
func Stop() error {
	keepPolling = false
	wait()
	return nil
}

// DeleteMessage deletes a sqs.Message
var DeleteMessage = func(message *awsSqs.Message) (err error) {
	_, err = service.DeleteMessage(
		&awsSqs.DeleteMessageInput{
			QueueUrl:      &queueURL,
			ReceiptHandle: message.ReceiptHandle,
		},
	)
	return
}

// poll repeatedly waits and receives messages from SQS.
// It is the main go routine.
var poll = func(output chan<- *awsSqs.Message) {
	registerWaitGroup()
	defer done()

	params := makeReceiveMessageInput()

	for keepPolling {
		resp, err := receiveMessage(params)

		if err != nil {
			log.Printf("Fail to receive message from SQS\nparams: %v\nError: %v\n", params, err.Error())
			continue
		}

		populate(resp.Messages, output)
	}
}

// makeReceiveMessageInput makes the parameter to poll SQS.
var makeReceiveMessageInput = func() *awsSqs.ReceiveMessageInput {
	return &awsSqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(DefaultNumberOfMessages),
		WaitTimeSeconds:     aws.Int64(DefaultWaitTimeSeconds),
	}
}

// receiveMessage will wait and receive messages from SQS.
// Need unit testing
var receiveMessage = func(params *awsSqs.ReceiveMessageInput) (*awsSqs.ReceiveMessageOutput, error) {
	return service.ReceiveMessage(params)
}

// populate populates messages into output channel.
var populate = func(messages []*awsSqs.Message, output chan<- *awsSqs.Message) {
	for _, message := range messages {
		output <- message
	}
}

// wait is used to block until all goroutines in wait group are done.
var wait = func() {
	wg.Wait()
}

// registerWaitGroup is used by a goroutine that wants to be a member of wait group.
var registerWaitGroup = func() {
	wg.Add(1)
}

// done is used by a goroutine that disclaims it is done.
var done = func() {
	wg.Done()
}
