package pbehavior

//go:generate mockgen -destination=../../../mocks/lib/canopsis/pbehavior/service.go git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior Service

import (
	"context"
	"encoding/json"
	libtypes "git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/timespan"
	"github.com/rs/zerolog"
	"time"
)

// Service computes pbehavior timespans and figures out state
// of provided entity by computed data.
type Service interface {
	Resolve(ctx context.Context, entity *libtypes.Entity, t time.Time) (ResolveResult, error)
	Compute(ctx context.Context, span timespan.Span) error
	Recompute(ctx context.Context, pbehaviorID string) error
	GetSpan() timespan.Span
	GetPbehaviorStatus(ctx context.Context, pbehaviorIDs []string, t time.Time) (map[string]bool, error)
}

// service uses TypeComputer to compute data and TypeResolver to resolve type by this data.
type service struct {
	resolver       *typeResolver
	computer       TypeComputer
	matcher        EntityMatcher
	logger         zerolog.Logger
	workerPoolSize int
}

// NewService creates new service.
func NewService(
	modelProvider ModelProvider,
	matcher EntityMatcher,
	logger zerolog.Logger,
	workerPoolSize ...int,
) Service {
	poolSize := DefaultPoolSize

	if len(workerPoolSize) == 1 {
		poolSize = workerPoolSize[0]
	} else if len(workerPoolSize) > 1 {
		panic("too much arguments")
	}

	return &service{
		resolver: NewTypeResolver(
			matcher,
			timespan.Span{},
			nil,
			nil,
			"",
			logger,
		),
		computer:       NewTypeComputer(modelProvider, logger, poolSize),
		matcher:        matcher,
		logger:         logger,
		workerPoolSize: poolSize,
	}
}

func (s *service) Compute(ctx context.Context, span timespan.Span) error {
	res, err := s.computer.Compute(ctx, span)
	if err != nil {
		return err
	}

	s.resolver = NewTypeResolver(
		s.matcher,
		span,
		res.computedPbehaviors,
		res.typesByID,
		res.defaultTypes[TypeActive],
		s.logger,
	)

	return nil
}

func (s *service) Recompute(ctx context.Context, pbehaviorID string) error {
	res, err := s.computer.Recompute(ctx, s.GetSpan(), pbehaviorID)
	if err != nil {
		return err
	}

	s.resolver.UpdateData(pbehaviorID, res)

	return nil
}

func (s *service) Resolve(ctx context.Context, entity *libtypes.Entity, t time.Time) (ResolveResult, error) {
	return s.resolver.Resolve(ctx, entity, t)
}

func (s *service) GetPbehaviorStatus(
	ctx context.Context,
	pbehaviorIDs []string,
	t time.Time,
) (map[string]bool, error) {
	return s.resolver.GetPbehaviorStatus(ctx, pbehaviorIDs, t)
}

func (s *service) GetSpan() timespan.Span {
	return s.resolver.GetSpan()
}

func (s *service) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.resolver)
}

func (s *service) UnmarshalJSON(b []byte) error {
	resolver := NewTypeResolver(
		s.matcher,
		timespan.Span{},
		nil,
		nil,
		"",
		s.logger,
	)
	err := json.Unmarshal(b, resolver)
	if err != nil {
		return err
	}

	s.resolver = resolver

	return nil
}
