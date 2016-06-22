package messagejson

import (
	"encoding/json"
	"log"
)

// NewsFeed is used to deserialize JSON message from SQS.
//
// It composites SocialFeed and inherited methods from there.
type NewsFeed struct {
	SocialFeed
}

// NewNewsFeed creates a NewsFeed instance by reading from a JSON
// string and validates the NewsFeed instance. If the validation
// fails, error will be returned.
var NewNewsFeed = func(jsonStr *string) (*NewsFeed, error) {
	nf := new(NewsFeed)

	if err := json.Unmarshal([]byte(*jsonStr), nf); err != nil {
		log.Println("Cannot parse message as JSON:\n", err)
		return nil, err
	}

	if err := nf.ValidateUserID(); err != nil {
		return nil, err
	}

	if nf.ValidateID() != nil {
		if err := nf.GenerateID(); err != nil {
			return nil, err
		}
	}

	if err := nf.ValidateSummary(); err != nil {
		return nil, err
	}

	return nf, nil
}
