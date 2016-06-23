package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	awsSqs "github.com/aws/aws-sdk-go/service/sqs"

	"github.com/theKono/orchid/cfg"
	"github.com/theKono/orchid/collection"
	"github.com/theKono/orchid/consumer"
	"github.com/theKono/orchid/sqs"
)

// init intializes logging, configuration.
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// main is the entry point for the worker.
func main() {
	log.Println("notification worker starts")

	messageCh := make(chan *awsSqs.Message, cfg.Parallel)
	if err := sqs.Start(messageCh); err != nil {
		log.Fatalln("Fail to poll SQS\n", err.Error())
	}

	messageConsumers := make([]consumer.MessageConsumer, cfg.Parallel)
	for i := 0; i < cfg.Parallel; i++ {
		messageConsumers[i] = consumer.ConsumeNotification
	}

	consumers := collection.New(messageConsumers)
	if err := consumers.Start(messageCh); err != nil {
		log.Fatalln("Fail to consume SQS\n", err.Error())
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		log.Println("Shutting down...")

		if err := sqs.Stop(); err != nil {
			log.Println("Fail to stop sqs\n", err.Error())
		}

		close(messageCh)
		consumers.Wait()
	}
}
