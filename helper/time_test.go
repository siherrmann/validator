package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestISO8601StringToTime(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "RFC3339 with timezone offset and nanoseconds",
			input:     "2026-01-19T14:37:45.673212+01:00",
			wantError: false,
		},
		{
			name:      "RFC3339 with Z suffix",
			input:     "2026-01-19T14:37:45Z",
			wantError: false,
		},
		{
			name:      "RFC3339 with Z suffix and nanoseconds",
			input:     "2026-01-19T14:37:45.123456789Z",
			wantError: false,
		},
		{
			name:      "Zero time with Z suffix",
			input:     "0001-01-01T00:00:00Z",
			wantError: false,
		},
		{
			name:      "RFC3339 with negative timezone offset",
			input:     "2026-01-19T14:37:45-05:00",
			wantError: false,
		},
		{
			name:      "DB format - local time with microseconds",
			input:     "2026-01-19T14:37:45.123456",
			wantError: false,
		},
		{
			name:      "DB format - UTC time with microseconds",
			input:     "2026-01-19T14:37:45.123456Z",
			wantError: false,
		},
		{
			name:      "DB format - local time with milliseconds",
			input:     "2026-01-19T14:37:45.123",
			wantError: false,
		},
		{
			name:      "DB format - UTC time with milliseconds",
			input:     "2026-01-19T14:37:45.123Z",
			wantError: false,
		},
		{
			name:      "Invalid format",
			input:     "not-a-date",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ISO8601StringToTime(tt.input)
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.False(t, result.IsZero() && tt.input != "0001-01-01T00:00:00Z")
		})
	}
}

func TestUnixStringToTime(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "Valid Unix timestamp",
			input:     "1737292665",
			wantError: false,
		},
		{
			name:      "Zero timestamp",
			input:     "0",
			wantError: false,
		},
		{
			name:      "Invalid format with letters",
			input:     "12345abc",
			wantError: true,
		},
		{
			name:      "Empty string",
			input:     "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnixStringToTime(tt.input)
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.False(t, result.IsZero() && tt.input != "0")
		})
	}
}
