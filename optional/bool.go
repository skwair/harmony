package optional

import "encoding/json"

// Bool represents an optional boolean value.
type Bool struct {
	b   bool
	nil bool
}

// MarshalJSON implements the json.Marshaler interface.
func (b *Bool) MarshalJSON() ([]byte, error) {
	if b.nil {
		return []byte(`null`), nil
	}

	return json.Marshal(b.b)
}

// NewBool returns a new optional bool set to b.
func NewBool(b bool) *Bool {
	return &Bool{
		b: b,
	}
}

// NewNilBool returns a new optional bool set to nil.
func NewNilBool() *Bool {
	return &Bool{
		nil: true,
	}
}
