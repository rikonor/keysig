package utils

import (
	"fmt"
	"time"
)

// DurationToMSFloat64 converts a time.Duration to ms float64
// Example result is 37.12312
func DurationToMSFloat64(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / float64((1000 * 1000))
}

// DurationToMSString converts a time.Duration to ms string
// Example result is "37.12312"
func DurationToMSString(d time.Duration) string {
	return fmt.Sprint(DurationToMSFloat64(d))
}
