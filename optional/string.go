package optional

import "encoding/json"

// String represents an optional string.
type String struct {
	s   string
	nil bool
}

// MarshalJSON implements the json.Marshaler interface.
func (s *String) MarshalJSON() ([]byte, error) {
	if s.nil {
		return []byte(`null`), nil
	}

	return json.Marshal(s.s)
}

// NewString returns a new optional string set to s.
func NewString(s string) *String {
	return &String{
		s: s,
	}
}

// NewNilString returns a new optional string set to nil.
func NewNilString() *String {
	return &String{
		nil: true,
	}
}

// StringSlice represents an optional string slice.
type StringSlice struct {
	ss  []string
	nil bool
}

// MarshalJSON implements the json.Marshaler interface.
func (s *StringSlice) MarshalJSON() ([]byte, error) {
	if s.nil {
		return []byte(`"null"`), nil
	}

	return json.Marshal(s.ss)
}

// NewStringSlice returns a new optional string slice set to ss.
// If ss is nil, it is equivalent to NewNilStringSlice.
func NewStringSlice(ss []string) *StringSlice {
	if ss == nil {
		return NewNilStringSlice()
	}

	return &StringSlice{
		ss: ss,
	}
}

// NewNilStringSlice returns a new optional string slice set to nil.
func NewNilStringSlice() *StringSlice {
	return &StringSlice{
		nil: true,
	}
}
