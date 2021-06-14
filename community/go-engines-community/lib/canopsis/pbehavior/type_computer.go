package pbehavior

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/timespan"
	"github.com/rs/zerolog"
	"github.com/teambition/rrule-go"
)

const (
	getPbehaviorsStep = iota
	getTypesStep
	getExceptionsStep
	getReasonsStep
	computePbehaviorsStep
)

const DefaultPoolSize = 10

// TypeComputer is used to compute periodical behavior timespans for provided interval.
type TypeComputer interface {
	// Compute calculates types for provided timespan.
	Compute(ctx context.Context, span timespan.Span) (*ComputeResult, error)
	Recompute(ctx context.Context, span timespan.Span, pbehaviorID string) (*ComputedPbehavior, error)
}

// typeComputer computes periodical behavior timespans for provided interval.
type typeComputer struct {
	modelProvider ModelProvider
	// workerPoolSize restricts amount of goroutine which can be used during data computing.
	workerPoolSize int
	logger         zerolog.Logger
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
	Types   []computedType `json:"t"`
	Created int64          `json:"c"`
}

// computedType represents computed typeID for determined time span.
type computedType struct {
	Span timespan.Span `json:"s"`
	ID   string        `json:"t"`
}

// ComputeResult represents computed data.
type ComputeResult struct {
	computedPbehaviors map[string]ComputedPbehavior
	typesByID          map[string]*Type
	defaultTypes       map[string]string
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
	logger zerolog.Logger,
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
		logger:         logger,
		workerPoolSize: poolSize,
	}
}

// pbhComputeResult uses for concurrency in compute.
type pbhComputeResult struct {
	id  string
	res ComputedPbehavior
	err error
}

// Compute fetches models from storage and computes data.
func (c *typeComputer) Compute(
	ctx context.Context,
	span timespan.Span,
) (*ComputeResult, error) {
	stepChan := make(chan int, 1)
	defer close(stepChan)
	stepChan <- getPbehaviorsStep

	var (
		pbehaviorsByID map[string]*PBehavior
		models         models
		computed       *ComputeResult
		res            map[string]ComputedPbehavior
	)

	logPrefix := "TypeComputer::Compute:"
	c.logger.Debug().Msgf("%s compute started", logPrefix)
	defer c.logger.Debug().Msgf("%s compute finished", logPrefix)

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug().Msgf("%s received STOP SIGNAL - abort", logPrefix)
			return nil, ctx.Err()
		case step := <-stepChan:
			c.logger.Debug().Msgf("%s step %v - begin", logPrefix, step)
			nextStep := -1
			var err error

			switch step {
			case getPbehaviorsStep:
				pbehaviorsByID, err = c.modelProvider.GetEnabledPbehaviors(ctx)
				if err != nil {
					break
				}

				nextStep = getTypesStep
			case getTypesStep:
				models.typesByID, err = c.modelProvider.GetTypes(ctx)
				if err != nil {
					break
				}

				models.defaultTypes, err = resolveDefaultTypes(models.typesByID)
				if err != nil {
					break
				}

				nextStep = getExceptionsStep
			case getExceptionsStep:
				models.exceptionsByID, err = c.modelProvider.GetExceptions(ctx)
				if err != nil {
					break
				}

				nextStep = getReasonsStep
			case getReasonsStep:
				models.reasonsByID, err = c.modelProvider.GetReasons(ctx)
				if err != nil {
					break
				}

				nextStep = computePbehaviorsStep
			case computePbehaviorsStep:
				res, err = c.computePbehaviors(ctx, span, pbehaviorsByID, models)
				if err != nil {
					break
				}

				computed = &ComputeResult{
					computedPbehaviors: res,
					typesByID:          models.typesByID,
					defaultTypes:       models.defaultTypes,
				}
			}

			if err != nil {
				c.logger.Debug().Msgf("%s step %v - error: %v", logPrefix, step, err)
				return nil, err
			}

			c.logger.Debug().Msgf("%s step %v - end", logPrefix, step)
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
	pbehaviorID string,
) (*ComputedPbehavior, error) {
	stepChan := make(chan int, 1)
	defer close(stepChan)
	stepChan <- getPbehaviorsStep

	var pbehavior *PBehavior
	var models models
	var computed *ComputedPbehavior
	logPrefix := fmt.Sprintf("TypeComputer::Recompute:%s:", pbehaviorID)

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug().Msgf("%s received STOP SIGNAL - abort", logPrefix)
			return nil, ctx.Err()
		case step := <-stepChan:
			c.logger.Debug().Msgf("%s step %v - begin", logPrefix, step)
			nextStep := -1
			var err error

			switch step {
			case getPbehaviorsStep:
				pbehavior, err = c.modelProvider.GetEnabledPbehavior(ctx, pbehaviorID)
				if err != nil || pbehavior == nil {
					break
				}

				nextStep = getTypesStep
			case getTypesStep:
				models.typesByID, err = c.modelProvider.GetTypes(ctx)
				if err != nil {
					break
				}

				models.defaultTypes, err = resolveDefaultTypes(models.typesByID)
				if err != nil {
					break
				}

				nextStep = getExceptionsStep
			case getExceptionsStep:
				models.exceptionsByID, err = c.modelProvider.GetExceptions(ctx)
				if err != nil {
					break
				}

				nextStep = getReasonsStep
			case getReasonsStep:
				models.reasonsByID, err = c.modelProvider.GetReasons(ctx)
				if err != nil {
					break
				}

				nextStep = computePbehaviorsStep
			case computePbehaviorsStep:
				computed, err = c.computePbehavior(pbehavior, span, models)
			}

			if err != nil {
				c.logger.Debug().Msgf("%s step %v - error: %v", logPrefix, step, err)
				return nil, err
			}

			c.logger.Debug().Msgf("%s step %v - end", logPrefix, step)
			if nextStep == -1 {
				return computed, nil
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

// computePbehaviors calculates Types for provided date concurrently.
func (c *typeComputer) computePbehaviors(
	ctx context.Context,
	span timespan.Span,
	pbehaviorsByID map[string]*PBehavior,
	models models,
) (map[string]ComputedPbehavior, error) {
	resCh := c.runWorkers(ctx, span, pbehaviorsByID, models)
	res := make(map[string]ComputedPbehavior)
	var err error

	for v := range resCh {
		if v.err == nil {
			res[v.id] = v.res
		} else {
			err = v.err
		}
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

// runWorkers creates res chan and runs workers to fill it.
func (c *typeComputer) runWorkers(
	ctx context.Context,
	span timespan.Span,
	pbehaviorsByID map[string]*PBehavior,
	models models,
) <-chan pbhComputeResult {
	resCh := make(chan pbhComputeResult)
	workerChan := c.createWorkerCh(ctx, pbehaviorsByID)

	go func() {
		defer close(resCh)
		wg := sync.WaitGroup{}

		for worker := 0; worker < c.workerPoolSize; worker++ {
			wg.Add(1)

			go func(ctx context.Context, id int) {
				workerCtx, cancel := context.WithCancel(ctx)
				defer cancel()
				defer wg.Done()

				logPrefix := fmt.Sprintf("TypeComputer::Compute: step %v - worker %d", computePbehaviorsStep, id)
				c.logger.Debug().Msgf("%s started\n", logPrefix)
				defer c.logger.Debug().Msgf("%s finished\n", logPrefix)

				for {
					select {
					case <-workerCtx.Done():
						resCh <- pbhComputeResult{
							err: workerCtx.Err(),
						}
						c.logger.Debug().Msgf("%s context was cancelled\n", logPrefix)
						return
					case p, ok := <-workerChan:
						if !ok {
							return
						}

						c.logger.Debug().Msgf("%s received message from worker channel\n", logPrefix)
						res, err := c.computePbehavior(p, span, models)
						if err != nil {
							resCh <- pbhComputeResult{
								err: err,
							}
						} else if res != nil {
							resCh <- pbhComputeResult{
								id:  p.ID,
								res: *res,
							}

							c.logger.Debug().Msgf("%s send result to channel\n", logPrefix)
						}

						c.logger.Debug().Msgf("%s finished message processing\n", logPrefix)
					}
				}
			}(ctx, worker)
		}

		wg.Wait()
	}()

	return resCh
}

// createWorkerCh creates worker chan and fills it.
func (c *typeComputer) createWorkerCh(ctx context.Context, pbehaviorsByID map[string]*PBehavior) <-chan *PBehavior {
	workerChan := make(chan *PBehavior)

	go func() {
		defer close(workerChan)
		defer c.logger.Debug().Msgf("TypeComputer::Compute: %v - send is over, waiting workers to complete jobs\n", computePbehaviorsStep)

		for _, pbehavior := range pbehaviorsByID {
			select {
			case <-ctx.Done():
				c.logger.Debug().Msgf("TypeComputer::Compute: %v - stop workers\n", computePbehaviorsStep)
				return
			default:
				c.logger.Debug().Msgf("TypeComputer::Compute: %v - send a work to worker channel\n", computePbehaviorsStep)
				workerChan <- pbehavior
			}
		}
	}()

	return workerChan
}

// computePbehavior calculates Types for provided timespan.
func (c *typeComputer) computePbehavior(
	pbehavior *PBehavior,
	span timespan.Span,
	models models,
) (*ComputedPbehavior, error) {
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
		event = NewEvent(pbehavior.Start.In(location), stop)
	} else {
		rOption, err := rrule.StrToROption(pbehavior.RRule)
		if err != nil {
			return nil, err
		}

		event = NewRecEvent(pbehavior.Start.In(location), stop, rOption)
	}

	resByExdate, err := c.computeByExdate(pbehavior, event, span, models)
	if err != nil {
		return nil, err
	}

	resByRrule, err := c.computeByRrule(pbehavior, event, span)
	if err != nil {
		return nil, err
	}

	resByActiveType, err := c.computeByActiveType(pbehavior, event, span, models)
	if err != nil {
		return nil, err
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

		return &ComputedPbehavior{
			Filter:  pbehavior.Filter,
			Name:    pbehavior.Name,
			Reason:  reasonName,
			Types:   res,
			Created: pbehavior.Created.Unix(),
		}, nil
	}

	return nil, nil
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
		from := maxTime(span.From(), exdate.Begin.In(location))
		to := minTime(span.To(), exdate.End.In(location))
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
