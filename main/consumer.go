package main

import (
	"fmt"
	"log"
)

type Consumer interface {
	IsConsuming() bool
	Consume(chan Message)
	GetName() *string
}

type SqsConsumer struct {
	number    int
	isRunning bool
	NewsFeedWriter
}

func (sc *SqsConsumer) IsConsuming() bool {
	return sc.isRunning
}

func (sc *SqsConsumer) Consume(messageChannel chan Message) {
	sc.isRunning = true

	for message := range messageChannel {
		log.Println(*sc.GetName(), "Receive message", *message.GetBody())
		sc.Run(message)
	}

	sc.isRunning = false
}

func (sc *SqsConsumer) GetName() *string {
	ret := fmt.Sprintf("Worker-%v", sc.number)
	return &ret
}
