/*
Package discord contains every Discord resources used throughout
Harmony as well as some utility functions and methods to work with them.

Objects suffixed by "Settings" or "Parameters" are settings and parameters
objects are sent by the client to Discord. They all use the `optional` package
to handle optional fields properly. Convenience functions are provided to create
settings and parameters, see "New*Settings" and "New*Parameters" functions.
Most of other objects are sent from Discord to the client. Some of these objects
contain fields that are pointers. When this is the case, it means this field
is nullable, be sure to check whether it's set before accessing it.
*/
package discord
