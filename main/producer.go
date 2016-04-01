package main

import (
	"log"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type Producer interface {
	ShouldStop() bool
	Stop()
	Produce(chan Message)
}

type SqsProducer struct {
	shouldStop bool
	params     sqs.ReceiveMessageInput
}

func (sp *SqsProducer) ShouldStop() bool {
	return sp.shouldStop
}

func (sp *SqsProducer) Stop() {
	sp.shouldStop = true
}

func (sp *SqsProducer) Notify() {
	sp.Stop()
}

func (sp *SqsProducer) Produce(output chan Message) {
	params := InitSqsReceiveMessageParams()

	for !sp.ShouldStop() {
		log.Println("Read message...")

		for _, m := range pollSqs(params) {
			output <- m
		}
	}

	close(output)
}
