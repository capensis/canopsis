package pbehavior

import (
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EventManager interface {
	GetEvent(ResolveResult, types.Alarm, types.Entity, time.Time) types.Event
	GetEventType(resolveResult ResolveResult, curPbehaviorInfo types.PbehaviorInfo) (eventType string, output string)
}

type eventManager struct {
	connector string
}

func (r eventManager) GetEvent(resolveResult ResolveResult, alarm types.Alarm, entity types.Entity, now time.Time) types.Event {
	eventType, output := r.GetEventType(resolveResult, alarm.Value.PbehaviorInfo)
	if eventType == "" {
		return types.Event{}
	}

	event := types.Event{
		Connector:     r.connector,
		ConnectorName: r.connector,
		Component:     alarm.Value.Component,
		Resource:      alarm.Value.Resource,
		SourceType:    entity.Type,
		Timestamp:     datetime.CpsTime{Time: now},
		EventType:     eventType,
		Output:        output,
		PbehaviorInfo: NewPBehaviorInfo(datetime.CpsTime{Time: now}, resolveResult),
		Author:        canopsis.DefaultEventAuthor,
		Initiator:     types.InitiatorSystem,
	}

	return event
}

func (r eventManager) GetEventType(resolveResult ResolveResult, curPbehaviorInfo types.PbehaviorInfo) (string, string) {
	resolvedType := resolveResult.ResolvedType

	if resolvedType.ID != "" && resolvedType.ID == curPbehaviorInfo.TypeID && resolveResult.ResolvedPbhID == curPbehaviorInfo.ID ||
		resolvedType.ID == "" && curPbehaviorInfo.IsDefaultActive() {
		return "", ""
	}

	var eventType string
	var output string
	if resolvedType.ID == "" {
		eventType = types.EventTypePbhLeave
		output = fmt.Sprintf(
			"Pbehavior %s. Type: %s. Reason: %s.",
			curPbehaviorInfo.Name,
			curPbehaviorInfo.TypeName,
			curPbehaviorInfo.ReasonName,
		)
	} else {
		if curPbehaviorInfo.IsDefaultActive() {
			eventType = types.EventTypePbhEnter
		} else {
			eventType = types.EventTypePbhLeaveAndEnter
		}

		output = fmt.Sprintf(
			"Pbehavior %s. Type: %s. Reason: %s.",
			resolveResult.ResolvedPbhName,
			resolvedType.Name,
			resolveResult.ResolvedPbhReasonName,
		)
	}

	return eventType, output
}

func NewEventManager(connector string) EventManager {
	return eventManager{
		connector: connector,
	}
}

func NewPBehaviorInfo(time datetime.CpsTime, result ResolveResult) types.PbehaviorInfo {
	if result.ResolvedType.ID == "" {
		return types.PbehaviorInfo{}
	}

	return types.PbehaviorInfo{
		Timestamp:     &time,
		ID:            result.ResolvedPbhID,
		Name:          result.ResolvedPbhName,
		ReasonName:    result.ResolvedPbhReasonName,
		ReasonID:      result.ResolvedPbhReasonID,
		TypeID:        result.ResolvedType.ID,
		TypeName:      result.ResolvedType.Name,
		CanonicalType: result.ResolvedType.Type,
	}
}
