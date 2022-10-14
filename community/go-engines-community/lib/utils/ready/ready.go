package ready

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/rs/zerolog"
)

func CheckAll(
	ctx context.Context,
	retryDelay time.Duration,
	retries int,
	withPostgres, withTechPostgres bool,
	logger zerolog.Logger,
) error {
	logger.Info().Msg("checking")
	err := Check(ctx, CheckRedis, "redis", retryDelay, retries, logger)
	if err != nil {
		return err
	}
	err = Check(ctx, CheckMongo, "mongo", retryDelay, retries, logger)
	if err != nil {
		return err
	}
	err = Check(ctx, CheckAMQP, "amqp", retryDelay, retries, logger)
	if err != nil {
		return err
	}

	if withPostgres {
		err = Check(ctx, CheckPostgres, "postgres", retryDelay, retries, logger)
		if err != nil {
			return err
		}
	}
	if withTechPostgres {
		err = Check(ctx, CheckTechPostgres, "tech postgres", retryDelay, retries, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

// Check calls checker function for a job to do.
// name is a display string
// retryDelay is the time to sleep before doing another attempt
// retries is the max number of tries to attempt
func Check(
	ctx context.Context,
	checker func(ctx context.Context, logger zerolog.Logger) error,
	name string,
	retryDelay time.Duration,
	retries int,
	logger zerolog.Logger,
) error {
	infinite := retries == 0
	for retry := 0; retry < retries || infinite; retry++ {
		err := checker(ctx, logger)
		if err == nil {
			return nil
		}

		logger.Info().Err(err).Msgf("%v: %v/%v", name, retry+1, retries)

		t := time.NewTimer(retryDelay)
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
		}
	}

	return fmt.Errorf("check %v failed after %v tries", name, retries)
}

func CheckRedis(ctx context.Context, logger zerolog.Logger) error {
	r, err := redis.NewSession(ctx, 0, logger, 0, 0)
	if err != nil {
		return err
	}

	_ = r.Close()
	return nil
}

func CheckMongo(ctx context.Context, logger zerolog.Logger) error {
	c, err := mongo.NewClient(ctx, 0, 0, logger)
	if err != nil {
		return err
	}

	_ = c.Disconnect(ctx)
	return nil
}

func CheckAMQP(_ context.Context, logger zerolog.Logger) error {
	q, err := amqp.NewConnection(logger, 0, 0)
	if err != nil {
		return err
	}

	_ = q.Close()
	return nil
}

func CheckPostgres(ctx context.Context, _ zerolog.Logger) error {
	p, err := postgres.NewPool(ctx, 0, 0)
	if err != nil {
		return err
	}

	p.Close()
	return nil
}

func CheckTechPostgres(ctx context.Context, _ zerolog.Logger) error {
	p, err := postgres.NewTechMetricsPool(ctx, 0, 0)
	if err != nil {
		return err
	}

	p.Close()
	return nil
}
