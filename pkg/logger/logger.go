package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel string

const (
	// INFO represents informational messages
	INFO LogLevel = "INFO"
	// WARN represents warning messages
	WARN LogLevel = "WARN"
	// ERROR represents error messages
	ERROR LogLevel = "ERROR"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// Logger provides structured logging functionality
type Logger struct {
	output *log.Logger
}

// New creates a new Logger instance
func New() *Logger {
	return &Logger{
		output: log.New(os.Stdout, "", 0),
	}
}

// Info logs an informational message
func (l *Logger) Info(message string, data map[string]interface{}) {
	l.log(INFO, message, data)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, data map[string]interface{}) {
	l.log(WARN, message, data)
}

// Error logs an error message
func (l *Logger) Error(message string, data map[string]interface{}) {
	l.log(ERROR, message, data)
}

// log outputs a structured log entry
func (l *Logger) log(level LogLevel, message string, data map[string]interface{}) {
	entry := LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple logging if JSON marshaling fails
		l.output.Printf("%s [%s] %s: %v", entry.Timestamp, entry.Level, message, data)
		return
	}

	l.output.Println(string(jsonData))
}