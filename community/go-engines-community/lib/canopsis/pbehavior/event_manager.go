package pbehavior

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EventManager interface {
	GetEvent(ResolveResult, types.Entity, datetime.CpsTime) (types.Event, error)
}

type eventManager struct {
	eventGenerator libevent.Generator
}

func (r *eventManager) GetEvent(resolveResult ResolveResult, entity types.Entity, ts datetime.CpsTime) (types.Event, error) {
	pbhInfo := NewPBehaviorInfo(ts, resolveResult)
	eventType, output := r.getEventType(pbhInfo, entity.PbehaviorInfo)
	if eventType == "" {
		return types.Event{}, nil
	}

	event, err := r.eventGenerator.Generate(entity)
	if err != nil {
		return event, err
	}

	event.Timestamp = ts
	event.EventType = eventType
	event.Output = output
	event.PbehaviorInfo = pbhInfo
	event.Author = pbhInfo.Author
	event.Initiator = types.InitiatorSystem

	return event, nil
}

func (r *eventManager) getEventType(newPbehaviorInfo, curPbehaviorInfo types.PbehaviorInfo) (string, string) {
	if newPbehaviorInfo.TypeID != "" && newPbehaviorInfo.TypeID == curPbehaviorInfo.TypeID && newPbehaviorInfo.ID == curPbehaviorInfo.ID ||
		newPbehaviorInfo.TypeID == "" && curPbehaviorInfo.IsDefaultActive() {
		return "", ""
	}

	var eventType string
	var output string
	if newPbehaviorInfo.TypeID == "" {
		eventType = types.EventTypePbhLeave
		output = curPbehaviorInfo.GetStepMessage()
	} else {
		if curPbehaviorInfo.IsDefaultActive() {
			eventType = types.EventTypePbhEnter
		} else {
			eventType = types.EventTypePbhLeaveAndEnter
		}

		output = newPbehaviorInfo.GetStepMessage()
	}

	return eventType, output
}

func NewEventManager(eventGenerator libevent.Generator) EventManager {
	return &eventManager{
		eventGenerator: eventGenerator,
	}
}

func NewPBehaviorInfo(time datetime.CpsTime, result ResolveResult) types.PbehaviorInfo {
	if result.Type.ID == "" {
		return types.PbehaviorInfo{}
	}

	pbhInfo := types.PbehaviorInfo{
		Timestamp:     &time,
		ID:            result.ID,
		Name:          result.Name,
		ReasonName:    result.ReasonName,
		ReasonID:      result.ReasonID,
		TypeID:        result.Type.ID,
		TypeName:      result.Type.Name,
		CanonicalType: result.Type.Type,
		Author:        canopsis.DefaultEventAuthor,
	}

	return pbhInfo
}
