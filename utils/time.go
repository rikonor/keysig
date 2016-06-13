package utils

import (
	"fmt"
	"time"
)

// DurationToMSString converts a time.Duration to ms string
// Example result is 37.12312
func DurationToMSString(d time.Duration) string {
	return fmt.Sprint(float64(d.Nanoseconds()) / float64((1000 * 1000)))
}
