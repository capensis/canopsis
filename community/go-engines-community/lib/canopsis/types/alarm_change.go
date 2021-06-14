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
	AlarmChangeTypeAckremove         AlarmChangeType = "ackremove"
	AlarmChangeTypeCancel            AlarmChangeType = "cancel"
	AlarmChangeTypeUncancel          AlarmChangeType = "uncancel"
	AlarmChangeTypeAssocTicket       AlarmChangeType = "assocticket"
	AlarmChangeTypeSnooze            AlarmChangeType = "snooze"
	AlarmChangeTypeUnsnooze          AlarmChangeType = "unsnooze"
	AlarmChangeTypeResolve           AlarmChangeType = "resolve"
	AlarmChangeTypeDone              AlarmChangeType = "done"
	AlarmChangeTypeComment           AlarmChangeType = "comment"
	AlarmChangeTypeChangeState       AlarmChangeType = "changestate"
	AlarmChangeTypePbhEnter          AlarmChangeType = "pbhenter"
	AlarmChangeTypePbhLeave          AlarmChangeType = "pbhleave"
	AlarmChangeTypePbhLeaveAndEnter  AlarmChangeType = "pbhleaveandenter"
	AlarmChangeTypeUpdateStatus      AlarmChangeType = "changestatus"
	AlarmChangeTypeActivate          AlarmChangeType = "activate"

	// AlarmChangeTypeDeclareTicket is used for manual declareticket trigger which is designed
	// to trigger webhook with declare ticket parameter.
	AlarmChangeTypeDeclareTicket AlarmChangeType = "declareticket"
	// AlarmChangeTypeDeclareTicketWebhook is triggered after declare ticket creation by webhook.
	AlarmChangeTypeDeclareTicketWebhook AlarmChangeType = "declareticketwebhook"

	// Following consts are used for instruction.
	AlarmChangeTypeInstructionStart       AlarmChangeType = "instructionstart"
	AlarmChangeTypeInstructionPause       AlarmChangeType = "instructionpause"
	AlarmChangeTypeInstructionResume      AlarmChangeType = "instructionresume"
	AlarmChangeTypeInstructionComplete    AlarmChangeType = "instructioncomplete"
	AlarmChangeTypeInstructionAbort       AlarmChangeType = "instructionabort"
	AlarmChangeTypeInstructionFail        AlarmChangeType = "instructionfail"
	AlarmChangeTypeInstructionJobStart    AlarmChangeType = "instructionjobstart"
	AlarmChangeTypeInstructionJobComplete AlarmChangeType = "instructionjobcomplete"
	AlarmChangeTypeInstructionJobAbort    AlarmChangeType = "instructionjobabort"
	AlarmChangeTypeInstructionJobFail     AlarmChangeType = "instructionjobfail"
)

// AlarmChange is a struct containing the type of change that occured on an
// alarm, as well as its previous state.
type AlarmChange struct {
	Type                 AlarmChangeType
	PreviousState        CpsNumber
	PreviousStateChange  CpsTime
	PreviousStatus       CpsNumber
	PreviousStatusChange CpsTime
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
	var triggers []string

	switch ac.Type {
	case AlarmChangeTypeCreateAndPbhEnter:
		triggers = append(triggers, string(AlarmChangeTypeCreate), string(AlarmChangeTypePbhEnter))
	case AlarmChangeTypePbhLeaveAndEnter:
		triggers = append(triggers, string(AlarmChangeTypePbhEnter), string(AlarmChangeTypePbhLeave))
	default:
		t := string(ac.Type)
		if t != "" {
			triggers = append(triggers, t)
		}
	}

	return triggers
}
