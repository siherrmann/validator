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
// It checks the format of the string and parses it accordingly.
// It supports various formats including local time, UTC time, with and without microseconds.
func ISO8601StringToTime(in string) (time.Time, error) {
	layout := ""
	// Iso8601 date string in local time (yyyy-MM-ddTHH:mm:ss.mmmuuu)
	match, _ := regexp.MatchString("^[-:.T0-9]{26}$", in)
	if match {
		layout = "2006-01-02T15:04:05.000000"
	}

	// Iso8601 date string in UTC time (yyyy-MM-ddTHH:mm:ss.mmmuuuZ)
	match, _ = regexp.MatchString("^[-:.T0-9]{26}Z$", in)
	if match {
		layout = "2006-01-02T15:04:05.000000Z"
	}

	// Iso8601 date string in local time without microseconds (yyyy-MM-ddTHH:mm:ss.mmm)
	match, _ = regexp.MatchString("^[-:.T0-9]{23}$", in)
	if match {
		layout = "2006-01-02T15:04:05.000"
	}

	// Iso8601 date string in UTC time without microseconds (yyyy-MM-ddTHH:mm:ss.mmmZ)
	match, _ = regexp.MatchString("^[-:.T0-9]{23}Z$", in)
	if match {
		layout = "2006-01-02T15:04:05.000Z"
	}

	date, err := time.Parse(layout, in)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing iso8601 string to time: %v", err)
	}
	return date, nil
}
