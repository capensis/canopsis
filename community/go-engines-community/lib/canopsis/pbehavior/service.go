package pbehavior

//go:generate mockgen -destination=../../../mocks/lib/canopsis/pbehavior/pbehavior.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior Service,ModelProvider,EventManager,Store,EntityTypeResolver,ComputedEntityTypeResolver,TypeComputer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

// Service computes pbehavior timespans and figures out state
// of provided entity by computed data.
type Service interface {
	Compute(ctx context.Context, span timespan.Span) (ComputedEntityTypeResolver, int, error)
	Recompute(ctx context.Context) (ComputedEntityTypeResolver, error)
	RecomputeByIds(ctx context.Context, pbehaviorIds []string) (ComputedEntityTypeResolver, error)
}

// service uses TypeComputer to compute data and TypeResolver to resolve type by this data.
type service struct {
	dbClient       mongo.DbClient
	computer       TypeComputer
	store          Store
	workerPoolSize int

	lockClient   redis.LockClient
	lockKey      string
	lockDuration time.Duration
	lockBackoff  time.Duration
	lockRetries  int

	logger zerolog.Logger
}

// NewService creates new service.
func NewService(
	dbClient mongo.DbClient,
	computer TypeComputer,
	store Store,
	lockClient redis.LockClient,
	logger zerolog.Logger,
) Service {
	return &service{
		dbClient:       dbClient,
		store:          store,
		computer:       computer,
		workerPoolSize: DefaultPoolSize,
		lockClient:     lockClient,
		lockKey:        redis.RecomputeLockKey,
		lockDuration:   redis.RecomputeLockDuration,
		lockBackoff:    time.Second,
		lockRetries:    10,
		logger:         logger,
	}
}

func (s *service) Compute(ctx context.Context, span timespan.Span) (ComputedEntityTypeResolver, int, error) {
	currentSpan, err := s.store.GetSpan(ctx)
	if err == nil {
		if currentSpan.To().Sub(span.From()) >= span.To().Sub(span.From())/2 {
			r, err := s.load(ctx, currentSpan)
			if err != nil {
				return nil, 0, err
			}

			return r, -1, nil
		}
	} else if !errors.Is(err, ErrNoComputed) {
		return nil, 0, err
	}

	return s.compute(ctx, &span)
}

func (s *service) Recompute(ctx context.Context) (ComputedEntityTypeResolver, error) {
	r, _, err := s.compute(ctx, nil)
	return r, err
}

func (s *service) RecomputeByIds(ctx context.Context, pbehaviorIds []string) (_ ComputedEntityTypeResolver, resErr error) {
	lock, err := s.lockClient.Obtain(ctx, s.lockKey, s.lockDuration, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(s.lockBackoff), s.lockRetries),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot obtain lock: %w", err)
	}

	defer func() {
		err = lock.Release(ctx)
		if err != nil && !errors.Is(err, redislock.ErrLockNotHeld) && resErr == nil {
			resErr = fmt.Errorf("cannot release lock: %w", err)
		}
	}()

	span, err := s.store.GetSpan(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.computer.ComputeByIds(ctx, span, pbehaviorIds)
	if err != nil {
		return nil, err
	}

	for _, pbehaviorID := range pbehaviorIds {
		computedPbehavior := res.ComputedPbehaviors[pbehaviorID]

		if computedPbehavior.Name == "" {
			err = s.store.DelComputedPbehavior(ctx, pbehaviorID)
			if err != nil {
				return nil, err
			}
		} else {
			err = s.store.SetComputedPbehavior(ctx, pbehaviorID, computedPbehavior)
			if err != nil {
				return nil, err
			}
		}
	}

	return s.load(ctx, span)
}

func (s *service) compute(ctx context.Context, span *timespan.Span) (_ ComputedEntityTypeResolver, _ int, resErr error) {
	lock, err := s.lockClient.Obtain(ctx, s.lockKey, s.lockDuration, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(s.lockBackoff), s.lockRetries),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("cannot obtain lock: %w", err)
	}

	defer func() {
		err = lock.Release(ctx)
		if err != nil && !errors.Is(err, redislock.ErrLockNotHeld) && resErr == nil {
			resErr = fmt.Errorf("cannot release lock: %w", err)
		}
	}()

	if span == nil {
		currentSpan, err := s.store.GetSpan(ctx)
		if err != nil {
			return nil, 0, fmt.Errorf("cannot get span: %w", err)
		}
		span = &currentSpan
	}

	res, err := s.computer.Compute(ctx, *span)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot compute data: %w", err)
	}

	err = s.store.SetSpan(ctx, *span)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot save span: %w", err)
	}

	err = s.store.SetComputed(ctx, res)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot save computed data: %w", err)
	}

	resolver := NewTypeResolver(
		*span,
		res.ComputedPbehaviors,
		res.TypesByID,
		res.DefaultActiveType,
		s.logger,
	)
	getter := NewComputedEntityGetter(s.dbClient)
	queries := s.getQueries(res.ComputedPbehaviors)
	err = getter.Compute(ctx, queries)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot compute entity getter: %w", err)
	}

	return NewComputedEntityTypeResolver(getter, resolver), len(res.ComputedPbehaviors), nil
}

func (s *service) load(ctx context.Context, span timespan.Span) (ComputedEntityTypeResolver, error) {
	data, err := s.store.GetComputed(ctx)
	if err != nil {
		return nil, err
	}

	resolver := NewTypeResolver(
		span,
		data.ComputedPbehaviors,
		data.TypesByID,
		data.DefaultActiveType,
		s.logger,
	)
	getter := NewComputedEntityGetter(s.dbClient)
	queries := s.getQueries(data.ComputedPbehaviors)
	err = getter.Compute(ctx, queries)
	if err != nil {
		return nil, err
	}

	return NewComputedEntityTypeResolver(getter, resolver), nil
}

func (s *service) getQueries(computed map[string]ComputedPbehavior) []bson.M {
	queries := make([]bson.M, 0, len(computed))
	for id, pbehavior := range computed {
		if len(pbehavior.EntityPattern) == 0 {
			continue
		}

		query, err := db.EntityPatternToMongoQuery(pbehavior.EntityPattern, "")
		if err != nil {
			s.logger.Err(err).Str("pbehavior", id).Msg("pbehavior has invalid pattern")
			continue
		}

		queries = append(queries, query)
	}

	return queries
}
