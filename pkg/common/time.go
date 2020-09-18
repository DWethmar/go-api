package common

import "time"

var defaultTimePrecision = time.Microsecond

// Now gets the current time with default time precision
func Now() time.Time {
	return time.Now().Truncate(defaultTimePrecision)
}

// DefaultTimePrecision truncates the time to use microseconds.
func DefaultTimePrecision(t *time.Time) time.Time {
	return t.Truncate(defaultTimePrecision)
}
