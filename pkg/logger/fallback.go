package logger

import (
	"fmt"
	"time"
)

// fallbackLogger is a minimal logger used before SetGlobalLogger is called.
// It prevents crashes and prints to stdout.
type fallbackLogger struct{}

func NewFallbackLogger() Logger {
	return &fallbackLogger{}
}

func (f *fallbackLogger) Log(entry LogEntry) error {
	fmt.Printf("[FALLBACK] [%s] %s\n", entry.Level, entry.Message)
	return nil
}
func (f *fallbackLogger) Info(ts time.Time, msg string, meta map[string]interface{}) error {
	return f.Log(LogEntry{Level: LevelInfo, Message: msg})
}
func (f *fallbackLogger) Error(ts time.Time, msg string, meta map[string]interface{}) error {
	return f.Log(LogEntry{Level: LevelError, Message: msg})
}
func (f *fallbackLogger) Warn(ts time.Time, msg string, meta map[string]interface{}) error {
	return f.Log(LogEntry{Level: LevelWarn, Message: msg})
}
func (f *fallbackLogger) Debug(ts time.Time, msg string, meta map[string]interface{}) error {
	return f.Log(LogEntry{Level: LevelDebug, Message: msg})
}
func (f *fallbackLogger) Close() error { return nil }
