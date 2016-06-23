// +build unit

package messagejson

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNotification_ValidateID(t *testing.T) {
	// When ID is not zero
	n := &Notification{SocialFeed{ID: 1}}
	if err := n.ValidateID(); err != nil {
		t.Fatal("Expect ValidateID() not to be an error\n", err)
	}

	// When ID is zero
	n.ID = 0
	if err := n.ValidateID(); err == nil {
		t.Fatal("Expect ValidateID() to be an error")
	}
}

func TestNotification_ValidateUserID(t *testing.T) {
	// When UserID is not zero
	n := &Notification{SocialFeed{UserID: 1}}
	if err := n.ValidateUserID(); err != nil {
		t.Fatal("Expect ValidateUserID() not to be an error\n", err)
	}

	// When UserID is zero
	n.UserID = 0
	if err := n.ValidateUserID(); err == nil {
		t.Fatal("Expect ValidateUserID() to be an error")
	}
}

func TestNotification_ValidateSummary(t *testing.T) {
	// When Summary is not a JSON
	n := &Notification{SocialFeed{Summary: "[[]"}}
	if err := n.ValidateSummary(); err == nil {
		t.Fatal("Expect ValidateSummary() to be an error")
	}

	// When Summary["id"] does not exist
	n.Summary = "{}"
	if err := n.ValidateSummary(); err == nil {
		t.Fatal("Expect ValidateSummary() to be an error")
	}

	// When Summary["id"] does not equal string(NewsfeedID)
	n.ID = 1
	n.Summary = `{"id": 1}`
	if err := n.ValidateSummary(); err == nil {
		t.Fatal("Expect n.ValidateSummary() to be an error")
	}

	n.Summary = `{"id": "1"}`
	if err := n.ValidateSummary(); err != nil {
		t.Fatal("Expect ValidateSummary() not to be an error\n", err)
	}
}

func TestNotification_GenerateID(t *testing.T) {
	// When summary is a bad JSON
	n := &Notification{SocialFeed{Summary: "["}}
	if n.GenerateID() == nil {
		t.Fatal("Expect GenerateID() to return error")
	}

	n = &Notification{SocialFeed{Summary: `{"key": 1000000000}`}}
	n.GenerateID()

	if n.Summary != fmt.Sprintf(`{"id":"%v","key":1000000000}`, n.ID) {
		t.Fatal("Summary is bad\n", n.Summary)
	}
}

func TestNewNotification(t *testing.T) {
	// It basically parse the JSON string
	m := map[string]interface{}{
		"id":              1,
		"user_id":         2,
		"observable_type": "ot",
		"observable_id":   "oi",
		"event_type":      "et",
		"target_type":     "tt",
		"target_id":       "ti",
		"summary":         `{"id":"1"}`,
	}
	b, _ := json.Marshal(m)
	jsonStr := string(b)
	n, err := NewNotification(&jsonStr)

	if err != nil {
		t.Fatal("Expect NewNotification() not to return error\n", err)
	}

	table := [][]interface{}{
		[]interface{}{"ID", m["id"], int(n.ID)},
		[]interface{}{"UserID", m["user_id"], int(n.UserID)},
		[]interface{}{"ObservableType", m["observable_type"], n.ObservableType},
		[]interface{}{"ObservableID", m["observable_id"], n.ObservableID},
		[]interface{}{"EventType", m["event_type"], n.EventType},
		[]interface{}{"TargetType", m["target_type"], n.TargetType},
		[]interface{}{"TargetID", m["target_id"], n.TargetID},
		[]interface{}{"Summary", m["summary"], n.Summary},
	}
	for _, tuple := range table {
		if tuple[1] != tuple[2] {
			t.Fatalf("Expect n.%v to equal %T(%v), but got %T(%v)\n", tuple[0], tuple[1], tuple[1], tuple[2], tuple[2])
		}
	}

	// When `id` is not provided
	delete(m, "id")
	b, _ = json.Marshal(m)
	jsonStr = string(b)
	n, err = NewNotification(&jsonStr)

	if err != nil {
		t.Fatal("Expect NewNotification() not to return error\n", err)
	}

	if n.ID == 0 {
		t.Fatal("Expect ID not to equal 0")
	}

	// When `user_id` is not provided
	delete(m, "user_id")
	b, _ = json.Marshal(m)
	jsonStr = string(b)
	n, err = NewNotification(&jsonStr)

	if err == nil {
		t.Fatal("Expect NewNotification() to return error")
	}

	// When `summary` is not a JSON
	m["user_id"] = 1
	m["summary"] = "[[]"
	b, _ = json.Marshal(m)
	jsonStr = string(b)
	n, err = NewNotification(&jsonStr)

	if err == nil {
		t.Fatal("Expect NewNotification() to return error")
	}
}
