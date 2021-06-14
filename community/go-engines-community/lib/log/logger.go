// Package log defines the default loggers.
package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// NewLogger returns the default logger, that should be used by all the
// engines.
// The returned logger is thread-safe, and may be used in multiple goroutines.
func NewLogger(debug bool) zerolog.Logger {
	logLevel := zerolog.InfoLevel
	if debug {
		logLevel = zerolog.DebugLevel
	}

	// The writer should be thread-safe so that the logger can be used in
	// multiple goroutines. This writer is thread-safe, since it writes to
	// os.Stdout which is an os.File.
	// It may be necessary to wrap other writers with zerolog.SyncWriter.
	// For more details, read :
	// https://godoc.org/github.com/rs/zerolog#SyncWriter
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	return zerolog.New(writer).Level(logLevel).With().Timestamp().Caller().Logger()
}

// NewTestLogger returns the default test logger, that should be used by all
// the engines.
func NewTestLogger() zerolog.Logger {
	return NewLogger(false)
}
