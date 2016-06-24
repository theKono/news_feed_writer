// +build unit

package collection

import (
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/theKono/orchid/consumer"
)

type SimpleConsumer struct{}

func (sc SimpleConsumer) Consume(<-chan *sqs.Message, *sync.WaitGroup) {
	return
}

func TestNew(t *testing.T) {
	consumers := []consumer.MessageConsumer{SimpleConsumer{}}
	c := New(consumers)

	if !(len(c.Consumers) == len(consumers) && c.Consumers[0] == consumers[0]) {
		t.Fatal("Expect Consumers to eq `%v`, but got `%v`", consumers, c.Consumers)
	}
}

func TestStart(t *testing.T) {
	c := make(chan *sqs.Message)

	// When there is no consumer
	collection := New([]consumer.MessageConsumer{})
	if err := collection.Start(c); err == nil {
		t.Fatal("Expect Start to return error")
	}

	// When there is consumer
	called := false
	var mh = func(<-chan *sqs.Message, *sync.WaitGroup) {
		called = true
	}
	collection = New([]consumer.MessageConsumer{consumer.MessageHandler(mh)})
	collection.Start(c)
	time.Sleep(time.Millisecond)
	if !called {
		t.Fatal("Expect mh is called")
	}
}

func TestWait(t *testing.T) {
	collection := new(ConsumerCollection)
	collection.wg = new(sync.WaitGroup)

	blocking := 0

	go func() {
		collection.wg.Add(1)
		blocking = 1
		time.Sleep(5 * time.Millisecond)
		collection.wg.Done()
		blocking = 2
	}()

	go func() {
		time.Sleep(time.Millisecond) // Let the first goroutine starts first
		collection.Wait()
		if blocking != 2 {
			t.Fatalf("Expect blocking to be 2, but got `%v`", blocking)
		}
	}()

	// Wait until 2 goroutines finish
	time.Sleep(10 * time.Millisecond)
}
