package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// CustomTime handles multiple time formats from Certimate API including null
type CustomTime struct {
	time.Time
	Valid bool
}

// UnmarshalJSON handles multiple time formats
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")

	if str == "" || str == "null" {
		ct.Valid = false
		return nil
	}

	// List of formats to try
	formats := []string{
		"2006-01-02 15:04:05.999Z",   // Certimate format with milliseconds
		"2006-01-02 15:04:05Z",       // Certimate format without milliseconds
		"2006-01-02 15:04:05.999",    // Without Z
		"2006-01-02 15:04:05",        // Simple datetime
		time.RFC3339,                  // Standard RFC3339
		time.RFC3339Nano,              // RFC3339 with nanoseconds
	}

	var err error
	for _, format := range formats {
		ct.Time, err = time.Parse(format, str)
		if err == nil {
			ct.Valid = true
			return nil
		}
	}

	return fmt.Errorf("cannot parse time: %s", str)
}

// MarshalJSON returns the time in RFC3339 format
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if !ct.Valid || ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(time.RFC3339))), nil
}

// IsZero returns true if the time is not set
func (ct CustomTime) IsZero() bool {
	return !ct.Valid || ct.Time.IsZero()
}

// NewCustomTime creates a CustomTime from time.Time
func NewCustomTime(t time.Time) CustomTime {
	return CustomTime{Time: t, Valid: true}
}

// StringArray handles both string (semicolon-separated) and array formats
type StringArray []string

// UnmarshalJSON handles both string and array formats
func (sa *StringArray) UnmarshalJSON(data []byte) error {
	// Try as array first
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*sa = arr
		return nil
	}

	// Try as string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Split by semicolon
	if str == "" {
		*sa = []string{}
	} else {
		*sa = strings.Split(str, ";")
	}

	return nil
}

// MarshalJSON returns the array as JSON
func (sa StringArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(sa))
}
