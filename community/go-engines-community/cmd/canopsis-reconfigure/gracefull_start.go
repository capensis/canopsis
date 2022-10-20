package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils/ready"
	"github.com/rs/zerolog"
)

// Retry parameters
const (
	EnvCpsStartMaxRetry     = "CPS_MAX_RETRY"
	EnvCpsStartRetryDelay   = "CPS_MAX_DELAY"
	EnvCpsStartWaitFirst    = "CPS_WAIT_FIRST_ATTEMPT"
	DefaultMaxRetry         = 10
	DefaultRetryDelay       = 30
	DefaultWaitFirstAttempt = 10
)

func GracefulStart(ctx context.Context, withPostgres, withTechPostgres bool, logger zerolog.Logger) error {
	maxRetry := getEnv(EnvCpsStartMaxRetry, DefaultMaxRetry, logger)
	retryDelay := time.Second * time.Duration(getEnv(EnvCpsStartRetryDelay, DefaultRetryDelay, logger))
	waitFirstAttempt := time.Second * time.Duration(getEnv(EnvCpsStartWaitFirst, DefaultWaitFirstAttempt, logger))

	logger.Info().Msgf("waiting %s before first attempt", waitFirstAttempt.String())
	t := time.NewTimer(waitFirstAttempt)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
	}

	return ready.CheckAll(ctx, retryDelay, maxRetry, withPostgres, withTechPostgres, logger)
}

func getEnv(paramName string, defaultValue int, logger zerolog.Logger) int {
	value := os.Getenv(paramName)
	if value == "" {
		return defaultValue
	}

	res, err := strconv.Atoi(value)
	if err != nil {
		logger.Warn().Err(err).Msgf("%q must be an integer, use default value instead", paramName)
		return defaultValue
	}

	if res <= 0 {
		logger.Warn().Err(err).Msgf("%q must be greater than zero, use default value instead", paramName)
		return defaultValue
	}

	return res
}
