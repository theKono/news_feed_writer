// +build unit

package mysql

import (
	"math/rand"
	"testing"

	"github.com/theKono/orchid/model/messagejson"
)

func TestTimeline_Shard(t *testing.T) {
	n := &Timeline{}

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

func TestNewTimeline(t *testing.T) {
	input := &messagejson.Timeline{
		messagejson.SocialFeed{
			ID:      rand.Int63(),
			UserID:  rand.Int31(),
			Summary: "s",
		},
	}
	output, _ := NewTimeline(input)
	table := [][]interface{}{
		[]interface{}{"TimelineID", input.ID, output.TimelineID},
		[]interface{}{"Kid", input.UserID, output.Kid},
		[]interface{}{"Summary", input.Summary, output.Summary},
		[]interface{}{"InScrapbook", false, output.InScrapbook},
	}

	for _, tuple := range table {
		if tuple[1] != tuple[2] {
			t.Fatalf("Expect %v to equal `%v`, but got `%v`", tuple...)
		}
	}
}
