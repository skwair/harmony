package discord

import (
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
	var ts time.Time
	if err := ts.UnmarshalJSON(data); err != nil {
		return err
	}

	t.Time = ts
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
