// +build unit

package messagejson

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTimeline_ValidateID(t *testing.T) {
	// When ID is not zero
	nf := &Timeline{SocialFeed{ID: 1}}
	if err := nf.ValidateID(); err != nil {
		t.Fatal("Expect ValidateID() not to be an error\n", err)
	}

	// When ID is zero
	nf.ID = 0
	if err := nf.ValidateID(); err == nil {
		t.Fatal("Expect ValidateID() to be an error")
	}
}

func TestTimeline_ValidateUserID(t *testing.T) {
	// When UserID is not zero
	nf := &Timeline{SocialFeed{UserID: 1}}
	if err := nf.ValidateUserID(); err != nil {
		t.Fatal("Expect ValidateUserID() not to be an error\n", err)
	}

	// When UserID is zero
	nf.UserID = 0
	if err := nf.ValidateUserID(); err == nil {
		t.Fatal("Expect ValidateUserID() to be an error")
	}
}

func TestTimeline_ValidateSummary(t *testing.T) {
	// When Summary is not a JSON
	nf := &Timeline{SocialFeed{Summary: "[[]"}}
	if err := nf.ValidateSummary(); err == nil {
		t.Fatal("Expect ValidateSummary() to be an error")
	}

	// When Summary["id"] does not exist
	nf.Summary = "{}"
	if err := nf.ValidateSummary(); err == nil {
		t.Fatal("Expect ValidateSummary() to be an error")
	}

	// When Summary["id"] does not equal string(NewsfeedID)
	nf.ID = 1
	nf.Summary = `{"id": 1}`
	if err := nf.ValidateSummary(); err == nil {
		t.Fatal("Expect nf.ValidateSummary() to be an error")
	}

	nf.Summary = `{"id": "1"}`
	if err := nf.ValidateSummary(); err != nil {
		t.Fatal("Expect ValidateSummary() not to be an error\n", err)
	}
}

func TestTimeline_GenerateID(t *testing.T) {
	// When summary is a bad JSON
	nf := &Timeline{SocialFeed{Summary: "["}}
	if nf.GenerateID() == nil {
		t.Fatal("Expect GenerateID() to return error")
	}

	nf = &Timeline{SocialFeed{Summary: `{"key": 1000000000}`}}
	nf.GenerateID()

	if nf.Summary != fmt.Sprintf(`{"id":"%v","key":1000000000}`, nf.ID) {
		t.Fatal("Summary is bad\n", nf.Summary)
	}
}

func TestNewTimeline(t *testing.T) {
	// It basically parse the JSON string
	m := map[string]interface{}{
		"id":      1,
		"user_id": 2,
		"summary": `{"id":"1"}`,
	}
	b, _ := json.Marshal(m)
	jsonStr := string(b)
	nf, err := NewTimeline(&jsonStr)

	if err != nil {
		t.Fatal("Expect NewTimeline() not to return error\n", err)
	}

	table := [][]interface{}{
		[]interface{}{"ID", m["id"], int(nf.ID)},
		[]interface{}{"UserID", m["user_id"], int(nf.UserID)},
		[]interface{}{"Summary", m["summary"], nf.Summary},
	}
	for _, tuple := range table {
		if tuple[1] != tuple[2] {
			t.Fatalf("Expect nf.%v to equal %T(%v), but got %T(%v)\n", tuple[0], tuple[1], tuple[1], tuple[2], tuple[2])
		}
	}

	// When `id` is not provided
	delete(m, "id")
	b, _ = json.Marshal(m)
	jsonStr = string(b)
	nf, err = NewTimeline(&jsonStr)

	if err != nil {
		t.Fatal("Expect NewTimeline() not to return error\n", err)
	}

	if nf.ID == 0 {
		t.Fatal("Expect ID not to equal 0")
	}

	// When `user_id` is not provided
	delete(m, "user_id")
	b, _ = json.Marshal(m)
	jsonStr = string(b)
	nf, err = NewTimeline(&jsonStr)

	if err == nil {
		t.Fatal("Expect NewTimeline() to return error")
	}

	// When `summary` is not a JSON
	m["user_id"] = 1
	m["summary"] = "[[]"
	b, _ = json.Marshal(m)
	jsonStr = string(b)
	nf, err = NewTimeline(&jsonStr)

	if err == nil {
		t.Fatal("Expect NewTimeline() to return error")
	}
}
