package messagejson

import (
	"encoding/json"
	"log"
)

// Notification is used to deserialize JSON message from SQS.
//
// It composites SocialFeed and inherited methods from there.
type Notification struct {
	SocialFeed
}

// NewNotification creates a Notification instance by reading from a
// JSON string and validates the NewsFeed instance. If the validation
// fails, error will be returned.
var NewNotification = func(jsonStr *string) (*Notification, error) {
	n := new(Notification)

	if err := json.Unmarshal([]byte(*jsonStr), n); err != nil {
		log.Println("Cannot parse message as JSON:\n", err)
		return nil, err
	}

	if err := n.ValidateUserID(); err != nil {
		return nil, err
	}

	if n.ValidateID() != nil {
		if err := n.GenerateID(); err != nil {
			return nil, err
		}
	}

	if err := n.ValidateSummary(); err != nil {
		return nil, err
	}

	return n, nil
}
