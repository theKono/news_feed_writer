package consumer

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"
)

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
