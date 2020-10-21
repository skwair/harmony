/*
Package discord contains every Discord resources used throughout
Harmony as well as some utility functions and methods to work with them.

There are two types of objects in this package: those suffixed by "Settings"
or "Parameters" and others. Settings and Parameters objects are sent by the
client to Discord. They all use the `optional` package to handle optional
fields properly. Convenience functions are provided to create Settings and
Parameters, see "New*Settings" and "New*Parameters" functions. Every other
object is sent from Discord to the client. Some of these objects fields can
be pointers. When this is the case, it means this field is nullable, be sure
to check whether it's set before accessing it.
*/
package discord
