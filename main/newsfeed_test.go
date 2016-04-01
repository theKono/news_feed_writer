package main

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestGenerateId(t *testing.T) {
	id := GenerateId(0)
	now := time.Now().UnixNano() / 1.0e6
	timePart := (id >> (ShardBits + RandomBits))
	shardPart := (id >> RandomBits) % int64(math.Pow(2, ShardBits))

	if now-timePart > 1000 {
		t.Errorf("Expect %v - %v <= 1000, got %v", now, timePart, now-timePart)
	}

	if shardPart != 0 {
		t.Errorf("Expect %v, got %v", 0, shardPart)
	}
}

func TestGetShard(t *testing.T) {
	n := Newsfeed{Kid: 1}
	if n.GetShard() != 1 {
		t.Errorf("Expect n.GetShard() to 1, got %v", n.GetShard())
	}

	n = Newsfeed{Kid: 2}
	if n.GetShard() != 0 {
		t.Errorf("Expect n.GetShard() to 0, got %v", n.GetShard())
	}
}

func TestValidate(t *testing.T) {
	// When no kid information
	n := Newsfeed{Summary: "{}"}
	if err := n.Validate(); err == nil {
		t.Errorf("Expect n.Validate() to raise")
	}

	// When summary is a bad json
	n = Newsfeed{Kid: 1, Summary: "[]]"}
	if err := n.Validate(); err == nil {
		t.Errorf("Expect n.Validate() to raise")
	}

	// When newsfeed_id is not present
	n = Newsfeed{Kid: 1, Summary: "{}"}
	if err := n.Validate(); err != nil {
		t.Errorf("Expect n.Validate() not to raise")
	}

	if n.NewsfeedID == 0 {
		t.Errorf("Expect n.NewsfeedID is not 0")
	}

	expectedSummary := fmt.Sprintf(`{"id":"%v"}`, n.NewsfeedID)
	if n.Summary != expectedSummary {
		t.Errorf("Expect n.Summary to %q, got %q", expectedSummary, n.Summary)
	}
}

func TestSave(t *testing.T) {
	defer shards[0].Delete(&Newsfeed{})
	defer shards[1].Delete(&Newsfeed{})

	// When Newsfeed should go to shard 0
	shards[0].Delete(&Newsfeed{})
	n := Newsfeed{Kid: 2, Summary: "{}"}
	n.Save()

	record := Newsfeed{}
	shards[0].First(&record, n.NewsfeedID)

	if record.NewsfeedID != n.NewsfeedID {
		t.Errorf(
			"Expect record.NewsfeedID to %v, got %v",
			n.NewsfeedID,
			record.NewsfeedID,
		)
	}

	// When Newsfeed should go to shard 1
	shards[1].Delete(&Newsfeed{})
	n = Newsfeed{Kid: 1, Summary: "{}"}
	n.Save()

	record = Newsfeed{}
	shards[1].First(&record, n.NewsfeedID)

	if record.NewsfeedID != n.NewsfeedID {
		t.Errorf(
			"Expect record.NewsfeedID to %v, got %v",
			n.NewsfeedID,
			record.NewsfeedID,
		)
	}

	// When any error occurred
	shards[0].Delete(&Newsfeed{})
	n = Newsfeed{Kid: 2, Summary: "{}"}
	n.Save()
	if err := n.Save(); err == nil {
		t.Errorf("Expect n.Save() to raise exception")
	}
}
