// Package logger provides a centralized structured logger built on log/slog.
//
// Call Init() once after environment variables are loaded. In DEV mode it
// emits human-readable text at Debug level; in production it emits JSON at
// Info level. All other packages should import this package and use its
// exported functions instead of the stdlib log package.
package logger

import (
	"log/slog"
	"os"
)

// Init configures the global logger based on the ENV environment variable.
// Must be called after godotenv.Load() so that ENV is already set.
func Init() {
	var handler slog.Handler
	if os.Getenv("ENV") == "DEV" {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	slog.SetDefault(slog.New(handler))
}

func Debug(msg string, args ...any) { slog.Debug(msg, args...) }
func Info(msg string, args ...any)  { slog.Info(msg, args...) }
func Warn(msg string, args ...any)  { slog.Warn(msg, args...) }
func Error(msg string, args ...any) { slog.Error(msg, args...) }

// Fatal logs at Error level then exits. Mirrors log.Fatal behaviour.
func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
