package main

import (
	"log"
	"os"
	"time"

	"github.com/theKono/news_feed_writer/config"
)

func initLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func initSimplePollSignaler(
	ps PollSignaler,
	signalChannel chan os.Signal,
	notifier Notifier,
) {
	ps.PollSignal(signalChannel, notifier)
}

func initConsumers(consumerCount int, messageChannel chan Message) []Consumer {
	if consumerCount <= 0 {
		log.Fatal("Consumer count is 0")
	}

	ret := make([]Consumer, consumerCount)

	for i := 0; i < consumerCount; i++ {
		nfw := SqsConsumer{
			number: i,
			NewsFeedWriter: &SimpleNewsfeedWriter{
				NewsfeedCreator: &SimpleNewsfeedCreator{},
			},
		}
		go nfw.Consume(messageChannel)
		ret[i] = &nfw
	}

	return ret
}

func initProducer(producer Producer, messageChannel chan Message) {
	producer.Produce(messageChannel)
}

func waitConsumers(consumers []Consumer) {
	for _, consumer := range consumers {
		for consumer.IsConsuming() {
			time.Sleep(time.Millisecond)
		}

		log.Println(*consumer.GetName(), "is down")
	}
}

func main() {
	parallel := config.GetParallel()
	messageChannel := make(chan Message, parallel*2)
	signalChannel := make(chan os.Signal)
	sps := SimplePollSignaler{}
	sqsProducer := SqsProducer{}
	sqsConsumers := initConsumers(parallel, messageChannel)

	initLogger()
	initSimplePollSignaler(&sps, signalChannel, &sqsProducer)
	initProducer(&sqsProducer, messageChannel)
	waitConsumers(sqsConsumers)
}
