package log

import (
	"fmt"
	"log"
	"os"
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
		s.Println("[DEBUG]", fmt.Sprint(v...))
	}
}

func (s *std) Debugf(format string, v ...interface{}) {
	if s.level >= LevelDebug {
		s.Println("[DEBUG]", fmt.Sprintf(format, v...))
	}
}

func (s *std) Info(v ...interface{}) {
	if s.level >= LevelInfo {
		s.Println("[INFO]", fmt.Sprint(v...))
	}
}

func (s *std) Infof(format string, v ...interface{}) {
	if s.level >= LevelInfo {
		s.Println("[INFO]", fmt.Sprintf(format, v...))
	}
}

func (s *std) Error(v ...interface{}) {
	if s.level >= LevelError {
		s.Println("[ERROR]", fmt.Sprint(v...))
	}
}

func (s *std) Errorf(format string, v ...interface{}) {
	if s.level >= LevelError {
		s.Println("[ERROR]", fmt.Sprintf(format, v...))
	}
}

// Level defines the level from which log should be displayed.
type Level int

const (
	LevelDebug Level = 2
	LevelInfo  Level = 1
	LevelError Level = 0
)

// NewStd returns a new logger for Harmony based on the standard logger.
func NewStd(l Level) Logger { return &std{Logger: log.New(os.Stdout, "", 0), level: l} }
