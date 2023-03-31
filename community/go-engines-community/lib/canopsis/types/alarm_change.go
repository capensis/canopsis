package types

import (
	"time"
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
	// Following change types are used for job execution.
	AlarmChangeTypeInstructionJobStart    AlarmChangeType = "instructionjobstart"
	AlarmChangeTypeInstructionJobComplete AlarmChangeType = "instructionjobcomplete"
	AlarmChangeTypeInstructionJobFail     AlarmChangeType = "instructionjobfail"

	// Following change types are used for junit.
	AlarmChangeTypeJunitTestSuiteUpdate AlarmChangeType = "junittestsuiteupdate"
	AlarmChangeTypeJunitTestCaseUpdate  AlarmChangeType = "junittestcaseupdate"

	// AlarmChangeTypeEntityToggled is used to update entity service's counters on disable/enable entity actions.
	AlarmChangeTypeEntityToggled AlarmChangeType = "entitytoggled"
)

// AlarmChange is a struct containing the type of change that occured on an
// alarm, as well as its previous state.
type AlarmChange struct {
	Type                            AlarmChangeType
	PreviousState                   CpsNumber
	PreviousStateChange             CpsTime
	PreviousStatus                  CpsNumber
	PreviousStatusChange            CpsTime
	PreviousPbehaviorTypeID         string
	PreviousPbehaviorCannonicalType string
}

func NewAlarmChange() AlarmChange {
	return AlarmChange{
		Type:                 AlarmChangeTypeNone,
		PreviousState:        AlarmStateOK,
		PreviousStateChange:  CpsTime{Time: time.Now()},
		PreviousStatus:       AlarmStatusOff,
		PreviousStatusChange: CpsTime{Time: time.Now()},
	}
}

func (ac *AlarmChange) GetTriggers() []string {
	return GetTriggers(ac.Type)
}

func GetTriggers(t AlarmChangeType) []string {
	var triggers []string

	switch t {
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
		AlarmChangeTypeAutoDeclareTicketWebhookFail:
		// not a trigger
	case AlarmChangeTypeDeclareTicketWebhook,
		AlarmChangeTypeAutoDeclareTicketWebhook:
		triggers = append(triggers, string(AlarmChangeTypeDeclareTicketWebhook))
	default:
		trigger := string(t)
		if trigger != "" {
			triggers = append(triggers, trigger)
		}
	}

	return triggers
}
