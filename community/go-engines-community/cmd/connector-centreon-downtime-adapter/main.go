package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"github.com/rs/zerolog"
)

const (
	EnvVarCanopsisUsername = "CPS_API_USERNAME"
	EnvVarCanopsisPassword = "CPS_API_PASSWORD"

	reconnectCount    = 3
	reconnectInterval = 8 * time.Millisecond

	prefetchCount = 10000
	prefetchSize  = 0

	defaultLocation = "Europe/Paris"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var flags Flags
	flags.ParseArgs()

	if flags.Version {
		canopsis.PrintVersionInfo()
		return
	}

	logger := log.NewLogger(flags.Debug)
	err := run(ctx, flags, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("adapter failed")
	}

	logger.Info().Msg("adapter stopped")
}

func run(ctx context.Context, flags Flags, logger zerolog.Logger) (resErr error) {
	config, err := LoadConfig(flags.ConfigPath)
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}

	config.Api.Username = os.Getenv(EnvVarCanopsisUsername)
	config.Api.Password = os.Getenv(EnvVarCanopsisPassword)

	amqpConn, err := amqp.NewConnection(logger, reconnectCount, reconnectInterval)
	if err != nil {
		return fmt.Errorf("cannot connect to rmq: %w", err)
	}

	defer func() {
		if err := amqpConn.Close(); err != nil && resErr == nil {
			resErr = err
		}
	}()

	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		return fmt.Errorf("cannot connect to rmq: %w", err)
	}

	err = amqpChannel.Qos(prefetchCount, prefetchSize, false)
	if err != nil {
		return fmt.Errorf("cannot modify rmq channel: %w", err)
	}

	dt, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return errors.New("unknown type of http.DefaultTransport")
	}

	transport := dt.Clone()
	if config.Api.InsecureSkipverify {
		logger.Warn().Msg("adapter accepts any certificate from canopsis api")
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}
	httpClient := &http.Client{Transport: transport}

	consumer := NewEventConsumer(amqpChannel, config, httpClient, logger)
	logger.Info().Msg("adapter started")

	return consumer.Start(ctx)
}
