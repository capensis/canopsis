package pbehavior

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"golang.org/x/sync/errgroup"
)

const DefaultPoolSize = 100

// TypeComputer is used to compute all periodical behaviors' timespans for provided interval.
type TypeComputer interface {
	// Compute calculates types for provided timespan.
	Compute(ctx context.Context, span timespan.Span) (ComputeResult, error)
	ComputeByIds(ctx context.Context, span timespan.Span, pbehaviorIds []string) (ComputeResult, error)
}

type typeComputer struct {
	modelProvider ModelProvider
	// workerPoolSize restricts amount of goroutine which can be used during data computing.
	workerPoolSize int
}

// ComputedPbehavior represents all computed types for periodical behavior.
// Computed types are sorted:
// - time spans which are defined by exdate
// - time spans which are defined by rrule
// - time spans which are defined by default inactive interval of active pbehavior
// For example, for active daily periodical behavior at 10:00-12:00 and date 2020-06-01:
// [2020-06-01T10:00, 2020-06-01T12:00] ActiveTypeID
// [2020-06-01T00:00, 2020-06-02T00:00] InactiveTypeID
type ComputedPbehavior struct {
	Name    string         `json:"n"`
	Reason  string         `json:"r"`
	Filter  string         `json:"f"`
	Types   []ComputedType `json:"t"`
	Created int64          `json:"c"`
	Color   string         `json:"-"`
}

// ComputeResult represents computed data.
type ComputeResult struct {
	ComputedPbehaviors map[string]ComputedPbehavior
	TypesByID          map[string]Type
	DefaultActiveType  string
}

// models contains all required models for computing.
type models struct {
	typesByID      map[string]Type
	defaultTypes   map[string]string
	exceptionsByID map[string]Exception
	reasonsByID    map[string]Reason
}

// NewTypeComputer creates new type resolver.
func NewTypeComputer(
	modelProvider ModelProvider,
) TypeComputer {
	return &typeComputer{
		modelProvider:  modelProvider,
		workerPoolSize: DefaultPoolSize,
	}
}

// pbhComputeResult uses for concurrency in compute.
type pbhComputeResult struct {
	id  string
	res ComputedPbehavior
}

// Compute fetches models from storage and computes data.
func (c *typeComputer) Compute(
	ctx context.Context,
	span timespan.Span,
) (ComputeResult, error) {
	pbehaviorsByID, err := c.modelProvider.GetEnabledPbehaviors(ctx, span)
	if err != nil {
		return ComputeResult{}, fmt.Errorf("cannot fetch pbehaviors: %w", err)
	}

	models := models{}
	models.typesByID, err = c.modelProvider.GetTypes(ctx)
	if err != nil {
		return ComputeResult{}, fmt.Errorf("cannot fetch pbehavior types: %w", err)
	}

	models.defaultTypes, err = ResolveDefaultTypes(models.typesByID)
	if err != nil {
		return ComputeResult{}, fmt.Errorf("cannot fetch defult pbehavior types: %w", err)
	}

	models.exceptionsByID, err = c.modelProvider.GetExceptions(ctx)
	if err != nil {
		return ComputeResult{}, fmt.Errorf("cannot fetch pbehavior exeptions: %w", err)
	}

	models.reasonsByID, err = c.modelProvider.GetReasons(ctx)
	if err != nil {
		return ComputeResult{}, fmt.Errorf("cannot fetch pbehavior reasons: %w", err)
	}

	res, err := c.runWorkers(ctx, span, pbehaviorsByID, models)
	if err != nil {
		return ComputeResult{}, err
	}

	computed := ComputeResult{
		ComputedPbehaviors: res,
		TypesByID:          models.typesByID,
		DefaultActiveType:  models.defaultTypes[TypeActive],
	}

	return computed, nil
}

func (c *typeComputer) ComputeByIds(
	ctx context.Context,
	span timespan.Span,
	pbehaviorIds []string,
) (ComputeResult, error) {
	pbehaviorsByID, err := c.modelProvider.GetEnabledPbehaviorsByIds(ctx, pbehaviorIds, span)
	if err != nil {
		return ComputeResult{}, fmt.Errorf("cannot fetch pbehaviors: %w", err)
	}

	models := models{}
	models.typesByID, err = c.modelProvider.GetTypes(ctx)
	if err != nil {
		return ComputeResult{}, fmt.Errorf("cannot fetch pbehavior types: %w", err)
	}

	models.defaultTypes, err = ResolveDefaultTypes(models.typesByID)
	if err != nil {
		return ComputeResult{}, fmt.Errorf("cannot fetch defult pbehavior types: %w", err)
	}

	models.exceptionsByID, err = c.modelProvider.GetExceptions(ctx)
	if err != nil {
		return ComputeResult{}, fmt.Errorf("cannot fetch pbehavior exeptions: %w", err)
	}

	models.reasonsByID, err = c.modelProvider.GetReasons(ctx)
	if err != nil {
		return ComputeResult{}, fmt.Errorf("cannot fetch pbehavior reasons: %w", err)
	}

	res, err := c.runWorkers(ctx, span, pbehaviorsByID, models)
	if err != nil {
		return ComputeResult{}, err
	}

	computed := ComputeResult{
		ComputedPbehaviors: res,
		TypesByID:          models.typesByID,
		DefaultActiveType:  models.defaultTypes[TypeActive],
	}

	return computed, nil
}

// ResolveDefaultTypes finds default types which uses :
// - active default type if there aren't any behaviors
// - inactive default type if there is behavior with active type
func ResolveDefaultTypes(typesByID map[string]Type) (map[string]string, error) {
	defaultTypesByType := map[string]Type{
		TypeInactive: {},
		TypeActive:   {},
	}

	for _, t := range typesByID {
		for dt := range defaultTypesByType {
			if t.Type == dt && (defaultTypesByType[dt].ID == "" ||
				defaultTypesByType[dt].Priority > t.Priority) {
				defaultTypesByType[dt] = t
			}
		}
	}

	defaultTypes := make(map[string]string)
	for dt := range defaultTypesByType {
		if defaultTypesByType[dt].ID == "" {
			return nil, fmt.Errorf("no default type %v", dt)
		}
		defaultTypes[dt] = defaultTypesByType[dt].ID
	}

	return defaultTypes, nil
}

// runWorkers creates res chan and runs workers to fill it.
func (c *typeComputer) runWorkers(
	ctx context.Context,
	span timespan.Span,
	pbehaviorsByID map[string]PBehavior,
	models models,
) (map[string]ComputedPbehavior, error) {
	eventComputer := NewEventComputer(models.typesByID, models.defaultTypes)
	resCh := make(chan pbhComputeResult)
	workerChan := make(chan PBehavior)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(workerChan)

		for _, pbehavior := range pbehaviorsByID {
			select {
			case <-ctx.Done():
				return nil
			case workerChan <- pbehavior:
			}
		}

		return nil
	})

	for worker := 0; worker < c.workerPoolSize; worker++ {
		g.Go(func() error {
			for p := range workerChan {
				res, err := c.computePbehavior(p, span, eventComputer, models)
				if err != nil {
					return err
				}

				if len(res.Types) > 0 {
					resCh <- pbhComputeResult{
						id:  p.ID,
						res: res,
					}
				}
			}

			return nil
		})
	}

	go func() {
		_ = g.Wait() // check error in wrapper func
		close(resCh)
	}()

	res := make(map[string]ComputedPbehavior)
	for v := range resCh {
		res[v.id] = v.res
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return res, nil
}

// computePbehavior calculates Types for provided timespan.
func (c *typeComputer) computePbehavior(
	pbehavior PBehavior,
	span timespan.Span,
	eventComputer EventComputer,
	models models,
) (ComputedPbehavior, error) {
	var start, end types.CpsTime
	if pbehavior.Start != nil {
		start = *pbehavior.Start
	}
	if pbehavior.Stop != nil {
		end = *pbehavior.Stop
	}
	exdates, err := c.getExdates(pbehavior, models)
	if err != nil {
		return ComputedPbehavior{}, err
	}
	params := PbhEventParams{
		Start:   start,
		End:     end,
		RRule:   pbehavior.RRule,
		Type:    pbehavior.Type,
		Exdates: exdates,
	}
	compitedTypes, err := eventComputer.Compute(params, span)
	if err != nil {
		return ComputedPbehavior{}, err
	}

	if len(compitedTypes) > 0 {
		reason, ok := models.reasonsByID[pbehavior.Reason]
		reasonName := pbehavior.Reason
		if ok {
			reasonName = reason.Name
		}

		return ComputedPbehavior{
			Filter:  pbehavior.Filter,
			Name:    pbehavior.Name,
			Reason:  reasonName,
			Types:   compitedTypes,
			Created: pbehavior.Created.Unix(),
			Color:   pbehavior.Color,
		}, nil
	}

	return ComputedPbehavior{}, nil
}

// getExdates returns pbehavior exdate list.
func (c *typeComputer) getExdates(
	pbehavior PBehavior,
	models models,
) ([]Exdate, error) {
	res := make([]Exdate, len(pbehavior.Exdates))
	for i := range pbehavior.Exdates {
		res[i] = pbehavior.Exdates[i]
	}

	for _, id := range pbehavior.Exceptions {
		if exception, ok := models.exceptionsByID[id]; ok {
			for i := range exception.Exdates {
				res = append(res, exception.Exdates[i])
			}
		} else {
			return nil, fmt.Errorf("unknown exception %v", id)
		}
	}

	return res, nil
}
