package audit

// StringValues holds a pair of string values.
type StringValues struct {
	Old, New string
}

// IntValues holds a pair of integer values.
type IntValues struct {
	Old, New int
}

// BoolValues holds a pair of boolean values.
type BoolValues struct {
	Old, New bool
}
