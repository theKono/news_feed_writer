package main

import (
	"os"
	"testing"
)

type MockNotifier struct {
	called bool
}

func (mn *MockNotifier) Notify() {
	mn.called = true
}

func TestPollSignal(t *testing.T) {
	channel := make(chan os.Signal)
	mockNotifier := MockNotifier{}
	sps := SimplePollSignaler{}

	sps.PollSignal(channel, &mockNotifier)
	if mockNotifier.called {
		t.Error("Expect mockNotifier.called to be false")
	}

	channel <- os.Interrupt
	if !mockNotifier.called {
		t.Error("Expect mockNotifier.called to be true")
	}
}
