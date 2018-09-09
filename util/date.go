package util

import (
	"time"
)

// GetNow get the time string of time RFC3339
func GetNow() string {
	return time.Now().Format(time.RFC3339)
}

// GetUTCNow get the utc time string of time RFC3339
func GetUTCNow() string {
	return time.Now().UTC().Format(time.RFC3339)
}
