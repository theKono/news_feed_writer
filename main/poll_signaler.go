package main

import (
	"os"
	"os/signal"
)

type PollSignaler interface {
	PollSignal(chan os.Signal, Notifier)
}

type SimplePollSignaler struct{}

func (sc *SimplePollSignaler) PollSignal(
	signalChannel chan os.Signal,
	observer Notifier,
) {
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		observer.Notify()
	}()
}
