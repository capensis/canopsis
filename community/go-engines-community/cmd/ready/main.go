package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils/ready"
	"github.com/rs/zerolog"
)

func main() {
	var retries int
	var flagTimeout time.Duration
	var withPostgres, withTechPostgres bool
	flag.IntVar(&retries, "retries", 10, "number of retries per check. if 0, infinite number of retries")
	flag.DurationVar(&flagTimeout, "timeout", time.Second*60, "timeout after given duration. never timeout if 0s")
	flag.BoolVar(&withPostgres, "withPostgres", false, "check postgres")
	flag.BoolVar(&withTechPostgres, "withTechPostgres", false, "check tech postgres")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	var logger = log.NewLogger(false)

	if flagTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, flagTimeout)
		defer cancel()
	}

	retryDelay := time.Second
	err := ready.CheckAll(ctx, retryDelay, retries, withPostgres, withTechPostgres, logger)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Fatal().Err(err).Msg("failed to open one of required sessions before shutdown")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Fatal().Err(err).Msgf("failed to open one of required sessions after %s", flagTimeout.String())
		}
		logger.Fatal().Err(err).Msg("failed to open one of required sessions")
	}

	buildInfo := canopsis.GetBuildInfo()
	err = ready.Check(ctx, func(ctx context.Context, logger zerolog.Logger) error {
		dbClient, err := mongo.NewClient(ctx, 0, 0, logger)
		if err != nil {
			return err
		}

		defer func() {
			_ = dbClient.Disconnect(ctx)
		}()
		versionConfAdapter := config.NewVersionAdapter(dbClient)
		conf, err := versionConfAdapter.GetConfig(ctx)
		if err != nil {
			return err
		}

		if conf.Version != buildInfo.Version {
			return fmt.Errorf("expected version %q but got %q", buildInfo.Version, conf.Version)
		}

		return nil
	}, "reconfigure", retryDelay, retries, logger)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Fatal().Err(err).Msg("failed to check reconfigure before shutdown")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Fatal().Err(err).Msgf("failed to check reconfigure after %s", flagTimeout.String())
		}
		logger.Fatal().Err(err).Msg("failed to check reconfigure")
	}

	logger.Info().Msg("ready!")
}
