package optional

import "encoding/json"

// Int represents an optional int.
type Int struct {
	i   int
	nil bool
}

// MarshalJSON implements the json.Marshaler interface.
func (i *Int) MarshalJSON() ([]byte, error) {
	if i.nil {
		return []byte(`null`), nil
	}

	return json.Marshal(i.i)
}

// NewInt returns a new optional int set to i.
func NewInt(i int) *Int {
	return &Int{
		i: i,
	}
}

// NewNilInt returns a new optional int set to nil.
func NewNilInt() *Int {
	return &Int{
		nil: true,
	}
}
