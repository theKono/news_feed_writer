// +build timeline

package messagejson

import (
	"encoding/json"
	"log"
)

// Timeline is used to deserialize JSON message from SQS.
//
// It composites SocialFeed and inherited methods from there.
type Timeline struct {
	SocialFeed
}

// NewTimeline creates a Timeline instance by reading from a JSON
// string and validates the Timeline instance. If the validation
// fails, error will be returned.
var NewTimeline = func(jsonStr *string) (*Timeline, error) {
	nf := new(Timeline)

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
