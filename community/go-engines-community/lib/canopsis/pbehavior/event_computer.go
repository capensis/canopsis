package pbehavior

import (
	"fmt"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/teambition/rrule-go"
)

// EventComputer is used to compute periodical behavior timespans for provided interval.
type EventComputer interface {
	Compute(params PbhEventParams, span timespan.Span) ([]ComputedType, error)
}

type eventComputer struct {
	typesByID    map[string]Type
	defaultTypes map[string]string
}

func NewEventComputer(typesByID map[string]Type, defaultTypes map[string]string) EventComputer {
	return &eventComputer{
		typesByID:    typesByID,
		defaultTypes: defaultTypes,
	}
}

type PbhEventParams struct {
	ID      string
	Start   datetime.CpsTime
	End     datetime.CpsTime
	RRule   string
	Type    string
	Exdates []Exdate
}

// ComputedType represents type for determined time span.
type ComputedType struct {
	ID    string        `json:"t"`
	Span  timespan.Span `json:"s"`
	Color string        `json:"-"`
}

// Compute calculates types for provided timespan.
func (c *eventComputer) Compute(
	params PbhEventParams,
	span timespan.Span,
) ([]ComputedType, error) {
	var event Event
	location := span.From().Location()
	var stop time.Time
	if params.End.Unix() <= 0 {
		stop = span.To()
	} else {
		stop = params.End.Time
	}
	stop = stop.In(location)

	if params.RRule == "" {
		event = NewEvent(params.Start.Time.In(location), stop)
	} else {
		rOption, err := rrule.StrToROption(params.RRule)
		if err != nil {
			return nil, err
		}

		event = NewRecEvent(params.Start.Time.In(location), stop, rOption)
	}

	err := c.sortExdates(params.Exdates)
	if err != nil {
		return nil, err
	}
	computed, err := c.computeByRrule(event, span, c.typesByID[params.Type], params.Exdates)
	if err != nil {
		return nil, err
	}

	computedByActiveType, err := c.computeByActiveType(event, span, params.Type)
	if err != nil {
		return nil, err
	}

	computed = append(computed, computedByActiveType...)

	return computed, nil
}

// computeByRrule returns all time spans for pbehavior on the date which are defined by rrule with exdates.
func (c *eventComputer) computeByRrule(
	event Event,
	span timespan.Span,
	t Type,
	exdates []Exdate,
) ([]ComputedType, error) {
	location := event.span.From().Location()
	eventTimespans, err := GetTimeSpans(event, span)
	if err != nil {
		return nil, err
	}

	computed := make([]ComputedType, len(eventTimespans))
	for i := range eventTimespans {
		computed[i].Span = eventTimespans[i]
		computed[i].ID = t.ID
		computed[i].Color = t.Color
	}

	computedByExdate := make([]ComputedType, 0)

	for _, exdate := range exdates {
		computedWithoutExdates := make([]ComputedType, 0, len(computed))

		for _, curr := range computed {
			from := utils.MaxTime(curr.Span.From(), exdate.Begin.Time.In(location))
			to := utils.MinTime(curr.Span.To(), exdate.End.Time.In(location))
			if from.After(to) {
				computedWithoutExdates = append(computedWithoutExdates, curr)
				continue
			}

			exdateSpan := timespan.New(from, to)
			diffs := curr.Span.Diff(exdateSpan)
			for _, diff := range diffs {
				computedWithoutExdates = append(computedWithoutExdates, ComputedType{
					Span:  diff,
					ID:    curr.ID,
					Color: curr.Color,
				})
			}

			computedByExdate = append(computedByExdate, ComputedType{
				ID:   exdate.Type,
				Span: exdateSpan,
			})
		}

		computed = computedWithoutExdates
		if len(computed) == 0 {
			break
		}
	}

	computed = append(computed, computedByExdate...)

	return computed, nil
}

// computeByActiveType returns time span with default inactive type if pbehavior has
// active type and pbehavior is in action for the date.
func (c *eventComputer) computeByActiveType(
	event Event,
	span timespan.Span,
	typeID string,
) ([]ComputedType, error) {
	t, ok := c.typesByID[typeID]
	if !ok {
		return nil, fmt.Errorf("unknown type %s", typeID)
	}

	if t.Type != TypeActive {
		return nil, nil
	}

	computed := make([]ComputedType, 0)
	// Check each day in the span if behavior is in action for this day.
	for date := utils.DateOf(span.From()); date.Before(span.To()); date = date.AddDate(0, 0, 1) {
		// Interval from midnight to next day midnight.
		dateSpan := timespan.New(date, date.AddDate(0, 0, 1))
		eventTimespans, err := GetTimeSpans(event, dateSpan)
		if err != nil {
			return nil, err
		}
		if len(eventTimespans) == 0 {
			continue
		}

		timespans := []timespan.Span{
			timespan.New(
				utils.MaxTime(dateSpan.From(), span.From()),
				utils.MinTime(dateSpan.To(), span.To()),
			),
		}
		// Exclude event intervals from inactive intervals.
		for _, eventTimespan := range eventTimespans {
			timespansWithoutEvent := make([]timespan.Span, 0)
			for _, activeTimespan := range timespans {
				timespansWithoutEvent = append(timespansWithoutEvent, activeTimespan.Diff(eventTimespan)...)
			}

			timespans = timespansWithoutEvent
		}

		for _, activeTimespan := range timespans {
			inactiveType := c.typesByID[c.defaultTypes[TypeInactive]]
			computed = append(computed, ComputedType{
				Span:  activeTimespan,
				ID:    inactiveType.ID,
				Color: inactiveType.Color,
			})
		}
	}

	return computed, nil
}

func (c *eventComputer) sortExdates(exdates []Exdate) error {
	var err error
	getPriorityByID := func(id string) int {
		t, ok := c.typesByID[id]
		if !ok {
			err = fmt.Errorf("unknown type %v", id)
			return 0
		}

		return t.Priority
	}
	sort.Slice(exdates, func(i, j int) bool {
		return getPriorityByID(exdates[i].Type) > getPriorityByID(exdates[j].Type)
	})

	return err
}
