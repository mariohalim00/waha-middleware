package util

import (
	"math/rand"
	"time"
)

func GenerateRandomDuration(minDuration, maxDuration int) time.Duration {
	num := rand.Intn(maxDuration-minDuration+1) + minDuration
	duration := time.Duration(num) * time.Second
	return duration
}
