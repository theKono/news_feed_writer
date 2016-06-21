package consumer

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"
)

// MessageConsumer is an interface defining operations that a basic SQS
// message consumer should implement.
type MessageConsumer interface {
	Consume(<-chan *sqs.Message, *sync.WaitGroup)
}

// MessageHandler is a handy type to make a function a message consumer.
type MessageHandler func(<-chan *sqs.Message, *sync.WaitGroup)

// Consume implements the MessageConsumer interface.
func (mh MessageHandler) Consume(c <-chan *sqs.Message, wg *sync.WaitGroup) {
	mh(c, wg)
}
