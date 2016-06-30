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

func TestMeasureExecTime(t *testing.T) {
	targetFn := func() time.Duration {
		now := time.Now()
		time.Sleep(10 * time.Millisecond)
		return MeasureExecTime(now, "targetFn")
	}

	execTime := targetFn()
	if execTime-10*time.Millisecond >= time.Millisecond {
		t.Fatalf("Expect execTime to < 0.001 sec, but got `%v`", execTime-10*time.Millisecond)
	}
}
