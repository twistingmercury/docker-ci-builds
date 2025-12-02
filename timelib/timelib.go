package timelib

import "time"

// Now returns the current time as a formatted string.
func Now() string {
	return time.Now().String()
}
