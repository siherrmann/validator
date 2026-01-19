package helper

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// UnixStringToTime converts a Unix timestamp string to a time.Time object.
// It checks if the string is a valid Unix timestamp and parses it.
func UnixStringToTime(in string) (time.Time, error) {
	// Unix seconds date string
	match, _ := regexp.MatchString("^[0-9]{1,}$", in)
	if !match {
		return time.Time{}, fmt.Errorf("invalid unix time: %v", in)
	}

	seconds, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing unix string to time: %v", err)
	}
	return time.Unix(seconds, 0).UTC(), nil
}

// ISO8601StringToTime converts an ISO8601 date string to a time.Time object.
// It supports RFC3339 format timestamps including timezone offsets,
// as well as database formats without timezone information.
func ISO8601StringToTime(in string) (time.Time, error) {
	// RFC3339Nano handles all formats with timezone (Z or +/-HH:MM) including nanoseconds
	// Examples: 2026-01-19T14:37:45.673212+01:00, 2026-01-19T14:37:45.123Z, 2026-01-19T14:37:45Z
	if date, err := time.Parse(time.RFC3339Nano, in); err == nil {
		return date, nil
	}
	if date, err := time.Parse(time.RFC3339, in); err == nil {
		return date, nil
	}

	// Database formats without timezone - local time only
	layout := ""
	switch len(in) {
	case 26: // yyyy-MM-ddTHH:mm:ss.mmmuuu (microseconds)
		layout = "2006-01-02T15:04:05.000000"
	case 23: // yyyy-MM-ddTHH:mm:ss.mmm (milliseconds)
		layout = "2006-01-02T15:04:05.000"
	default:
		return time.Time{}, fmt.Errorf("error parsing iso8601 string to time: unsupported format")
	}

	date, err := time.Parse(layout, in)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing iso8601 string to time: %v", err)
	}
	return date, nil
}
