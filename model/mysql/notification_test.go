// +build unit

package mysql

import (
	"math/rand"
	"testing"

	"github.com/theKono/orchid/model/messagejson"
)

func TestNotification_Shard(t *testing.T) {
	n := &Notification{}

	// When Kid is an even
	n.Kid = 2
	if n.Shard() != 0 {
		t.Fatalf("Expect Shard() to be 0, but got `%v`", n.Shard())
	}

	n.Kid = 1
	if n.Shard() != 1 {
		t.Fatalf("Expect Shard() to be 1, but got `%v`", n.Shard())
	}
}

func TestNewNotification(t *testing.T) {
	input := &messagejson.Notification{
		messagejson.SocialFeed{
			ID:             rand.Int63(),
			UserID:         rand.Int31(),
			ObservableType: "ot",
			ObservableID:   "oi",
			EventType:      "et",
			TargetType:     "tt",
			TargetID:       "ti",
			Summary:        "s",
		},
	}
	output, _ := NewNotification(input)
	table := [][]interface{}{
		[]interface{}{"NotificationID", input.ID, output.NotificationID},
		[]interface{}{"Kid", input.UserID, output.Kid},
		[]interface{}{"ObservableType", input.ObservableType, output.ObservableType},
		[]interface{}{"ObservableID", input.ObservableID, output.ObservableID},
		[]interface{}{"EventType", input.EventType, output.EventType},
		[]interface{}{"TargetType", input.TargetType, output.TargetType},
		[]interface{}{"TargetID", input.TargetID, output.TargetID},
		[]interface{}{"Summary", input.Summary, output.Summary},
		[]interface{}{"State", int8(UnseenAndUnread), output.State},
	}

	for _, tuple := range table {
		if tuple[1] != tuple[2] {
			t.Fatalf("Expect %v to equal `%v`, but got `%v`", tuple...)
		}
	}
}
