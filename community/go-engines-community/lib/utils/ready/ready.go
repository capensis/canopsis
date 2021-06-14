package ready

import (
	"fmt"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/influx"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
)

var logger = log.NewLogger(false)

// Timeout is to be ran by a go routine, so when the duration goes out,
// timeout will call log.Fatal
func Timeout(t time.Duration) {
	time.Sleep(t)
	logger.Fatal().Msgf("timeout expired after %s", t.String())
}

// Check calls checker function for a job to do.
// name is a dispaly string
// retryDelay is the time to sleep before doing another attempt
// retries is the max number of tries to attempt
func Check(checker func() error, name string, retryDelay time.Duration, retries int) error {
	infinite := false

	if retries == 0 {
		infinite = true
	}

	for retry := 0; retry < retries || infinite; retry++ {
		err := checker()
		if err != nil {
			logger.Info().Err(err).Msgf("%v: %v/%v", name, retry+1, retries)
			time.Sleep(retryDelay)
		} else {
			return nil
		}
	}

	return fmt.Errorf("check %v failed after %v tries", name, retries)
}

// Abort log.Fatalf() when err != nil
func Abort(err error) {
	if err != nil {
		logger.Fatal().Err(err).Msg("hard failure")
	}
}

// CheckRedis ...
func CheckRedis() error {
	_, err := redis.NewSession(0, logger, 0, 0)
	return err
}

// CheckMongo connects to mongo with mongo.Timeout
func CheckMongo() error {
	_, err := mongo.NewSession(mongo.Timeout)
	return err
}

// CheckAMQP ...
func CheckAMQP() error {
	_, err := amqp.NewConnection(log.NewLogger(false), 0, 0)
	return err
}

// CheckInflux ...
func CheckInflux() error {
	_, err := influx.NewSession()
	return err
}
