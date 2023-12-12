package pbehaviorics

import (
	"fmt"
	"math"
	"time"

	pbehaviorapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/ics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/teambition/rrule-go"
)

// Service is used to implement ICS calendar generation.
type Service interface {
	// GenICSFrom creates ICS calendar from pbehavior.
	GenICSFrom(pbh pbehaviorapi.Response, maxPriority, minPriority int64) (ics.Calendar, error)
}

// NewService creates new service.
func NewService(timezoneConfigProvider config.TimezoneConfigProvider) Service {
	return &service{
		timezoneConfigProvider: timezoneConfigProvider,
	}
}

// service represents ICS calendar generation.
type service struct {
	timezoneConfigProvider config.TimezoneConfigProvider
}

func (s *service) GenICSFrom(
	pbh pbehaviorapi.Response,
	maxPriority, minPriority int64,
) (ics.Calendar, error) {
	location := s.timezoneConfigProvider.Get().Location
	e := ics.NewEvent()
	e.SetUID(pbh.ID)
	e.SetCreated(pbh.Created.Time.In(location))
	e.SetLastModified(pbh.Updated.Time.In(location))
	e.SetDTStart(pbh.Start.Time.In(location))
	if pbh.Stop != nil {
		e.SetDTEnd(pbh.Stop.Time.In(location))
	}
	e.SetDTStamp(time.Now().In(location))
	e.SetRRule(pbh.RRule)
	description := fmt.Sprintf("%s\\n%s (%s)", pbh.Reason.Name, pbh.Type.Name, pbh.Type.Type)
	e.SetDescription(description)
	e.SetSummary(pbh.Name)
	var rOption *rrule.ROption
	if pbh.RRule != "" {
		var err error
		rOption, err = rrule.StrToROption(pbh.RRule)
		if err != nil {
			return nil, err
		}
	}

	var event pbehavior.Event
	if pbh.Stop != nil {
		event = pbehavior.NewRecEvent(pbh.Start.Time, pbh.Stop.Time, rOption)
	}
	// Set exdate from exdates
	if len(pbh.Exdates) > 0 {
		for _, exd := range pbh.Exdates {
			if pbh.Stop == nil {
				event = pbehavior.NewRecEvent(pbh.Start.Time, exd.End.Time, rOption)
			}

			timespans, err := pbehavior.GetTimeSpans(event, timespan.New(exd.Begin.Time, exd.End.Time))
			if err != nil {
				return nil, err
			}

			for _, s := range timespans {
				e.SetExdate(s.From().In(location))
			}
		}
	}

	// Set exdate from exceptions
	if len(pbh.Exceptions) > 0 {
		for _, ex := range pbh.Exceptions {
			for _, exd := range ex.Exdates {
				if pbh.Stop == nil {
					event = pbehavior.NewRecEvent(pbh.Start.Time, exd.End.Time, rOption)
				}

				timespans, err := pbehavior.GetTimeSpans(event, timespan.New(exd.Begin.Time, exd.End.Time))
				if err != nil {
					return nil, err
				}

				for _, s := range timespans {
					e.SetExdate(s.From().In(location))
				}
			}
		}
	}

	// Set priority : transform priority to fit into allowed ICS priority interval.
	ratio := float64(ics.MaxPriority-ics.MinPriority) / float64(maxPriority-minPriority)
	priority := (float64(pbh.Type.Priority)-float64(minPriority))*ratio + float64(ics.MinPriority)
	e.SetPriority(int64(math.Ceil(priority)))

	cal := ics.NewCalendar()
	cal.AddComponent(e)

	return cal, nil
}
