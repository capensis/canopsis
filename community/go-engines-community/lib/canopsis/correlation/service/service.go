package service

//go:generate mockgen -destination=../../../../mocks/lib/canopsis/correlation/service/service.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation/service MetaAlarmService

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
)

// MetaAlarmService ...
type MetaAlarmService interface {
	CreateMetaAlarm(
		event types.Event,
		children []types.AlarmWithEntity,
		rule correlation.Rule,
	) (types.Event, error)
	AddChildToMetaAlarm(
		ctx context.Context,
		event types.Event,
		metaAlarm types.Alarm,
		childAlarm types.AlarmWithEntity,
		rule correlation.Rule,
	) (types.Event, error)
	AddMultipleChildsToMetaAlarm(
		ctx context.Context,
		event types.Event,
		metaAlarm types.Alarm,
		children []types.AlarmWithEntity,
		rule correlation.Rule,
	) (types.Event, error)
	RemoveMultipleChildToMetaAlarm(
		ctx context.Context,
		event types.Event,
		metaAlarm types.Alarm,
		children []types.AlarmWithEntity,
		rule correlation.Rule,
	) (types.Event, error)
}

type service struct {
	alarmAdapter        alarm.Adapter
	logger              zerolog.Logger
	alarmConfigProvider config.AlarmConfigProvider
}

type EventExtraInfosMeta struct {
	Rule     correlation.Rule
	Count    int64
	Children types.AlarmWithEntity
}

// NewMetaAlarmService instantiates meta-alarm service; receives alarmAdapter as adapter to db Alarm collection
func NewMetaAlarmService(alarmAdapter alarm.Adapter, alarmConfigProvider config.AlarmConfigProvider, logger zerolog.Logger) MetaAlarmService {
	return &service{
		alarmAdapter:        alarmAdapter,
		alarmConfigProvider: alarmConfigProvider,
		logger:              logger,
	}
}

func (s *service) CreateMetaAlarm(
	event types.Event,
	children []types.AlarmWithEntity,
	rule correlation.Rule,
) (types.Event, error) {
	if len(children) == 0 {
		return types.Event{}, correlation.ErrNoChildren
	}

	infos := EventExtraInfosMeta{
		Rule:     rule,
		Count:    int64(len(children)),
		Children: children[len(children)-1],
	}
	output, err := s.executeOutputTpl(infos)
	if err != nil {
		return types.Event{}, err
	}

	eventChildren := make([]string, 0, len(children))
	for i := 0; i < len(children); i++ {
		eventChildren = append(eventChildren, children[i].Entity.ID)
	}

	return types.Event{
		Timestamp:         event.Timestamp,
		Author:            event.Author,
		State:             event.State,
		Component:         "metaalarm",
		Connector:         "engine",
		ConnectorName:     "correlation",
		Resource:          correlation.DefaultMetaAlarmEntityPrefix + utils.NewID(),
		SourceType:        types.SourceTypeResource,
		EventType:         types.EventTypeMetaAlarm,
		MetaAlarmChildren: &eventChildren,
		MetaAlarmRuleID:   rule.ID,
		Output:            output,
		ExtraInfos: map[string]interface{}{
			"Meta": infos,
		},
	}, nil
}

// AddChildToMetaAlarm makes references from meta-alarm to child and from child to parent,
// updates mata-alarm's state to worst from children
func (s *service) AddChildToMetaAlarm(
	ctx context.Context,
	event types.Event,
	metaAlarm types.Alarm,
	child types.AlarmWithEntity,
	rule correlation.Rule,
) (types.Event, error) {
	if metaAlarm.HasChildByEID(child.Entity.ID) {
		return types.Event{}, correlation.ErrChildAlreadyExist
	}

	childAlarm := child.Alarm

	metaAlarm.AddChild(childAlarm.EntityID)
	childAlarm.AddParent(metaAlarm.EntityID)

	childrenCount, err := s.alarmAdapter.GetCountOpenedAlarmsByIDs(metaAlarm.Value.Children)
	if err != nil {
		return types.Event{}, err
	}

	infos := EventExtraInfosMeta{
		Rule:     rule,
		Count:    childrenCount,
		Children: child,
	}
	output, err := s.executeOutputTpl(infos)
	if err != nil {
		return types.Event{}, err
	}

	metaAlarm.UpdateOutput(output)

	if childAlarm.Value.State != nil {
		maCurrentState := metaAlarm.CurrentState()
		if childAlarm.Value.State.Value > maCurrentState {
			err = metaAlarm.PartialUpdateState(childAlarm.Value.LastUpdateDate, childAlarm.Value.State.Value, metaAlarm.Value.Output, s.alarmConfigProvider.Get())
			if err != nil {
				return types.Event{}, err
			}
		}
	}
	maActions, ticket := metaAlarm.GetAppliedActions()
	if err := childAlarm.ApplyActions(maActions, ticket); err != nil {
		s.logger.Warn().Err(err).Str("alarm-ID", childAlarm.ID).Msg("adding child to metaalarm")
		return types.Event{}, err
	}

	if err := childAlarm.PartialUpdateAddStepWithStep(types.NewMetaAlarmAttachStep(metaAlarm, rule.Name)); err != nil {
		s.logger.Err(err).Str("metaalarm", metaAlarm.EntityID).
			Str("child", childAlarm.EntityID).
			Msg("Failed to add metaalarmattach step to child")
		return types.Event{}, err
	}

	updatedAlarms := []types.Alarm{childAlarm, metaAlarm}
	err = s.alarmAdapter.PartialMassUpdateOpen(ctx, updatedAlarms)
	if err != nil {
		return types.Event{}, err
	}

	return genUpdatedMetaAlarmEvent(event, metaAlarm, infos, output), nil
}

func (s *service) AddMultipleChildsToMetaAlarm(
	ctx context.Context,
	event types.Event,
	metaAlarm types.Alarm,
	children []types.AlarmWithEntity,
	rule correlation.Rule,
) (types.Event, error) {
	worstState, worstStateDate := types.CpsNumber(types.AlarmStateOK), metaAlarm.Value.LastUpdateDate
	maActions, ticket := metaAlarm.GetAppliedActions()

	for i := 0; i < len(children); i++ {
		childAlarm := children[i].Alarm

		if metaAlarm.HasChildByEID(children[i].Entity.ID) {
			return types.Event{}, correlation.ErrChildAlreadyExist
		}

		metaAlarm.AddChild(childAlarm.EntityID)
		childAlarm.AddParent(metaAlarm.EntityID)

		if err := childAlarm.PartialUpdateAddStepWithStep(types.NewMetaAlarmAttachStep(metaAlarm, rule.Name)); err != nil {
			s.logger.Err(err).Str("metaalarm", metaAlarm.EntityID).
				Str("child", childAlarm.EntityID).
				Msg("Failed to add metaalarmattach step to child")
			return types.Event{}, err
		}

		if childAlarm.Value.State != nil && childAlarm.Value.State.Value > worstState {
			worstState, worstStateDate = childAlarm.Value.State.Value, childAlarm.Value.LastUpdateDate
		}

		if err := childAlarm.ApplyActions(maActions, ticket); err != nil {
			s.logger.Warn().Err(err).Str("alarm-ID", childAlarm.ID).Msg("adding children to metaalarm")
			return types.Event{}, err
		}

		children[i].Alarm = childAlarm
	}

	childrenCount, err := s.alarmAdapter.GetCountOpenedAlarmsByIDs(metaAlarm.Value.Children)
	if err != nil {
		return types.Event{}, err
	}
	infos := EventExtraInfosMeta{
		Rule:     rule,
		Count:    childrenCount,
		Children: children[len(children)-1],
	}
	output, err := s.executeOutputTpl(infos)
	if err != nil {
		return types.Event{}, err
	}
	metaAlarm.UpdateOutput(output)

	maCurrentState := metaAlarm.CurrentState()
	if worstState > maCurrentState {
		err = metaAlarm.PartialUpdateState(worstStateDate, worstState, metaAlarm.Value.Output, s.alarmConfigProvider.Get())
		if err != nil {
			return types.Event{}, err
		}
	}

	updated := make([]types.Alarm, len(children)+1)
	for i := 0; i < len(children); i++ {
		updated[i] = children[i].Alarm
	}

	updated[len(children)] = metaAlarm
	err = s.alarmAdapter.PartialMassUpdateOpen(ctx, updated)
	if err != nil {
		return types.Event{}, err
	}

	return genUpdatedMetaAlarmEvent(event, metaAlarm, infos, output), nil
}

func (s *service) RemoveMultipleChildToMetaAlarm(
	ctx context.Context,
	event types.Event,
	metaAlarm types.Alarm,
	children []types.AlarmWithEntity,
	rule correlation.Rule,
) (types.Event, error) {
	for _, child := range children {
		metaAlarm.RemoveChild(child.Alarm.EntityID)
		child.Alarm.RemoveParent(metaAlarm.EntityID)
	}

	metaAlarmChildren := make([]types.AlarmWithEntity, 0)
	err := s.alarmAdapter.GetOpenedAlarmsWithEntityByIDs(metaAlarm.Value.Children, &metaAlarmChildren)
	if err != nil {
		return types.Event{}, err
	}
	infos := EventExtraInfosMeta{
		Rule:     rule,
		Count:    int64(len(metaAlarmChildren)),
		Children: metaAlarmChildren[len(metaAlarmChildren)-1],
	}
	output, err := s.executeOutputTpl(infos)
	if err != nil {
		return types.Event{}, err
	}
	metaAlarm.UpdateOutput(output)

	updated := make([]types.Alarm, len(children)+1)
	for i := 0; i < len(children); i++ {
		updated[i] = children[i].Alarm
	}

	updated[len(children)] = metaAlarm
	err = s.alarmAdapter.PartialMassUpdateOpen(ctx, updated)
	if err != nil {
		return types.Event{}, err
	}

	return genUpdatedMetaAlarmEvent(event, metaAlarm, infos, output), nil
}

func genUpdatedMetaAlarmEvent(
	event types.Event,
	metaAlarm types.Alarm,
	infos EventExtraInfosMeta,
	output string,
) types.Event {
	metaAlarmEvent := types.Event{
		Timestamp:         event.Timestamp,
		Author:            event.Author,
		Component:         metaAlarm.Value.Component,
		Connector:         metaAlarm.Value.Connector,
		ConnectorName:     metaAlarm.Value.ConnectorName,
		Resource:          metaAlarm.Value.Resource,
		EventType:         types.EventTypeMetaAlarmUpdated,
		MetaAlarmRuleID:   infos.Rule.ID,
		MetaAlarmChildren: &metaAlarm.Value.Children,
		SourceType:        types.SourceTypeResource,
		Output:            output,
		ExtraInfos: map[string]interface{}{
			"Meta": infos,
		},
	}

	metaAlarmEvent.SourceType = metaAlarmEvent.DetectSourceType()

	return metaAlarmEvent
}

func (s *service) executeOutputTpl(
	data EventExtraInfosMeta,
) (string, error) {
	rule := data.Rule
	if rule.OutputTemplate == "" {
		return "", nil
	}

	tpl, err := template.New("template").Funcs(types.GetTemplateFunc()).Parse(rule.OutputTemplate)
	if err != nil {
		return "", fmt.Errorf("unable to parse output template for metaalarm rule %s: %+v", rule.ID, err)
	}

	b := strings.Builder{}
	err = tpl.Execute(&b, data)
	if err != nil {
		return "", fmt.Errorf(
			"unable to execute output template for metaalarm rule %s: %+v", rule.ID, err)
	}

	return b.String(), nil
}
