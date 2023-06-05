// Package log defines the default loggers.
package log

import (
	"context"
	"io"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
)

// NewLogger returns the default logger, that should be used by all the
// engines.
// The returned logger is thread-safe, and may be used in multiple goroutines.
func NewLogger(debug bool) zerolog.Logger {
	var (
		logger               zerolog.Logger
		loggerWriter, writer io.Writer
	)

	logLevel := zerolog.InfoLevel
	if debug {
		logLevel = zerolog.DebugLevel
	}

	// Default
	writer = os.Stdout
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	loggerWriter = consoleWriter

	cfg, err := loadLoggerConfig()
	if err != nil {
		logger = zerolog.New(loggerWriter).Level(logLevel).With().Timestamp().Caller().Logger()
		return logger
	}

	if cfg.Writer != "" {
		switch cfg.Writer {
		case "stderr":
			writer = os.Stderr
		case "journald":
			writer = journald.NewJournalDWriter()
		}
		loggerWriter = writer
	}

	if cfg.ConsoleWriter.Enabled {
		consoleWriter.Out = writer
		consoleWriter.NoColor = cfg.ConsoleWriter.NoColor
		if cfg.ConsoleWriter.TimeFormat != "" {
			consoleWriter.TimeFormat = cfg.ConsoleWriter.TimeFormat
		}
		if len(cfg.ConsoleWriter.PartsOrder) > 0 {
			consoleWriter.PartsOrder = cfg.ConsoleWriter.PartsOrder
		}
		loggerWriter = consoleWriter
	}

	// The writer should be thread-safe so that the logger can be used in
	// multiple goroutines. This writer is thread-safe, since it writes to
	// os.Stdout which is an os.File.
	// It may be necessary to wrap other writers with zerolog.SyncWriter.
	// For more details, read :
	// https://godoc.org/github.com/rs/zerolog#SyncWriter
	logger = zerolog.New(loggerWriter).Level(logLevel).With().Timestamp().Caller().Logger()
	return logger
}

// NewTestLogger returns the default test logger, that should be used by all
// the engines.
func NewTestLogger() zerolog.Logger {
	return NewLogger(false)
}

func loadLoggerConfig() (*config.SectionLogger, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbClient, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
	if err != nil {
		return nil, err
	}

	cfg, err := config.NewAdapter(dbClient).GetConfig(ctx)
	dbErr := dbClient.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	if dbErr != nil {
		return nil, dbErr
	}

	return &cfg.Logger, nil
}
