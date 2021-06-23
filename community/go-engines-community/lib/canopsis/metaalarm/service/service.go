package service

//go:generate mockgen -destination=../../../../mocks/lib/canopsis/metaalarm/service/service.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/service MetaAlarmService

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"text/template"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
)

// MetaAlarmService ...
type MetaAlarmService interface {
	CreateMetaAlarm(
		event *types.Event,
		children []types.AlarmWithEntity,
		rule metaalarm.Rule,
	) (types.Event, error)
	AddChildToMetaAlarm(
		event *types.Event,
		metaAlarm types.Alarm,
		childAlarm types.AlarmWithEntity,
		rule metaalarm.Rule,
	) (types.Event, error)
	AddMultipleChildsToMetaAlarm(
		event *types.Event,
		metaAlarm types.Alarm,
		children []types.AlarmWithEntity,
		rule metaalarm.Rule,
	) (types.Event, error)
	RemoveMultipleChildToMetaAlarm(
		event *types.Event,
		metaAlarm types.Alarm,
		children []types.AlarmWithEntity,
		rule metaalarm.Rule,
	) (types.Event, error)
}

type service struct {
	alarmAdapter        alarm.Adapter
	logger              zerolog.Logger
	alarmConfigProvider config.AlarmConfigProvider
	ruleAdapter         metaalarm.RulesAdapter
}

const metaAlarmEntityPrefix = "meta-alarm-entity-"

type eventExtraInfosMeta struct {
	Rule     metaalarm.Rule
	Count    int
	Children types.AlarmWithEntity
}

// NewMetaAlarmService instantiates meta-alarm service; receives alarmAdapter as adapter to db Alarm collection
func NewMetaAlarmService(
	alarmAdapter alarm.Adapter, ruleApt metaalarm.RulesAdapter,
	alarmConfigProvider config.AlarmConfigProvider, logger zerolog.Logger) MetaAlarmService {
	return &service{
		alarmAdapter:        alarmAdapter,
		ruleAdapter:         ruleApt,
		alarmConfigProvider: alarmConfigProvider,
		logger:              logger,
	}
}

// CreateMetaAlarm ...
func (s *service) CreateMetaAlarm(
	event *types.Event,
	children []types.AlarmWithEntity,
	rule metaalarm.Rule,
) (types.Event, error) {
	var lastChild types.AlarmWithEntity
	if len(children) > 0 {
		lastChild = children[len(children)-1]
	}
	infos := eventExtraInfosMeta{
		Rule:     rule,
		Count:    len(children),
		Children: lastChild,
	}
	output, err := s.executeOutputTpl(infos)
	if err != nil {
		return types.Event{}, err
	}
	metaAlarmEvent := s.genCreateMetaAlarmEvent(*event, infos)
	metaAlarmEvent.Output = output
	for i := 0; i < len(children); i++ {
		*metaAlarmEvent.MetaAlarmChildren = append(*metaAlarmEvent.MetaAlarmChildren, children[i].Alarm.EntityID)
	}

	return metaAlarmEvent, nil
}

func (s *service) genCreateMetaAlarmEvent(
	event types.Event,
	infos eventExtraInfosMeta,
) types.Event {
	chlds := make([]string, 0)
	resource := metaAlarmEntityPrefix + utils.NewID()
	metaAlarmEvent := types.Event{
		Timestamp:         event.Timestamp,
		Author:            event.Author,
		State:             event.State,
		Component:         "metaalarm",
		Connector:         "engine",
		ConnectorName:     "correlation",
		Resource:          resource,
		SourceType:        types.SourceTypeResource,
		EventType:         types.EventTypeMetaAlarm,
		MetaAlarmChildren: &chlds,
		MetaAlarmRuleID:   infos.Rule.ID,
		ExtraInfos: map[string]interface{}{
			"Meta": infos,
		},
	}
	return metaAlarmEvent
}

func (s *service) genUpdatedMetaAlarmEvent(
	event types.Event,
	metaAlarm types.Alarm,
	infos eventExtraInfosMeta,
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
		ExtraInfos: map[string]interface{}{
			"Meta": infos,
		},
	}
	metaAlarmEvent.SourceType = metaAlarmEvent.DetectSourceType()

	return metaAlarmEvent
}

func (s *service) mkMetaAlarmAttachStep(metaAlarm types.Alarm) types.AlarmStep {
	ruleIdentifier := metaAlarm.Value.Meta
	rule, err := s.ruleAdapter.GetRule(metaAlarm.Value.Meta)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			s.logger.Err(err).Str("rule", metaAlarm.Value.Meta).Msg("Get rule had failed")
		}
	} else {
		ruleIdentifier = rule.Name
	}
	newStep := types.NewMetaAlarmAttachStep(metaAlarm, ruleIdentifier)
	return newStep
}

// AddChildToMetaAlarm makes references from meta-alarm to child and from child to parent,
// updates mata-alarm's state to worst from children
func (s *service) AddChildToMetaAlarm(
	event *types.Event,
	metaAlarm types.Alarm,
	child types.AlarmWithEntity,
	rule metaalarm.Rule,
) (types.Event, error) {
	childAlarm := child.Alarm
	isExistedAlarm := s.isExisted(metaAlarm, childAlarm)
	if !isExistedAlarm {
		metaAlarm.Value.Children = append(metaAlarm.Value.Children, childAlarm.EntityID)
		childAlarm.Value.Parents = append(childAlarm.Value.Parents, metaAlarm.EntityID)
	}
	childrenCount, err := s.alarmAdapter.GetCountOpenedAlarmsByIDs(metaAlarm.Value.Children)
	if err != nil {
		return types.Event{}, err
	}
	infos := eventExtraInfosMeta{
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
			metaAlarm.UpdateState(childAlarm.Value.State.Value, childAlarm.Value.LastUpdateDate)
		} else if isExistedAlarm && childAlarm.Value.State.Value < maCurrentState {
			alarm.UpdateToWorstState(&metaAlarm, []*types.Alarm{&childAlarm}, s.alarmAdapter, s.alarmConfigProvider.Get())
		}
	}
	maActions, ticket := metaAlarm.GetAppliedActions()
	if _, err := childAlarm.ApplyActions(maActions, ticket); err != nil {
		s.logger.Warn().Err(err).Str("alarm-ID", childAlarm.ID).Msg("adding child to metaalarm")
	}

	newStep := s.mkMetaAlarmAttachStep(metaAlarm)
	if err := childAlarm.Value.Steps.Add(newStep); err != nil {
		s.logger.Err(err).Str("metaalarm", metaAlarm.EntityID).
			Str("child", childAlarm.EntityID).
			Msg("Failed to add metaalarmattach step to child")
	}

	updatedAlarms := []types.Alarm{metaAlarm, childAlarm}
	err = s.alarmAdapter.MassUpdate(updatedAlarms, true)
	if err != nil {
		return types.Event{}, err
	}

	metaAlarmEvent := s.genUpdatedMetaAlarmEvent(*event, metaAlarm, infos)
	metaAlarmEvent.Output = output

	return metaAlarmEvent, nil
}

func (s *service) AddMultipleChildsToMetaAlarm(
	event *types.Event,
	metaAlarm types.Alarm,
	children []types.AlarmWithEntity,
	rule metaalarm.Rule,
) (types.Event, error) {
	worstState, worstStateDate := types.CpsNumber(types.AlarmStateOK), metaAlarm.Value.LastUpdateDate
	updateChildren := make([]*types.Alarm, 0, len(children))
	maActions, ticket := metaAlarm.GetAppliedActions()
	for i := 0; i < len(children); i++ {
		childAlarm := children[i].Alarm
		if !s.isExisted(metaAlarm, childAlarm) {
			metaAlarm.Value.Children = append(metaAlarm.Value.Children, childAlarm.EntityID)
			childAlarm.Value.Parents = append(childAlarm.Value.Parents, metaAlarm.EntityID)

			newStep := s.mkMetaAlarmAttachStep(metaAlarm)
			if err := childAlarm.Value.Steps.Add(newStep); err != nil {
				s.logger.Err(err).Str("metaalarm", metaAlarm.EntityID).
					Str("child", childAlarm.EntityID).
					Msg("Failed to add metaalarmattach step to child")
			}
			children[i].Alarm = childAlarm
		} else {
			updateChildren = append(updateChildren, &childAlarm)
		}
		if childAlarm.Value.State != nil && childAlarm.Value.State.Value > worstState {
			worstState, worstStateDate = childAlarm.Value.State.Value, childAlarm.Value.LastUpdateDate
		}
		if done, err := childAlarm.ApplyActions(maActions, ticket); err != nil {
			s.logger.Warn().Err(err).Str("alarm-ID", childAlarm.ID).Msg("adding children to metaalarm")
		} else if done {
			children[i].Alarm = childAlarm
		}
	}

	childrenCount, err := s.alarmAdapter.GetCountOpenedAlarmsByIDs(metaAlarm.Value.Children)
	if err != nil {
		return types.Event{}, err
	}
	infos := eventExtraInfosMeta{
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
		metaAlarm.UpdateState(worstState, worstStateDate)
	} else if worstState < maCurrentState && len(updateChildren) > 0 {
		alarm.UpdateToWorstState(&metaAlarm, updateChildren, s.alarmAdapter, s.alarmConfigProvider.Get())
	}

	updated := make([]types.Alarm, len(children))
	for i := 0; i < len(children); i++ {
		updated[i] = children[i].Alarm
	}

	updated = append(updated, metaAlarm)
	err = s.alarmAdapter.MassUpdate(updated, true)
	if err != nil {
		return types.Event{}, err
	}

	metaAlarmEvent := s.genUpdatedMetaAlarmEvent(*event, metaAlarm, infos)
	metaAlarmEvent.Output = output

	return metaAlarmEvent, nil
}

func (s *service) RemoveMultipleChildToMetaAlarm(
	event *types.Event,
	metaAlarm types.Alarm,
	children []types.AlarmWithEntity,
	rule metaalarm.Rule,
) (types.Event, error) {
	for i := 0; i < len(children); i++ {
		if childrenIndex := s.indexOfChildren(metaAlarm, children[i].Alarm); childrenIndex != -1 {
			metaAlarm.Value.Children = s.removeIndex(metaAlarm.Value.Children, childrenIndex)
		}

		if parentIndex := s.indexOfParents(metaAlarm, children[i].Alarm); parentIndex != -1 {
			children[i].Alarm.Value.Parents = s.removeIndex(children[i].Alarm.Value.Parents, parentIndex)
		}
	}

	metaAlarmChildren := make([]types.AlarmWithEntity, 0)
	err := s.alarmAdapter.GetOpenedAlarmsWithEntityByIDs(metaAlarm.Value.Children, &metaAlarmChildren)
	if err != nil {
		return types.Event{}, err
	}
	infos := eventExtraInfosMeta{
		Rule:     rule,
		Count:    len(metaAlarmChildren),
		Children: metaAlarmChildren[len(children)-1],
	}
	output, err := s.executeOutputTpl(infos)
	if err != nil {
		return types.Event{}, err
	}
	metaAlarm.UpdateOutput(output)

	updated := make([]types.Alarm, len(children))
	for i := 0; i < len(children); i++ {
		updated[i] = children[i].Alarm
	}

	updated = append(updated, metaAlarm)
	err = s.alarmAdapter.MassUpdate(updated, true)
	if err != nil {
		return types.Event{}, err
	}

	metaAlarmEvent := s.genUpdatedMetaAlarmEvent(*event, metaAlarm, infos)
	metaAlarmEvent.Output = output

	return metaAlarmEvent, nil
}

func (s *service) removeIndex(sl []string, idx int) []string {
	sl[idx] = sl[len(sl)-1]
	sl[len(sl)-1] = ""
	sl = sl[:len(sl)-1]
	return sl
}

func (s *service) indexOfParents(metaAlarm types.Alarm, alarm types.Alarm) int {
	for i := 0; i < len(alarm.Value.Parents); i++ {
		if alarm.Value.Parents[i] == metaAlarm.EntityID {
			return i
		}
	}
	return -1
}

func (s *service) indexOfChildren(metaAlarm types.Alarm, alarm types.Alarm) int {
	for i := 0; i < len(metaAlarm.Value.Children); i++ {
		if metaAlarm.Value.Children[i] == alarm.EntityID {
			return i
		}
	}
	return -1
}

func (s *service) isExisted(metaAlarm types.Alarm, alarm types.Alarm) bool {
	for _, childrenID := range metaAlarm.Value.Children {
		if childrenID == alarm.EntityID {
			return true
		}
	}

	return false
}

func (s *service) executeOutputTpl(
	data eventExtraInfosMeta,
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
