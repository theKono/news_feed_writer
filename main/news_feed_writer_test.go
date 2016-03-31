package main

import (
	"errors"
	"testing"
)

func TestCreateNewsfeed(t *testing.T) {
	var jsonString string
	var n *Newsfeed
	var err error

	// When it is a bad json
	snc := SimpleNewsfeedCreator{}
	jsonString = "{}}"
	if n, err = snc.CreateNewsfeed(&jsonString); err == nil {
		t.Errorf("Expect snc.CreateNewsfeed to raise error")
	}

	// When it is a good json
	jsonString = `{
        "newsfeed_id": 1,
        "kid": 2,
        "observable_type": "ot",
        "observable_id": "oi",
        "event_type": "e",
        "target_type": "tt",
        "target_id": "ti",
        "summary": "s"
    }`
	if n, err = snc.CreateNewsfeed(&jsonString); err != nil {
		t.Errorf("Expect snc.CreateNewsfeed not to raise error")
	}

	if n.NewsfeedID != 1 {
		t.Errorf("Expect n.NewsfeedId to 1, got %v", n.NewsfeedID)
	}

	if n.Kid != 2 {
		t.Errorf("Expect n.Kid to 2, got %v", n.Kid)
	}

	if n.ObservableType != "ot" {
		t.Errorf(`Expect n.ObservableType to "ot", got %q`, n.ObservableType)
	}

	if n.ObservableID != "oi" {
		t.Errorf(`Expect n.ObservableID to "oi", got %q`, n.ObservableID)
	}

	if n.EventType != "e" {
		t.Errorf(`Expect n.EventType to "e", got %q`, n.EventType)
	}

	if n.TargetType != "tt" {
		t.Errorf(`Expect n.TargetType to "tt", got %q`, n.TargetType)
	}

	if n.TargetID != "ti" {
		t.Errorf(`Expect n.TargetID to "ti", got %q`, n.TargetID)
	}

	if n.Summary != "s" {
		t.Errorf(`Expect n.Summary to "s", got %q`, n.Summary)
	}
}

type MockNewsfeedCreator struct{}

func (mnc *MockNewsfeedCreator) CreateNewsfeed(s *string) (*Newsfeed, error) {
	n := Newsfeed{NewsfeedID: 1, Kid: 1, Summary: `{"id":"1"}`}

	if *s != "" {
		return &n, errors.New(*s)
	}

	return &n, nil
}

type MockMessage struct {
	Body          string
	ReceiptHandle string
	QueueURL      string
}

func (mm *MockMessage) GetBody() *string {
	return &mm.Body
}

func (mm *MockMessage) GetReceiptHandle() *string {
	return &mm.ReceiptHandle
}

func (mm *MockMessage) GetQueueUrl() *string {
	return &mm.QueueURL
}

func TestRunWhenItIsABadJson(t *testing.T) {
	defer nukeDb()

	callDeleteMessage := false
	nukeDb()

	deleteMessage = func(m Message) {
		callDeleteMessage = true
	}

	m := MockMessage{Body: "b", ReceiptHandle: "r", QueueURL: "q"}
	snw := SimpleNewsfeedWriter{NewsfeedCreator: &MockNewsfeedCreator{}}
	snw.Run(&m)

	if !callDeleteMessage {
		t.Errorf("Expect deleteMessage is called")
	}
}

func TestRunWhenSaveIsFailed(t *testing.T) {
	defer nukeDb()

	callDeleteMessage := false
	nukeDb()

	deleteMessage = func(m Message) {
		callDeleteMessage = true
	}

	empty := ""
	mn := MockNewsfeedCreator{}
	n, _ := mn.CreateNewsfeed(&empty)
	n.Save()

	m := MockMessage{Body: "", ReceiptHandle: "r", QueueURL: "q"}
	snw := SimpleNewsfeedWriter{NewsfeedCreator: &MockNewsfeedCreator{}}
	snw.Run(&m)

	if !callDeleteMessage {
		t.Errorf("Expect deleteMessage is called")
	}
}

func nukeDb() {
	shards[0].Delete(&Newsfeed{})
	shards[1].Delete(&Newsfeed{})
}
