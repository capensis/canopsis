package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils/ready"
)

func main() {
	var flagRetries int
	var flagTimeout time.Duration
	var withPostgres, withTechPostgres bool
	flag.IntVar(&flagRetries, "retries", 10, "number of retries per check. if 0, infinite number of retries")
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

	err := ready.CheckAll(ctx, time.Second, flagRetries, withPostgres, withTechPostgres, logger)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Fatal().Err(err).Msg("failed to open one of required sessions before shutdown")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Fatal().Err(err).Msgf("failed to open one of required sessions after %s", flagTimeout.String())
		}
		logger.Fatal().Err(err).Msg("failed to open one of required sessions")
	}

	logger.Info().Msg("ready!")
}
