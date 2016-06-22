package messagejson

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/theKono/orchid/util"
)

// SocialFeed acts as the super class of NewsFeed and Notification
type SocialFeed struct {
	ID             int64  `json:"id"`
	UserID         int32  `json:"user_id"`
	ObservableType string `json:"observable_type"`
	ObservableID   string `json:"observable_id"`
	EventType      string `json:"event_type"`
	TargetType     string `json:"target_type"`
	TargetID       string `json:"target_id"`
	Summary        string `json:"summary"`
}

// ValidateID validates SocialFeed.ID is present.
func (sf *SocialFeed) ValidateID() error {
	if sf.ID != 0 {
		return nil
	}

	return errors.New("ID is required")
}

// GenerateID generates a new SocialFeed.ID and sync SocialFeed.Summary.
func (sf *SocialFeed) GenerateID() error {
	var summary map[string]interface{}
	var jsonStr []byte
	var err error

	if err = json.Unmarshal([]byte(sf.Summary), &summary); err != nil {
		return errors.New("Summary is a bad json")
	}

	sf.ID = util.GenerateID(int(sf.UserID % 2))

	summary["id"] = fmt.Sprint(sf.ID)
	if jsonStr, err = json.Marshal(summary); err != nil {
		return errors.New("Cannot marshal to json string")
	}
	sf.Summary = string(jsonStr)

	return nil
}

// ValidateUserID validates SocialFeed.UserID is present.
func (sf *SocialFeed) ValidateUserID() error {
	if sf.UserID != 0 {
		return nil
	}

	return errors.New("UserID is required")
}

// ValidateSummary validates SocialFeed.Summary.
// 1) SocialFeed.Summary should be a valid JSON string.
// 2) SocialFeed.Summary must have "id" included after deserializing.
// 3) The former "id" must be a stringified SocialFeed.ID.
func (sf *SocialFeed) ValidateSummary() error {
	var summary map[string]interface{}

	if err := json.Unmarshal([]byte(sf.Summary), &summary); err != nil {
		return errors.New("Summary is a bad json")
	}

	if _, ok := summary["id"]; !ok {
		return errors.New(`Summary["id"] is required`)
	}

	if summary["id"] != fmt.Sprint(sf.ID) {
		return errors.New(`Summary["id"] is not equal to ID`)
	}

	return nil
}
