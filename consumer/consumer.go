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

// DecorateConsumeFn is a decorator for consumer function.
//
// It handles the WaitGroup, receives message from the input channel.
// The consumer just focuses on the task of consuming a SQS message.
// To delete a message or not, it fully depends on the consumer itself.
func DecorateConsumeFn(f func(*sqs.Message) error) MessageHandler {
	return MessageHandler(
		func(input <-chan *sqs.Message, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()

			for m := range input {
				f(m)
			}
		},
	)
}
