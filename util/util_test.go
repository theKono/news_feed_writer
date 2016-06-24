// +build unit

package util

import (
	"math/rand"
	"testing"
	"time"
)

func TestGenerateID(t *testing.T) {
	shard := rand.Int() % 2
	id := GenerateID(shard)
	now := time.Now()
	timePart := id >> 22
	s := (id >> 10) % 4096

	if d := now.UnixNano()/1.0e6 - timePart; d >= 1000 {
		t.Fatalf("Expect time discrepency is less than 1 second, but got `%v`", d)
	}

	if int(s) != shard {
		t.Fatalf("Expect shard to be `%v`, but got `%v`", shard, s)
	}
}
