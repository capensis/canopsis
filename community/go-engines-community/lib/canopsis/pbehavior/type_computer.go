package pbehavior

//go:generate easyjson -no_std_marshalers

import (
	"context"
	"fmt"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/teambition/rrule-go"
	"golang.org/x/sync/errgroup"
)

const (
	getPbehaviorsStep = iota
	getTypesStep
	getExceptionsStep
	getReasonsStep
	computePbehaviorsStep
)

const DefaultPoolSize = 100

// TypeComputer is used to compute periodical behavior timespans for provided interval.
type TypeComputer interface {
	// Compute calculates types for provided timespan.
	Compute(ctx context.Context, span timespan.Span) (ComputeResult, error)
	Recompute(ctx context.Context, span timespan.Span, pbehaviorIds []string) (map[string]ComputedPbehavior, error)
}

// typeComputer computes periodical behavior timespans for provided interval.
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
//easyjson:json
type ComputedPbehavior struct {
	Name    string         `json:"n"`
	Reason  string         `json:"r"`
	Filter  string         `json:"f"`
	Types   []computedType `json:"t"`
	Created int64          `json:"c"`
}

// computedType represents computed typeID for determined time span.
type computedType struct {
	Span timespan.Span `json:"s"`
	ID   string        `json:"t"`
}

//easyjson:json
type Types struct {
	T map[string]*Type
}

// ComputeResult represents computed data.
type ComputeResult struct {
	computedPbehaviors map[string]ComputedPbehavior
	typesByID          map[string]*Type
	defaultActiveType  string
}

// models contains all required models for computing.
type models struct {
	typesByID      map[string]*Type
	defaultTypes   map[string]string
	exceptionsByID map[string]*Exception
	reasonsByID    map[string]*Reason
}

// NewTypeComputer creates new type resolver.
func NewTypeComputer(
	modelProvider ModelProvider,
	workerPoolSize ...int,
) TypeComputer {
	poolSize := DefaultPoolSize

	if len(workerPoolSize) == 1 {
		poolSize = workerPoolSize[0]
	} else if len(workerPoolSize) > 1 {
		panic("too much arguments")
	}

	return &typeComputer{
		modelProvider:  modelProvider,
		workerPoolSize: poolSize,
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
	stepChan := make(chan int, 1)
	defer close(stepChan)
	stepChan <- getPbehaviorsStep

	var (
		pbehaviorsByID map[string]*PBehavior
		models         models
		computed       ComputeResult
		res            map[string]ComputedPbehavior
	)

	for {
		select {
		case <-ctx.Done():
			return computed, nil
		case step := <-stepChan:
			nextStep := -1
			var err error

			switch step {
			case getPbehaviorsStep:
				pbehaviorsByID, err = c.modelProvider.GetEnabledPbehaviors(ctx)
				if err != nil {
					err = fmt.Errorf("cannot fetch pbehaviors: %w", err)
					break
				}

				nextStep = getTypesStep
			case getTypesStep:
				models.typesByID, err = c.modelProvider.GetTypes(ctx)
				if err != nil {
					err = fmt.Errorf("cannot fetch pbehavior types: %w", err)
					break
				}

				models.defaultTypes, err = resolveDefaultTypes(models.typesByID)
				if err != nil {
					err = fmt.Errorf("cannot fetch defult pbehavior types: %w", err)
					break
				}

				nextStep = getExceptionsStep
			case getExceptionsStep:
				models.exceptionsByID, err = c.modelProvider.GetExceptions(ctx)
				if err != nil {
					err = fmt.Errorf("cannot fetch pbehavior exeptions: %w", err)
					break
				}

				nextStep = getReasonsStep
			case getReasonsStep:
				models.reasonsByID, err = c.modelProvider.GetReasons(ctx)
				if err != nil {
					err = fmt.Errorf("cannot fetch pbehavior reasons: %w", err)
					break
				}

				nextStep = computePbehaviorsStep
			case computePbehaviorsStep:
				res, err = c.runWorkers(ctx, span, pbehaviorsByID, models)
				if err != nil {
					break
				}

				computed = ComputeResult{
					computedPbehaviors: res,
					typesByID:          models.typesByID,
					defaultActiveType:  models.defaultTypes[TypeActive],
				}
			}

			if err != nil {
				return computed, err
			}

			if nextStep == -1 {
				return computed, nil
			}

			stepChan <- nextStep
		}
	}
}

func (c *typeComputer) Recompute(
	ctx context.Context,
	span timespan.Span,
	pbehaviorIds []string,
) (map[string]ComputedPbehavior, error) {
	stepChan := make(chan int, 1)
	defer close(stepChan)
	stepChan <- getPbehaviorsStep

	var (
		pbehaviorsByID map[string]*PBehavior
		models         models
		res            map[string]ComputedPbehavior
	)

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case step := <-stepChan:
			nextStep := -1
			var err error

			switch step {
			case getPbehaviorsStep:
				pbehaviorsByID, err = c.modelProvider.GetEnabledPbehaviorsByIds(ctx, pbehaviorIds)
				if err != nil {
					err = fmt.Errorf("cannot fetch pbehavior: %w", err)
					break
				}
				if len(pbehaviorsByID) == 0 {
					break
				}

				nextStep = getTypesStep
			case getTypesStep:
				models.typesByID, err = c.modelProvider.GetTypes(ctx)
				if err != nil {
					err = fmt.Errorf("cannot fetch pbehavior types: %w", err)
					break
				}

				models.defaultTypes, err = resolveDefaultTypes(models.typesByID)
				if err != nil {
					err = fmt.Errorf("cannot fetch default pbehavior types: %w", err)
					break
				}

				nextStep = getExceptionsStep
			case getExceptionsStep:
				models.exceptionsByID, err = c.modelProvider.GetExceptions(ctx)
				if err != nil {
					err = fmt.Errorf("cannot fetch pbehavior exceptions: %w", err)
					break
				}

				nextStep = getReasonsStep
			case getReasonsStep:
				models.reasonsByID, err = c.modelProvider.GetReasons(ctx)
				if err != nil {
					err = fmt.Errorf("cannot fetch pbehavior reasons: %w", err)
					break
				}

				nextStep = computePbehaviorsStep
			case computePbehaviorsStep:
				res, err = c.runWorkers(ctx, span, pbehaviorsByID, models)
				if err != nil {
					break
				}
			}

			if err != nil {
				return nil, err
			}

			if nextStep == -1 {
				return res, nil
			}

			stepChan <- nextStep
		}
	}
}

// resolveDefaultTypes finds default types which uses :
// - active default type if there aren't any behaviors
// - inactive default type if there is behavior with active type
func resolveDefaultTypes(typesByID map[string]*Type) (map[string]string, error) {
	defaultTypesByType := map[string]*Type{
		TypeInactive: nil,
		TypeActive:   nil,
	}

	for _, t := range typesByID {
		for dt := range defaultTypesByType {
			if t.Type == dt && (defaultTypesByType[dt] == nil ||
				defaultTypesByType[dt].Priority > t.Priority) {
				defaultTypesByType[dt] = t
			}
		}
	}

	defaultTypes := make(map[string]string)
	for dt := range defaultTypesByType {
		if defaultTypesByType[dt] == nil {
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
	pbehaviorsByID map[string]*PBehavior,
	models models,
) (map[string]ComputedPbehavior, error) {
	resCh := make(chan pbhComputeResult)
	workerChan := make(chan *PBehavior)
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
				res, err := c.computePbehavior(p, span, models)
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
	pbehavior *PBehavior,
	span timespan.Span,
	models models,
) (ComputedPbehavior, error) {
	var event Event
	location := span.From().Location()
	var stop time.Time
	if pbehavior.Stop == nil {
		stop = span.To()
	} else {
		stop = pbehavior.Stop.Time
	}
	stop = stop.In(location)

	if pbehavior.RRule == "" {
		event = NewEvent(pbehavior.Start.Time.In(location), stop)
	} else {
		rOption, err := rrule.StrToROption(pbehavior.RRule)
		if err != nil {
			return ComputedPbehavior{}, err
		}

		event = NewRecEvent(pbehavior.Start.Time.In(location), stop, rOption)
	}

	resByExdate, err := c.computeByExdate(pbehavior, event, span, models)
	if err != nil {
		return ComputedPbehavior{}, err
	}

	resByRrule, err := c.computeByRrule(pbehavior, event, span)
	if err != nil {
		return ComputedPbehavior{}, err
	}

	resByActiveType, err := c.computeByActiveType(pbehavior, event, span, models)
	if err != nil {
		return ComputedPbehavior{}, err
	}

	res := resByExdate
	res = append(res, resByRrule...)
	res = append(res, resByActiveType...)

	if len(res) > 0 {
		reason, ok := models.reasonsByID[pbehavior.Reason]
		reasonName := pbehavior.Reason
		if ok {
			reasonName = reason.Name
		}

		return ComputedPbehavior{
			Filter:  pbehavior.Filter,
			Name:    pbehavior.Name,
			Reason:  reasonName,
			Types:   res,
			Created: pbehavior.Created.Unix(),
		}, nil
	}

	return ComputedPbehavior{}, nil
}

// computeByExdate returns all time spans for pbehavior on the date which are defined by exdate.
func (c *typeComputer) computeByExdate(
	pbehavior *PBehavior,
	event Event,
	span timespan.Span,
	models models,
) ([]computedType, error) {
	exdateList, err := c.getSortedExdate(pbehavior, models)
	if err != nil {
		return nil, err
	}

	location := event.span.From().Location()
	res := make([]computedType, 0)
	for _, exdate := range exdateList {
		from := maxTime(span.From(), exdate.Begin.Time.In(location))
		to := minTime(span.To(), exdate.End.Time.In(location))
		if from.After(to) {
			continue
		}

		view := timespan.New(from, to)
		timespans, err := GetTimeSpans(event, view)
		if err != nil {
			return nil, err
		}

		for i := range timespans {
			res = append(res, computedType{
				Span: timespans[i],
				ID:   exdate.Type,
			})
		}

	}

	return res, nil
}

// computeByRrule returns all time spans for pbehavior on the date which are defined by rrule.
func (c *typeComputer) computeByRrule(
	pbehavior *PBehavior,
	event Event,
	span timespan.Span,
) ([]computedType, error) {
	timespans, err := GetTimeSpans(event, span)
	if err != nil {
		return nil, err
	}

	res := make([]computedType, len(timespans))

	for i := range timespans {
		res[i].Span = timespans[i]
		res[i].ID = pbehavior.Type
	}

	return res, nil
}

// computeByActiveType returns time span with default inactive type if pbehavior has
// active type and pbehavior is in action for the date.
func (c *typeComputer) computeByActiveType(
	pbehavior *PBehavior,
	event Event,
	span timespan.Span,
	models models,
) ([]computedType, error) {
	t, ok := models.typesByID[pbehavior.Type]
	if !ok {
		return nil, fmt.Errorf("unknown type %s of pbh %+v", pbehavior.Type, pbehavior)
	}

	if t.Type != TypeActive {
		return nil, nil
	}

	res := make([]computedType, 0)
	// Check each day in the span if behavior is in action for this day.
	for date := dateOf(span.From()); date.Before(span.To()); date = date.AddDate(0, 0, 1) {
		// Interval from midnight to next day midnight.
		dateSpan := timespan.New(date, date.AddDate(0, 0, 1))
		timespans, err := GetTimeSpans(event, dateSpan)
		if err != nil {
			return nil, err
		}

		if len(timespans) > 0 {
			res = append(res, computedType{
				Span: timespan.New(
					maxTime(dateSpan.From(), span.From()),
					minTime(dateSpan.To(), span.To()),
				),
				ID: models.defaultTypes[TypeInactive],
			})
		}
	}

	return res, nil
}

// getSortedExdate returns pbehavior exdate list which is sorted by its type priority.
func (c *typeComputer) getSortedExdate(
	pbehavior *PBehavior,
	models models,
) ([]*Exdate, error) {
	res := make([]*Exdate, len(pbehavior.Exdates))
	for i := range pbehavior.Exdates {
		res[i] = &pbehavior.Exdates[i]
	}

	for _, id := range pbehavior.Exceptions {
		if exception, ok := models.exceptionsByID[id]; ok {
			for i := range exception.Exdates {
				res = append(res, &exception.Exdates[i])
			}
		} else {
			return nil, fmt.Errorf("unknown exception %v", id)
		}
	}

	var unknownTypeErr error
	getPriorityByID := func(id string) int {
		t, ok := models.typesByID[id]
		if !ok {
			unknownTypeErr = fmt.Errorf("unknown type %v", id)
			return 0
		}

		return t.Priority
	}

	sort.Slice(res, func(i, j int) bool {
		return getPriorityByID(res[i].Type) > getPriorityByID(res[j].Type)
	})

	if unknownTypeErr != nil {
		return nil, unknownTypeErr
	}

	return res, nil
}
