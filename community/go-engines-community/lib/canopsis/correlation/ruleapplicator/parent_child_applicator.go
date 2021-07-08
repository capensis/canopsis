package ruleapplicator

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"github.com/rs/zerolog"
)

type ParentChildApplicator struct {
	alarmAdapter     alarm.Adapter
	metaAlarmService service.MetaAlarmService
	logger           zerolog.Logger
}

func (a ParentChildApplicator) Apply(ctx context.Context, event types.Event, rule correlation.Rule) ([]types.Event, error) {
	var metaAlarmEvent types.Event

	if event.SourceType == types.SourceTypeComponent {
		//skip is component alarm is already a meta-alarm
		if event.Alarm.IsMetaAlarm() {
			return nil, nil
		}

		resourceAlarms, err := a.alarmAdapter.GetAllOpenedResourceAlarmsByComponent(event.Component)
		if err != nil {
			return nil, err
		}

		if len(resourceAlarms) == 0 {
			return nil, nil
		}

		componentAlarm := *event.Alarm
		componentAlarm.SetMeta(rule.ID)

		metaAlarmEvent, err = a.metaAlarmService.AddMultipleChildsToMetaAlarm(ctx, event, componentAlarm, resourceAlarms, rule)
		if err != nil {
			return nil, err
		}
	}

	if event.SourceType == types.SourceTypeResource {
		// Check if component alarm exists and if it's already a meta-alarm
		componentAlarm, err := a.alarmAdapter.GetLastAlarm(event.Connector, event.ConnectorName, event.Component)
		if err != nil {
			if _, ok := err.(errt.NotFound); !ok {
				return nil, err
			}

			return nil, nil
		}

		if !componentAlarm.IsMalfunctioning() {
			return nil, nil
		}

		if !componentAlarm.IsMetaAlarm() {
			// transform to meta-alarm
			componentAlarm.SetMeta(rule.ID)
		}

		metaAlarmEvent, err = a.metaAlarmService.AddChildToMetaAlarm(
			ctx,
			event,
			componentAlarm,
			types.AlarmWithEntity{Alarm: *event.Alarm, Entity: *event.Entity},
			rule,
		)
		if err != nil {
			return nil, err
		}
	}

	if metaAlarmEvent.EventType != "" {
		return []types.Event{metaAlarmEvent}, nil
	}

	return nil, nil
}

func NewParentChildApplicator(alarmAdapter alarm.Adapter, metaAlarmService service.MetaAlarmService, logger zerolog.Logger) ParentChildApplicator {
	return ParentChildApplicator{
		alarmAdapter:     alarmAdapter,
		metaAlarmService: metaAlarmService,
		logger:           logger,
	}
}
