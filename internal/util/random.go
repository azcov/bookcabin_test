package util

import (
	"math/rand"
	"time"
)

func RandomDuration(minMs, maxMs int) time.Duration {
	diff := maxMs - minMs
	return time.Duration(minMs+rand.Intn(diff)) * time.Millisecond
}

func RandomFailure(failRate float64) bool {
	return rand.Float64() < failRate
}
