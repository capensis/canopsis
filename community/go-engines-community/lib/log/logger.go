// Package log defines the default loggers.
package log

import (
	"context"
	"flag"
	"io"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
)

const (
	Stderr   = "stderr"
	Journald = "journald"
)

type Options struct {
	Debug  bool
	Writer string
}

func BindCmdFlags(opts *Options) {
	flag.BoolVar(&opts.Debug, "d", false, "Debug mode")
	flag.StringVar(&opts.Writer, "cps.logger", "", "Logger output destination. Overrides the \"Canopsis.logger.Writer\" setting from the TOML config file.")
}

// NewLogger returns the default logger, that should be used by all the
// engines.
// The returned logger is thread-safe, and may be used in multiple goroutines.
func NewLogger(ctx context.Context, opts Options) zerolog.Logger {
	var (
		logger               zerolog.Logger
		loggerWriter, writer io.Writer
	)

	logLevel := zerolog.InfoLevel
	if opts.Debug {
		logLevel = zerolog.DebugLevel
	}

	// Default
	writer = os.Stdout
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	loggerWriter = consoleWriter

	cfg, err := loadLoggerConfig(ctx)
	if err != nil {
		logger = zerolog.New(loggerWriter).Level(logLevel).With().Timestamp().Caller().Logger()
		return logger
	}

	writerStr := opts.Writer
	if writerStr == "" {
		writerStr = cfg.Writer
	}

	if writerStr != "" {
		switch writerStr {
		case Stderr:
			writer = os.Stderr
		case Journald:
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

func Writers() []string {
	return []string{
		Stderr,
		Journald,
	}
}

func loadLoggerConfig(ctx context.Context) (*config.SectionLogger, error) {
	dbClient, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
	if err != nil {
		return nil, err
	}

	cfg, err := config.NewAdapter(dbClient).GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	err = dbClient.Disconnect(ctx)
	if err != nil {
		return nil, err
	}

	return &cfg.Logger, nil
}
