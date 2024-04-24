package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GetRandomTimeBetween returns a random time value between two durations
func GetRandomTimeBetween(minDuration, maxDuration time.Duration) time.Time {
	// Convert durations to milliseconds
	minMilliseconds := minDuration.Milliseconds()
	maxMilliseconds := maxDuration.Milliseconds()

	// Generate a random number of milliseconds between min and max
	randomMilliseconds := rand.Int63n(maxMilliseconds-minMilliseconds+1) + minMilliseconds

	// Create a time.Duration from the random milliseconds
	randomDuration := time.Duration(randomMilliseconds) * time.Millisecond

	// Add the random duration to the minimum time
	return time.Now().Add(randomDuration)
}

func ParseMilliseconds(timeStr string) (int, error) {
	cleanTimeStr := strings.TrimSuffix(timeStr, "ms")
	return strconv.Atoi(cleanTimeStr)
}

// function that parses a string to a time.Duration
func ParseDuration(durationStr string) (time.Duration, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, err
	}
	return duration, nil
}
