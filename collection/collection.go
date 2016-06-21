package collection

import (
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/theKono/orchid/consumer"
)

// ConsumerCollection controls the MessageConsumers
type ConsumerCollection struct {
	Consumers []consumer.MessageConsumer
	wg        *sync.WaitGroup
}

// Start makes all consumers start to work. They will compete for
// sqs.Message from the input channel.
func (cc *ConsumerCollection) Start(input <-chan *sqs.Message) error {
	if len(cc.Consumers) == 0 {
		return errors.New("No consumer")
	}

	for _, consumer := range cc.Consumers {
		go consumer.Consume(input, cc.wg)
	}

	return nil
}

// Wait blocks until all consumers stop.
func (cc *ConsumerCollection) Wait() {
	cc.wg.Wait()
}

// New creates a ConsumerCollection.
func New(c []consumer.MessageConsumer) *ConsumerCollection {
	return &ConsumerCollection{Consumers: c, wg: new(sync.WaitGroup)}
}
