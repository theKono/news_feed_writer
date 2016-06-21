package messagejson

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/theKono/orchid/util"
)

// NewsFeed is used to deserialize JSON message from SQS.
type NewsFeed struct {
	ID             int64  `json:"id"`
	UserID         int32  `json:"user_id"`
	ObservableType string `json:"observable_type"`
	ObservableID   string `json:"observable_id"`
	EventType      string `json:"event_type"`
	TargetType     string `json:"target_type"`
	TargetID       string `json:"target_id"`
	Summary        string `json:"summary"`
}

// ValidateID validates NewsFeed.ID is present.
func (nf *NewsFeed) ValidateID() error {
	if nf.ID != 0 {
		return nil
	}

	return errors.New("ID is required")
}

// GenerateID generates a new NewsFeed.ID and sync NewsFeed.Summary.
func (nf *NewsFeed) GenerateID() error {
	var summary map[string]interface{}
	var jsonStr []byte
	var err error

	if err = json.Unmarshal([]byte(nf.Summary), &summary); err != nil {
		return errors.New("Summary is a bad json")
	}

	nf.ID = util.GenerateID(int(nf.UserID % 2))

	summary["id"] = fmt.Sprint(nf.ID)
	if jsonStr, err = json.Marshal(summary); err != nil {
		return errors.New("Cannot marshal to json string")
	}
	nf.Summary = string(jsonStr)

	return nil
}

// ValidateUserID validates NewsFeed.UserID is present.
func (nf *NewsFeed) ValidateUserID() error {
	if nf.UserID != 0 {
		return nil
	}

	return errors.New("UserID is required")
}

// ValidateSummary validates NewsFeed.Summary.
// 1) NewsFeed.Summary should be a valid JSON string.
// 2) NewsFeed.Summary must have "id" included after deserializing.
// 3) The former "id" must be a stringified NewsFeed.ID.
func (nf *NewsFeed) ValidateSummary() error {
	var summary map[string]interface{}

	if err := json.Unmarshal([]byte(nf.Summary), &summary); err != nil {
		return errors.New("Summary is a bad json")
	}

	if _, ok := summary["id"]; !ok {
		return errors.New(`Summary["id"] is required`)
	}

	if summary["id"] != fmt.Sprint(nf.ID) {
		return errors.New(`Summary["id"] is not equal to ID`)
	}

	return nil
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
