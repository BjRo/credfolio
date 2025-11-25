package logger

import (
	"io"
	"log"
	"os"
)

// Level represents log level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger provides structured logging
type Logger struct {
	level  Level
	output io.Writer
	debug  *log.Logger
	info   *log.Logger
	warn   *log.Logger
	err    *log.Logger
}

// New creates a new logger with the given level
func New(level Level) *Logger {
	return NewWithOutput(level, os.Stdout)
}

// NewWithOutput creates a new logger with custom output
func NewWithOutput(level Level, output io.Writer) *Logger {
	flags := log.LstdFlags | log.Lmsgprefix

	return &Logger{
		level:  level,
		output: output,
		debug:  log.New(output, "[DEBUG] ", flags),
		info:   log.New(output, "[INFO] ", flags),
		warn:   log.New(output, "[WARN] ", flags),
		err:    log.New(output, "[ERROR] ", flags),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.debug.Printf(msg, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.info.Printf(msg, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.warn.Printf(msg, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= LevelError {
		l.err.Printf(msg, args...)
	}
}

// SetLevel changes the logging level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// DefaultLogger is the default logger instance
var DefaultLogger = New(LevelInfo)

// Debug logs a debug message using the default logger
func Debug(msg string, args ...interface{}) {
	DefaultLogger.Debug(msg, args...)
}

// Info logs an info message using the default logger
func Info(msg string, args ...interface{}) {
	DefaultLogger.Info(msg, args...)
}

// Warn logs a warning message using the default logger
func Warn(msg string, args ...interface{}) {
	DefaultLogger.Warn(msg, args...)
}

// Error logs an error message using the default logger
func Error(msg string, args ...interface{}) {
	DefaultLogger.Error(msg, args...)
}
