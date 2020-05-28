/*
Package log defines an interface that can be implemented in order to provide a logger
for Harmony. A default implementation using Go's standard log package is also present.
*/
package log

import (
	"fmt"
	"io"
	"log"
)

// Logger is the interface to implement in order to provide a logger
// compatible with Harmony.
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Level() Level
}

type std struct {
	*log.Logger
	level Level
}

func (s *std) Debug(v ...interface{}) {
	if s.level >= LevelDebug {
		s.printWithPrefix("[DEBUG]", v...)
	}
}

func (s *std) Debugf(format string, v ...interface{}) {
	if s.level >= LevelDebug {
		s.printfWithPrefix("[DEBUG]", format, v...)
	}
}

func (s *std) Info(v ...interface{}) {
	if s.level >= LevelInfo {
		s.printWithPrefix("[INFO]", v...)
	}
}

func (s *std) Infof(format string, v ...interface{}) {
	if s.level >= LevelInfo {
		s.printfWithPrefix("[INFO]", format, v...)
	}
}

func (s *std) Warn(v ...interface{}) {
	if s.level >= LevelWarn {
		s.printWithPrefix("[WARN]", v...)
	}
}

func (s *std) Warnf(format string, v ...interface{}) {
	if s.level >= LevelWarn {
		s.printfWithPrefix("[WARN]", format, v...)
	}
}

func (s *std) Error(v ...interface{}) {
	if s.level >= LevelError {
		s.printWithPrefix("[ERROR]", v...)
	}
}

func (s *std) Errorf(format string, v ...interface{}) {
	if s.level >= LevelError {
		s.printfWithPrefix("[ERROR]", format, v...)
	}
}

func (s *std) Level() Level {
	return s.level
}

func (s *std) printWithPrefix(prefix string, v ...interface{}) {
	s.Println(prefix, fmt.Sprint(v...))
}

func (s *std) printfWithPrefix(prefix, format string, v ...interface{}) {
	s.Println(prefix, fmt.Sprintf(format, v...))
}

// Level defines the level from which log should be displayed.
type Level int

const (
	// LevelDebug traces everything Harmony does, it dumps every HTTP call
	// and logs every websocket message. Very useful for debugging or developing
	// new features.
	// Beware of debug level as it is very chatty and it will log sensitive
	// information such as bot tokens, voice connections secret keys, etc.
	LevelDebug Level = 3
	// LevelInfo is here to notify that something happened. There's generally nothing to do
	// about them, they are just here to inform about an event.
	// This is also the default log level of Harmony.
	LevelInfo Level = 2
	// LevelWarn is for important logs that indicates something wrong or unusual happened.
	// The application is still running but probably in a degraded way and might crash if
	// no action is taken to fix the issues reported.
	LevelWarn Level = 1
	// LevelError is for when something went really wrong, meaning the connection to the
	// Gateway is probably down and/or failed to reconnect. These are very often network
	// issues.
	LevelError Level = 0
)

// NewStd returns a new logger for Harmony based on the standard logger.
func NewStd(w io.Writer, l Level) Logger {
	return &std{Logger: log.New(w, "", log.LstdFlags), level: l}
}
