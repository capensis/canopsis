package metrics

const (
	TotalAlarmNumber        = "total_alarm_number"
	NonDisplayedAlarmNumber = "non_displayed_alarm_number"
	PbhAlarmNumber          = "pbh_alarm_number"
	InstructionAlarmNumber  = "instruction_alarm_number"
	TicketAlarmNumber       = "ticket_alarm_number"
	CorrelationAlarmNumber  = "correlation_alarm_number"
	AckAlarmNumber          = "ack_alarm_number"
	CancelAckAlarmNumber    = "cancel_ack_alarm_number"
	AckDuration             = "ack_duration"
	ResolveDuration         = "resolve_duration"
	SliDuration             = "sli_duration"
	AlarmStateChangeNumber  = "alarm_state_change_number"

	UserActivity = "user_activity"
	UserSessions = "user_sessions"
	TicketNumber = "ticket_number"

	ManualInstructionAssignedAlarms = "manual_instruction_assigned_alarms"
	ManualInstructionExecutedAlarms = "manual_instruction_executed_alarms"

	InstructionAssignedInstructions = "instruction_assigned_instructions"
	InstructionExecutedInstructions = "instruction_executed_instructions"

	EntityMetaData = "entities"
	UserMetaData   = "users"

	InstructionExecution             = "instruction_execution"
	InstructionExecutionHourly       = "instruction_execution_hourly"
	InstructionExecutionByModifiedOn = "instruction_execution_by_modified_on"

	NotAckedInHourAlarms      = "not_acked_in_hour_alarms"
	NotAckedInFourHoursAlarms = "not_acked_in_four_hours_alarms"
	NotAckedInDayAlarms       = "not_acked_in_day_alarms"

	PerfData     = "perf_data"
	PerfDataName = "perf_data_name"

	MessageRate       = "message_rate"
	MessageRateHourly = "message_rate_hourly"
)

// criteria type
const (
	EntityCriteriaType = iota
	UserCriteriaType
)
