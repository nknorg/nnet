package util

import (
	crand "crypto/rand"
	mrand "math/rand"
	"time"
)

// RandBytes returns a random []byte with length l
func RandBytes(len int) ([]byte, error) {
	b := make([]byte, len)
	_, err := crand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// RandDuration generates a random duration
func RandDuration(average time.Duration, delta float64) time.Duration {
	// uniform random number between (1 - delta) * average to (1 + delta) * average
	return time.Duration((1 - delta + 2*mrand.Float64()*delta) * float64(average))
}
