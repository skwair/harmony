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
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
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

func (s *std) printWithPrefix(prefix string, v ...interface{}) {
	s.Println(prefix, fmt.Sprint(v...))
}

func (s *std) printfWithPrefix(prefix, format string, v ...interface{}) {
	s.Println(prefix, fmt.Sprintf(format, v...))
}

// Level defines the level from which log should be displayed.
type Level int

const (
	LevelDebug Level = 2
	LevelInfo  Level = 1
	LevelError Level = 0
)

// NewStd returns a new logger for Harmony based on the standard logger.
func NewStd(w io.Writer, l Level) Logger { return &std{Logger: log.New(w, "", log.LstdFlags), level: l} }
