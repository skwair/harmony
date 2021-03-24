package discord

import (
	"strings"
	"time"
)

// Time wraps the standard time.Time structure and handles
// zero-value JSON marshaling properly.
type Time struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}

	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.ReplaceAll(str, `"`, "")

	tt, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	t.Time = tt
	return nil
}

// Std returns the standard time.Time value this Time represents.
func (t Time) Std() time.Time {
	return t.Time
}

// TimeFromStd converts a time.Time to a Time.
func TimeFromStd(t time.Time) Time {
	return Time{t}
}
