package log

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// A little JSONL file logger.
// IRL this should probably use a *real* library for this sort of thing.

var LOGGER *Logger = nil

// Log levels.
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// A log entry
type logEntry struct {
	Timestamp string   `json:"timestamp"`
	Level     LogLevel `json:"level"`
	Message   string   `json:"message"`
}

// A logger
type Logger struct {
	// The file to which it logs.
	file *os.File

	// Whether or not it does anything at all in practice.
	isLogging bool
}

// Initialize the Logger instance.
func Init(debug bool) error {
	var err error
	LOGGER, err = New(debug)
	return err
}

// Creates a new Logger instance.
func New(debug bool) (*Logger, error) {
	if LOGGER != nil {
		return LOGGER, nil
	}
	if debug {
		file, err := os.OpenFile("jkm.logs.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		return &Logger{file: file, isLogging: debug}, nil
	}
	return &Logger{file: nil, isLogging: debug}, nil
}

// Release the Logger's resources.
func (logger *Logger) Close() error {
	if logger.file == nil {
		return nil
	}
	return logger.file.Close()
}

// Log a message at the specified level with source file and line information
func (l *Logger) log(level LogLevel, msg string) error {
	if !l.isLogging || l.file == nil {
		return nil
	}

	entry := logEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = l.file.Write(append(jsonData, '\n'))
	return err
}

// Log the given message at info level
func (l *Logger) Log(msg string) error {
	return l.log(LevelInfo, msg)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) error {
	return l.log(LevelDebug, msg)
}

// Info logs an info message
func (l *Logger) Info(msg string) error {
	return l.log(LevelInfo, msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) error {
	return l.log(LevelWarn, msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) error {
	return l.log(LevelError, msg)
}

// Info through the singleton.
func Info(msg string) error {
	return LOGGER.Info(msg)
}

// Debug through the singleton
func Debug(msg string) error {
	return LOGGER.Debug(msg)
}

// Warn through the singleton
func Warn(msg string) error {
	return LOGGER.Warn(msg)
}

// Error through the singleton
func Error(msg string) error {
	return LOGGER.Error(msg)
}

// Convenience method to log an error
func LogError(err error) error {
	if err == nil {
		return nil
	}
	return Error(err.Error())
}

// Close the singleton.
func Close() error {
	return LOGGER.Close()
}

// Format a string and log at info level.
func Infof(format string, args ...any) error {
	return Info(fmt.Sprintf(format, args...))
}

// Format a string and log at debug level.
func Debugf(format string, args ...any) error {
	return Debug(fmt.Sprintf(format, args...))
}

// Format a string and log at warn level.
func Warnf(format string, args ...any) error {
	return Warn(fmt.Sprintf(format, args...))
}

// Format a string and log at error level.
func Errorf(format string, args ...any) error {
	return Error(fmt.Sprintf(format, args...))
}
