package types

import (
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

// AlarmChangeType is a type representing a change that can occur on an alarm.
type AlarmChangeType string

// An AlarmChangeType can have the following values:
const (
	AlarmChangeTypeNone              AlarmChangeType = ""
	AlarmChangeTypeStateIncrease     AlarmChangeType = "stateinc"
	AlarmChangeTypeStateDecrease     AlarmChangeType = "statedec"
	AlarmChangeTypeCreate            AlarmChangeType = "create"
	AlarmChangeTypeCreateAndPbhEnter AlarmChangeType = "createandpbhenter"
	AlarmChangeTypeAck               AlarmChangeType = "ack"
	AlarmChangeTypeDoubleAck         AlarmChangeType = "doubleack"
	AlarmChangeTypeAckremove         AlarmChangeType = "ackremove"
	AlarmChangeTypeCancel            AlarmChangeType = "cancel"
	AlarmChangeTypeUncancel          AlarmChangeType = "uncancel"
	AlarmChangeTypeAssocTicket       AlarmChangeType = "assocticket"
	AlarmChangeTypeSnooze            AlarmChangeType = "snooze"
	AlarmChangeTypeUnsnooze          AlarmChangeType = "unsnooze"
	AlarmChangeTypeResolve           AlarmChangeType = "resolve"
	AlarmChangeTypeComment           AlarmChangeType = "comment"
	AlarmChangeTypeChangeState       AlarmChangeType = "changestate"
	AlarmChangeTypePbhEnter          AlarmChangeType = "pbhenter"
	AlarmChangeTypePbhLeave          AlarmChangeType = "pbhleave"
	AlarmChangeTypePbhLeaveAndEnter  AlarmChangeType = "pbhleaveandenter"
	AlarmChangeTypeUpdateStatus      AlarmChangeType = "changestatus"
	AlarmChangeTypeActivate          AlarmChangeType = "activate"

	AlarmChangeTypeWebhookStart                 AlarmChangeType = "webhookstart"
	AlarmChangeTypeWebhookComplete              AlarmChangeType = "webhookcomplete"
	AlarmChangeTypeWebhookFail                  AlarmChangeType = "webhookfail"
	AlarmChangeTypeDeclareTicketWebhook         AlarmChangeType = "declareticketwebhook"
	AlarmChangeTypeDeclareTicketWebhookFail     AlarmChangeType = "declareticketwebhookfail"
	AlarmChangeTypeAutoWebhookStart             AlarmChangeType = "autowebhookstart"
	AlarmChangeTypeAutoWebhookComplete          AlarmChangeType = "autowebhookcomplete"
	AlarmChangeTypeAutoWebhookFail              AlarmChangeType = "autowebhookfail"
	AlarmChangeTypeAutoDeclareTicketWebhook     AlarmChangeType = "autodeclareticketwebhook"
	AlarmChangeTypeAutoDeclareTicketWebhookFail AlarmChangeType = "autodeclareticketwebhookfail"

	// Following change types are used for manual instruction execution.
	AlarmChangeTypeInstructionStart    AlarmChangeType = "instructionstart"
	AlarmChangeTypeInstructionPause    AlarmChangeType = "instructionpause"
	AlarmChangeTypeInstructionResume   AlarmChangeType = "instructionresume"
	AlarmChangeTypeInstructionComplete AlarmChangeType = "instructioncomplete"
	AlarmChangeTypeInstructionFail     AlarmChangeType = "instructionfail"
	// AlarmChangeTypeInstructionAbort is used for manual and auto instruction execution.
	AlarmChangeTypeInstructionAbort AlarmChangeType = "instructionabort"
	// Following change types are used for auto instruction execution.
	AlarmChangeTypeAutoInstructionStart    AlarmChangeType = "autoinstructionstart"
	AlarmChangeTypeAutoInstructionComplete AlarmChangeType = "autoinstructioncomplete"
	AlarmChangeTypeAutoInstructionFail     AlarmChangeType = "autoinstructionfail"
	// Following change types are used for auto instruction triggers.
	AlarmChangeTypeAutoInstructionResultOk   AlarmChangeType = "autoinstructionresultok"
	AlarmChangeTypeAutoInstructionResultFail AlarmChangeType = "autoinstructionresultfail"
	// Following change types are used for job execution.
	AlarmChangeTypeInstructionJobStart    AlarmChangeType = "instructionjobstart"
	AlarmChangeTypeInstructionJobComplete AlarmChangeType = "instructionjobcomplete"
	AlarmChangeTypeInstructionJobFail     AlarmChangeType = "instructionjobfail"

	// Following change types are used for junit.
	AlarmChangeTypeJunitTestSuiteUpdate AlarmChangeType = "junittestsuiteupdate"
	AlarmChangeTypeJunitTestCaseUpdate  AlarmChangeType = "junittestcaseupdate"

	AlarmChangeTypeEnabled AlarmChangeType = "enabled"

	// AlarmChangeTypeAutoInstructionActivate is used to activate alarm when an autoremediation triggered by create trigger is completed
	AlarmChangeTypeAutoInstructionActivate AlarmChangeType = "autoinstructionactivate"

	// AlarmChangeEventsCount is used for eventscount trigger and alarm's events_count value should be added to the end of a trigger name, e.g. "eventscount3"
	AlarmChangeEventsCount AlarmChangeType = "eventscount"
)

const MinimalEventsCountThreshold = 2

// AlarmChange is a struct containing the type of change that occurred on an
// alarm, as well as its previous state.
type AlarmChange struct {
	Type                            AlarmChangeType   `json:"Type"`
	PreviousState                   CpsNumber         `json:"PreviousState"`
	PreviousStateChange             datetime.CpsTime  `json:"PreviousStateChange"`
	PreviousStatus                  CpsNumber         `json:"PreviousStatus"`
	PreviousStatusChange            datetime.CpsTime  `json:"PreviousStatusChange"`
	PreviousPbehaviorTime           *datetime.CpsTime `json:"PreviousPbehaviorTime"`
	PreviousEntityPbehaviorTime     *datetime.CpsTime `json:"PreviousEntityPbehaviorTime"`
	PreviousPbehaviorTypeID         string            `json:"PreviousPbehaviorTypeID"`
	PreviousPbehaviorCannonicalType string            `json:"PreviousPbehaviorCannonicalType"`

	EventsCount int `json:"EventsCount"`
}

func NewAlarmChange() AlarmChange {
	return AlarmChange{
		Type:                 AlarmChangeTypeNone,
		PreviousState:        AlarmStateOK,
		PreviousStateChange:  datetime.NewCpsTime(),
		PreviousStatus:       AlarmStatusOff,
		PreviousStatusChange: datetime.NewCpsTime(),
	}
}

func NewAlarmChangeByAlarm(alarm Alarm, t ...AlarmChangeType) AlarmChange {
	alarmChangeType := AlarmChangeTypeNone
	if len(t) > 0 {
		alarmChangeType = t[0]
	}

	return AlarmChange{
		Type:                            alarmChangeType,
		PreviousState:                   alarm.Value.State.Value,
		PreviousStateChange:             alarm.Value.State.Timestamp,
		PreviousStatus:                  alarm.Value.Status.Value,
		PreviousStatusChange:            alarm.Value.Status.Timestamp,
		PreviousPbehaviorTypeID:         alarm.Value.PbehaviorInfo.TypeID,
		PreviousPbehaviorCannonicalType: alarm.Value.PbehaviorInfo.CanonicalType,
	}
}

func (ac *AlarmChange) GetTriggers() []string {
	var triggers []string

	switch ac.Type {
	case AlarmChangeTypeNone:
		if ac.EventsCount >= MinimalEventsCountThreshold {
			triggers = append(triggers, string(AlarmChangeEventsCount)+strconv.Itoa(ac.EventsCount))
		}
	case AlarmChangeTypeStateIncrease:
		triggers = append(triggers, string(AlarmChangeTypeStateIncrease))
		if ac.EventsCount >= MinimalEventsCountThreshold {
			triggers = append(triggers, string(AlarmChangeEventsCount)+strconv.Itoa(ac.EventsCount))
		}
	case AlarmChangeTypeStateDecrease:
		triggers = append(triggers, string(AlarmChangeTypeStateDecrease))
		if ac.EventsCount >= MinimalEventsCountThreshold {
			triggers = append(triggers, string(AlarmChangeEventsCount)+strconv.Itoa(ac.EventsCount))
		}
	case AlarmChangeTypeCreateAndPbhEnter:
		triggers = append(triggers, string(AlarmChangeTypeCreate), string(AlarmChangeTypePbhEnter))
	case AlarmChangeTypePbhLeaveAndEnter:
		triggers = append(triggers, string(AlarmChangeTypePbhEnter), string(AlarmChangeTypePbhLeave))
	case AlarmChangeTypeDoubleAck:
		triggers = append(triggers, string(AlarmChangeTypeAck))
	case AlarmChangeTypeWebhookStart,
		AlarmChangeTypeWebhookComplete,
		AlarmChangeTypeWebhookFail,
		AlarmChangeTypeDeclareTicketWebhookFail,
		AlarmChangeTypeAutoWebhookStart,
		AlarmChangeTypeAutoWebhookComplete,
		AlarmChangeTypeAutoWebhookFail,
		AlarmChangeTypeAutoDeclareTicketWebhookFail,
		AlarmChangeTypeAutoInstructionActivate:
		// not a trigger
	case AlarmChangeTypeDeclareTicketWebhook,
		AlarmChangeTypeAutoDeclareTicketWebhook:
		triggers = append(triggers, string(AlarmChangeTypeDeclareTicketWebhook))
	default:
		trigger := string(ac.Type)
		if trigger != "" {
			triggers = append(triggers, trigger)
		}
	}

	return triggers
}

func (ac *AlarmChange) IsZero() bool {
	if ac == nil {
		return true
	}

	return *ac == AlarmChange{}
}
