package pbehavior

//go:generate mockgen -destination=../../../mocks/lib/canopsis/pbehavior/pbehavior.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior Service,EntityMatcher,ModelProvider,EventManager,ComputedEntityMatcher,Store,EntityTypeResolver

import (
	"context"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/bsm/redislock"
	"time"
)

// Service computes pbehavior timespans and figures out state
// of provided entity by computed data.
type Service interface {
	Compute(ctx context.Context, span timespan.Span) (int, error)
	Recompute(ctx context.Context) error
	RecomputeByID(ctx context.Context, pbehaviorID string) error
	RecomputeByIds(ctx context.Context, pbehaviorIds []string) error

	Resolve(ctx context.Context, entityID string, t time.Time) (ResolveResult, error)
}

// service uses TypeComputer to compute data and TypeResolver to resolve type by this data.
type service struct {
	resolver       TypeResolver
	computer       TypeComputer
	matcher        ComputedEntityMatcher
	store          Store
	workerPoolSize int

	lockClient   redis.LockClient
	lockKey      string
	lockDuration time.Duration
	lockBackoff  time.Duration
	lockRetries  int
}

// NewService creates new service.
func NewService(
	modelProvider ModelProvider,
	matcher ComputedEntityMatcher,
	store Store,
	lockClient redis.LockClient,
) Service {
	return &service{
		store:          store,
		computer:       NewTypeComputer(modelProvider),
		matcher:        matcher,
		workerPoolSize: DefaultPoolSize,
		lockClient:     lockClient,
		lockKey:        redis.RecomputeLockKey,
		lockDuration:   redis.RecomputeLockDuration,
		lockBackoff:    time.Second,
		lockRetries:    10,
	}
}

func (s *service) Compute(ctx context.Context, span timespan.Span) (int, error) {
	currentSpan, err := s.store.GetSpan(ctx)
	if err == nil {
		if currentSpan.To().Sub(span.From()) >= span.To().Sub(span.From())/2 {
			err := s.load(ctx, currentSpan)
			if err != nil {
				return 0, err
			}

			return -1, nil
		}
	} else if !errors.Is(err, ErrNoComputed) {
		return 0, err
	}

	return s.compute(ctx, &span)
}

func (s *service) Recompute(ctx context.Context) error {
	_, err := s.compute(ctx, nil)
	return err
}

func (s *service) RecomputeByID(ctx context.Context, pbehaviorID string) (resErr error) {
	lock, err := s.lockClient.Obtain(ctx, s.lockKey, s.lockDuration, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(s.lockBackoff), s.lockRetries),
	})
	if err != nil {
		return fmt.Errorf("cannot obtain lock: %w", err)
	}

	defer func() {
		err = lock.Release(ctx)
		if err != nil && !errors.Is(err, redislock.ErrLockNotHeld) && resErr == nil {
			resErr = fmt.Errorf("cannot release lock: %w", err)
		}
	}()

	span, err := s.store.GetSpan(ctx)
	if err != nil {
		return err
	}
	res, err := s.computer.Recompute(ctx, span, []string{pbehaviorID})
	if err != nil {
		return err
	}

	computedPbehavior := res[pbehaviorID]
	if computedPbehavior.Name == "" {
		err = s.store.DelComputedPbehavior(ctx, pbehaviorID)
		if err != nil {
			return err
		}
	} else {
		err = s.store.SetComputedPbehavior(ctx, pbehaviorID, computedPbehavior)
		if err != nil {
			return err
		}
	}

	return s.load(ctx, span)
}

func (s *service) RecomputeByIds(ctx context.Context, pbehaviorIds []string) (resErr error) {
	lock, err := s.lockClient.Obtain(ctx, s.lockKey, s.lockDuration, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(s.lockBackoff), s.lockRetries),
	})
	if err != nil {
		return fmt.Errorf("cannot obtain lock: %w", err)
	}

	defer func() {
		err = lock.Release(ctx)
		if err != nil && !errors.Is(err, redislock.ErrLockNotHeld) && resErr == nil {
			resErr = fmt.Errorf("cannot release lock: %w", err)
		}
	}()

	span, err := s.store.GetSpan(ctx)
	if err != nil {
		return err
	}

	res, err := s.computer.Recompute(ctx, span, pbehaviorIds)
	if err != nil {
		return err
	}

	for _, pbehaviorID := range pbehaviorIds {
		computedPbehavior := res[pbehaviorID]

		if computedPbehavior.Name == "" {
			err = s.store.DelComputedPbehavior(ctx, pbehaviorID)
			if err != nil {
				return err
			}
		} else {
			err = s.store.SetComputedPbehavior(ctx, pbehaviorID, computedPbehavior)
			if err != nil {
				return err
			}
		}
	}

	return s.load(ctx, span)
}

func (s *service) Resolve(ctx context.Context, entityID string, t time.Time) (ResolveResult, error) {
	if s.resolver == nil {
		return ResolveResult{}, ErrNoComputed
	}

	pbhIDs, err := s.matcher.Match(ctx, entityID)
	if err != nil || len(pbhIDs) == 0 {
		return ResolveResult{}, err
	}

	return s.resolver.Resolve(ctx, t, pbhIDs)
}

func (s *service) compute(ctx context.Context, span *timespan.Span) (count int, resErr error) {
	lock, err := s.lockClient.Obtain(ctx, s.lockKey, s.lockDuration, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(s.lockBackoff), s.lockRetries),
	})
	if err != nil {
		return 0, fmt.Errorf("cannot obtain lock: %w", err)
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
			return 0, err
		}
		span = &currentSpan
	}

	res, err := s.computer.Compute(ctx, *span)
	if err != nil {
		return 0, err
	}

	err = s.store.SetSpan(ctx, *span)
	if err != nil {
		return 0, err
	}

	err = s.store.SetComputed(ctx, res)
	if err != nil {
		return 0, err
	}

	s.resolver = NewTypeResolver(
		*span,
		res.computedPbehaviors,
		res.typesByID,
		res.defaultActiveType,
	)

	filters, err := getFilters(res.computedPbehaviors)
	if err != nil {
		return 0, err
	}
	err = s.matcher.LoadAll(ctx, filters)
	if err != nil {
		return 0, err
	}

	return len(res.computedPbehaviors), nil
}

func (s *service) load(ctx context.Context, span timespan.Span) error {
	data, err := s.store.GetComputed(ctx)
	if err != nil {
		return err
	}

	s.resolver = NewTypeResolver(
		span,
		data.computedPbehaviors,
		data.typesByID,
		data.defaultActiveType,
	)

	filters, err := getFilters(data.computedPbehaviors)
	if err != nil {
		return err
	}

	return s.matcher.LoadAll(ctx, filters)
}

func getFilters(computed map[string]ComputedPbehavior) (map[string]interface{}, error) {
	filters := make(map[string]interface{}, len(computed))
	var err error
	for id, pbehavior := range computed {
		if len(pbehavior.OldMongoQuery) > 0 {
			filters[id] = pbehavior.OldMongoQuery
		} else {
			filters[id], err = pbehavior.Patten.ToMongoQuery("")
			if err != nil {
				return nil, err
			}
		}
	}

	return filters, nil
}
