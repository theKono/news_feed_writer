// +build unit

package messagejson

import (
	"fmt"
	"testing"
)

func TestSocialFeed_ValidateID(t *testing.T) {
	// When ID is not zero
	sf := &SocialFeed{ID: 1}
	if err := sf.ValidateID(); err != nil {
		t.Fatal("Expect ValidateID() not to be an error\n", err)
	}

	// When ID is zero
	sf.ID = 0
	if err := sf.ValidateID(); err == nil {
		t.Fatal("Expect ValidateID() to be an error")
	}
}

func TestSocialFeed_ValidateUserID(t *testing.T) {
	// When UserID is not zero
	sf := &SocialFeed{UserID: 1}
	if err := sf.ValidateUserID(); err != nil {
		t.Fatal("Expect ValidateUserID() not to be an error\n", err)
	}

	// When UserID is zero
	sf.UserID = 0
	if err := sf.ValidateUserID(); err == nil {
		t.Fatal("Expect ValidateUserID() to be an error")
	}
}

func TestSocialFeed_ValidateSummary(t *testing.T) {
	// When Summary is not a JSON
	sf := &SocialFeed{Summary: "[[]"}
	if err := sf.ValidateSummary(); err == nil {
		t.Fatal("Expect ValidateSummary() to be an error")
	}

	// When Summary["id"] does not exist
	sf.Summary = "{}"
	if err := sf.ValidateSummary(); err == nil {
		t.Fatal("Expect ValidateSummary() to be an error")
	}

	// When Summary["id"] does not equal string(NewsfeedID)
	sf.ID = 1
	sf.Summary = `{"id": 1}`
	if err := sf.ValidateSummary(); err == nil {
		t.Fatal("Expect sf.ValidateSummary() to be an error")
	}

	sf.Summary = `{"id": "1"}`
	if err := sf.ValidateSummary(); err != nil {
		t.Fatal("Expect ValidateSummary() not to be an error\n", err)
	}
}

func TestSocialFeed_GenerateID(t *testing.T) {
	// When summary is a bad JSON
	sf := &SocialFeed{Summary: "["}
	if sf.GenerateID() == nil {
		t.Fatal("Expect GenerateID() to return error")
	}

	sf = &SocialFeed{Summary: `{"key": 1000000000}`}
	sf.GenerateID()

	if sf.Summary != fmt.Sprintf(`{"id":"%v","key":1000000000}`, sf.ID) {
		t.Fatal("Summary is bad\n", sf.Summary)
	}
}
