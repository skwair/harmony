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
