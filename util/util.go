package util

import (
	"math/rand"
	"time"
)

const (
	shardBits  = 12
	randomBits = 10
	maxRandom  = 1024 // 2**randomBits
)

// GenerateID generates a int64 id.
//
// It is composed of unix timestamp, shard, and random number. For the
// first 42 bits is unix timestamp, the next 12 bits is shard, and the
// remaining 10 bits is a random number.
func GenerateID(shard int) int64 {
	t := time.Now().UnixNano()
	milliSecond := t / 1.0e6

	ret := milliSecond << (shardBits + randomBits)
	ret += int64(shard << randomBits)
	ret += int64(rand.New(rand.NewSource(t)).Intn(maxRandom))

	return ret
}
