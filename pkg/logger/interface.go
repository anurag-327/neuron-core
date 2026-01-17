package logger

import (
	"time"
)

// Logger is the main interface for pushing logs.
// Services should depend on this interface, not concrete structs.
type Logger interface {
	// Log sends a generic log entry
	Log(entry LogEntry) error

	// Helper functions for common log levels
	Info(Timestamp time.Time, msg string, meta map[string]interface{}) error
	Error(Timestamp time.Time, msg string, meta map[string]interface{}) error
	Warn(Timestamp time.Time, msg string, meta map[string]interface{}) error
	Debug(Timestamp time.Time, msg string, meta map[string]interface{}) error

	// Close cleans up underlying connections
	Close() error
}
