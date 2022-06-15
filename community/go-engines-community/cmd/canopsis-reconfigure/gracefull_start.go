package main

import (
	"fmt"
	"log"
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

func getValue(paramName string, defaultValue int) int {
	value := os.Getenv(paramName)

	if value == "" {
		return defaultValue
	}

	res, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("graceful start: getValue: %s=%v: %v", paramName, value, err)
	}
	return res
}

// GracefullStart will try to initialize every
func GracefullStart(logger zerolog.Logger) error {
	maxRetry := getValue(EnvCpsStartMaxRetry, DefaultMaxRetry)
	retryDelay, err := time.ParseDuration(fmt.Sprintf("%ds", getValue(EnvCpsStartRetryDelay, DefaultRetryDelay)))

	if err != nil {
		return fmt.Errorf("retry delay: %v", err)
	}

	waitFirstAttempt, err := time.ParseDuration(fmt.Sprintf("%ds", getValue(EnvCpsStartWaitFirst, DefaultWaitFirstAttempt)))
	if err != nil {
		return fmt.Errorf("wait first attempt: %v", err)
	}

	logger.Info().Msgf("waiting %s before first attempt", waitFirstAttempt.String())
	time.Sleep(waitFirstAttempt)

	logger.Info().Msg("checking")
	ready.Abort(ready.Check(ready.CheckRedis, "redis", retryDelay, maxRetry))
	ready.Abort(ready.Check(ready.CheckMongo, "mongo", retryDelay, maxRetry))
	ready.Abort(ready.Check(ready.CheckAMQP, "amqp", retryDelay, maxRetry))
	//ready.Abort(ready.Check(ready.CheckInflux, "influx", retryDelay, maxRetry))

	return nil
}
