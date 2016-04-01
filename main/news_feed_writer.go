package main

import (
	"encoding/json"
	"log"
)

type NewsFeedWriter interface {
	Run(Message)
}

type NewsfeedCreator interface {
	CreateNewsfeed(*string) (*Newsfeed, error)
}

type SimpleNewsfeedWriter struct {
	NewsfeedCreator
}

func (snfw *SimpleNewsfeedWriter) Run(message Message) {
	defer deleteMessage(message)

	n, err := snfw.CreateNewsfeed(message.GetBody())
	if err != nil {
		log.Println(
			"Error",
			"Cannot create Newsfeed from message, just delete it",
		)
		return
	}

	if err = n.Save(); err != nil {
		log.Println("Error", err, "Just delete it")
		return
	}
}

type SimpleNewsfeedCreator struct{}

func (snw *SimpleNewsfeedCreator) CreateNewsfeed(data *string) (*Newsfeed, error) {
	var nj NewsfeedJson

	if err := json.Unmarshal([]byte(*data), &nj); err != nil {
		log.Println("Cannot parse message as json")
		log.Println(data)
		return &Newsfeed{}, err
	}

	return nj.NewNewsfeed(), nil
}
