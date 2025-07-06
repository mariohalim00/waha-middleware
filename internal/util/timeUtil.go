package util

import (
	"math/rand"
	"time"
)

func GenerateRandomDuration(maxDuration int) time.Duration {
	num := rand.Intn(maxDuration) + 1
	duration := time.Duration(num) * time.Second
	return duration
}
