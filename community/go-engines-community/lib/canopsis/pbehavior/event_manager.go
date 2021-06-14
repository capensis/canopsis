package pbehavior

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"time"
)

type EventManager interface {
	GetEvent(ResolveResult, types.Alarm, time.Time) types.Event
}

type eventManager struct {
}

func (r eventManager) GetEvent(resolveResult ResolveResult, alarm types.Alarm, now time.Time) types.Event {
	resolvedType := resolveResult.ResolvedType
	curPbehaviorInfo := alarm.Value.PbehaviorInfo

	if resolvedType != nil && resolvedType.ID == curPbehaviorInfo.TypeID && resolveResult.ResolvedPbhID == curPbehaviorInfo.ID ||
		resolvedType == nil && curPbehaviorInfo.IsDefaultActive() {
		return types.Event{}
	}

	var eventType string
	var output string
	if resolvedType == nil {
		eventType = types.EventTypePbhLeave
		output = fmt.Sprintf(
			"Pbehavior %s. Type: %s. Reason: %s",
			curPbehaviorInfo.Name,
			curPbehaviorInfo.TypeName,
			curPbehaviorInfo.Reason,
		)
	} else {
		if curPbehaviorInfo.IsDefaultActive() {
			eventType = types.EventTypePbhEnter
		} else {
			eventType = types.EventTypePbhLeaveAndEnter
		}

		output = fmt.Sprintf(
			"Pbehavior %s. Type: %s. Reason: %s",
			resolveResult.ResolvedPbhName,
			resolvedType.Name,
			resolveResult.ResolvedPbhReason,
		)
	}

	event := types.Event{
		Connector:     alarm.Value.Connector,
		ConnectorName: alarm.Value.ConnectorName,
		Component:     alarm.Value.Component,
		Resource:      alarm.Value.Resource,
		Timestamp:     types.CpsTime{Time: now},
		EventType:     eventType,
		Output:        output,
	}

	event.PbehaviorInfo = NewPBehaviorInfo(resolveResult)
	event.SourceType = event.DetectSourceType()

	return event
}

func NewEventManager() EventManager {
	return eventManager{}
}

func NewPBehaviorInfo(result ResolveResult) types.PbehaviorInfo {
	if result.ResolvedType == nil {
		return types.PbehaviorInfo{}
	}

	return types.PbehaviorInfo{
		ID:            result.ResolvedPbhID,
		Name:          result.ResolvedPbhName,
		Reason:        result.ResolvedPbhReason,
		TypeID:        result.ResolvedType.ID,
		TypeName:      result.ResolvedType.Name,
		CanonicalType: result.ResolvedType.Type,
	}
}
