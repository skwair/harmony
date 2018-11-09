/*
Package optional defines optional versions of primitive types that can be nil. These are
useful when performing partial updates of some objects because it is impossible to tell
the difference between a value that is explicitly set to its zero value and a value
that is not set at all when marshaling a struct with native primitive types.
*/
package optional
