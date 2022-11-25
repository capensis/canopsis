package ruleapplicator

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ManualRule = "Manual alarm"

var manualRule *metaalarm.Rule

type ManualGroupApplicator struct {
	alarmAdapter     alarm.Adapter
	ruleAdapter      metaalarm.RulesAdapter
	metaAlarmService service.MetaAlarmService
	logger           zerolog.Logger
}

func NewManualGroupApplicator(alarmAdapter alarm.Adapter, metaAlarmService service.MetaAlarmService,
	ruleAdapter metaalarm.RulesAdapter, logger zerolog.Logger) ManualGroupApplicator {
	return ManualGroupApplicator{
		alarmAdapter:     alarmAdapter,
		ruleAdapter:      ruleAdapter,
		metaAlarmService: metaAlarmService,
		logger:           logger,
	}
}

type missingRequiredFields struct {
	fields []string
}

func (m missingRequiredFields) Error() string {
	return fmt.Sprintf("missing require fields %s", m.fields)
}

func (a ManualGroupApplicator) getOrCreateRule() (metaalarm.Rule, error) {
	if manualRule != nil {
		return *manualRule, nil
	}
	rule, err := a.ruleAdapter.GetManualRule()
	if err == nil {
		manualRule = &rule
		return rule, nil
	}
	rule = metaalarm.Rule{
		ID:       fmt.Sprintf("zgrp-%s", primitive.NewObjectID().Hex()),
		Type:     metaalarm.RuleManualGroup,
		Patterns: metaalarm.RulePatterns{},
		Config:   metaalarm.RuleConfig{},
		Name:     ManualRule,
	}
	err = a.ruleAdapter.Save(rule)
	if err == nil {
		manualRule = &rule
		return *manualRule, nil
	}
	return rule, err

}

func (a ManualGroupApplicator) retrieveListAssociatedAlarms(event *types.Event) ([]types.AlarmWithEntity, *[]types.Alarm, error) {
	if event.MetaAlarmParents == nil || event.MetaAlarmChildren == nil {
		return nil, nil, missingRequiredFields{[]string{"ma_children", "ma_parents"}}
	}

	var children []types.AlarmWithEntity
	err := a.alarmAdapter.GetOpenedAlarmsWithEntityByIDs(*event.MetaAlarmChildren, &children)
	if err != nil {
		a.logger.Error().Interface("alarm_id_list", *event.MetaAlarmChildren).
			Err(err).
			Msg("Failed to retrieve list of children alarms")
		return nil, nil, err
	}

	var metaalarms []types.Alarm
	err = a.alarmAdapter.GetOpenedAlarmsByIDs(*event.MetaAlarmParents, &metaalarms)
	if err != nil {
		a.logger.Error().Interface("meta_alarm_id_list", *event.MetaAlarmParents).
			Err(err).
			Msg("Failed to retrieve list of meta alarms")
		return nil, nil, err
	}

	return children, &metaalarms, nil
}

func (a ManualGroupApplicator) addAlarmsToGroups(
	event *types.Event,
	children []types.AlarmWithEntity,
	metaalarms *[]types.Alarm,
	rule metaalarm.Rule,
) ([]types.Event, error) {
	metaAlarmEvents := make([]types.Event, 0)
	for _, ma := range *metaalarms {
		metaAlarmEvent, err := a.metaAlarmService.AddMultipleChildsToMetaAlarm(event, ma, children, rule)
		if err != nil {
			return nil, err
		}

		metaAlarmEvents = append(metaAlarmEvents, metaAlarmEvent)
	}
	return metaAlarmEvents, nil
}

func (a ManualGroupApplicator) removeAlarmsToGroups(
	event *types.Event,
	children []types.AlarmWithEntity,
	metaalarms *[]types.Alarm,
	rule metaalarm.Rule,
) ([]types.Event, error) {
	metaAlarmEvents := make([]types.Event, 0)
	for _, ma := range *metaalarms {
		metaAlarmEvent, err := a.metaAlarmService.RemoveMultipleChildToMetaAlarm(event, ma, children, rule)
		if err != nil {
			return nil, err
		}

		metaAlarmEvents = append(metaAlarmEvents, metaAlarmEvent)
	}

	return metaAlarmEvents, nil
}

func (a ManualGroupApplicator) groupAlarms(event *types.Event) (types.Event, error) {
	var metaAlarmEvent types.Event
	var err error

	if event.MetaAlarmChildren == nil {
		return metaAlarmEvent, missingRequiredFields{[]string{"ma_children"}}
	}

	rule, err := a.getOrCreateRule()
	if err != nil {
		return metaAlarmEvent, err
	}

	var children []types.AlarmWithEntity
	err = a.alarmAdapter.GetOpenedAlarmsWithEntityByIDs(*event.MetaAlarmChildren, &children)
	if err != nil {
		return metaAlarmEvent, err
	}

	metaAlarmEvent, err = a.metaAlarmService.CreateMetaAlarm(event, children, rule)
	if err != nil {
		return metaAlarmEvent, err
	}

	return metaAlarmEvent, nil
}

func (a ManualGroupApplicator) Apply(ctx context.Context, event *types.Event, r metaalarm.Rule) ([]types.Event, error) {

	if event.EventType == types.EventManualMetaAlarmGroup {
		metaAlarmEvent, err := a.groupAlarms(event)
		if err == nil {
			metaAlarmEvent.ExtraInfos = event.ExtraInfos
			return []types.Event{metaAlarmEvent}, nil
		}

		return nil, err
	}

	var metaAlarmEvents []types.Event
	var err error
	children, metaAlarms, err := a.retrieveListAssociatedAlarms(event)
	if err != nil {
		return nil, err
	}

	if event.EventType == types.EventManualMetaAlarmUngroup {
		metaAlarmEvents, err = a.removeAlarmsToGroups(event, children, metaAlarms, r)
	} else if event.EventType == types.EventManualMetaAlarmUpdate {
		metaAlarmEvents, err = a.addAlarmsToGroups(event, children, metaAlarms, r)
	}

	return metaAlarmEvents, err
}
