package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Date time.Time

const DateFormat = "2006-01-02"

// UnmarshalJSON parses YYYY-MM-DD or RFC3339 format
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	// Try parsing YYYY-MM-DD
	t, err := time.Parse(DateFormat, s)
	if err == nil {
		*d = Date(t)
		return nil
	}
	// Fallback to RFC3339
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		*d = Date(t)
		return nil
	}
	return fmt.Errorf("invalid date format: %s (expected YYYY-MM-DD)", s)
}

// MarshalJSON formats as YYYY-MM-DD
func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(t.Format(DateFormat))
}

// Value implements driver.Valuer for database compatibility
func (d Date) Value() (driver.Value, error) {
	t := time.Time(d)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

// Scan implements sql.Scanner for database compatibility
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*d = Date(v)
		return nil
	case string:
		t, err := time.Parse(DateFormat, v)
		if err != nil {
			t, err = time.Parse(time.RFC3339, v)
		}
		if err != nil {
			return err
		}
		*d = Date(t)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Date", value)
	}
}

// Time returns the underlying time.Time
func (d Date) Time() time.Time {
	return time.Time(d)
}
