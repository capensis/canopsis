package builder

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

var (
	staticFields = map[string]map[string]string{
		PatternTypeAlarm: {
			"v.display_name":           pattern.FieldTypeString,
			"v.output":                 pattern.FieldTypeString,
			"v.long_output":            pattern.FieldTypeString,
			"v.initial_output":         pattern.FieldTypeString,
			"v.initial_long_output":    pattern.FieldTypeString,
			"v.connector":              pattern.FieldTypeString,
			"v.connector_name":         pattern.FieldTypeString,
			"v.component":              pattern.FieldTypeString,
			"v.resource":               pattern.FieldTypeString,
			"v.state.val":              pattern.FieldTypeInt,
			"v.status.val":             pattern.FieldTypeInt,
			"v.total_state_changes":    pattern.FieldTypeInt,
			"v.duration":               pattern.FieldTypeDuration,
			"v.creation_date":          pattern.FieldTypeTimestamp,
			"v.last_event_date":        pattern.FieldTypeTimestamp,
			"v.last_update_date":       pattern.FieldTypeTimestamp,
			"v.resolved":               pattern.FieldTypeTimestamp,
			"v.ack":                    pattern.FieldTypeReference,
			"v.ack.a":                  pattern.FieldTypeString,
			"v.ack.m":                  pattern.FieldTypeString,
			"v.ack.t":                  pattern.FieldTypeTimestamp,
			"v.ack.initiator":          pattern.FieldTypeString,
			"v.ticket":                 pattern.FieldTypeReference,
			"v.ticket.m":               pattern.FieldTypeString,
			"v.ticket.ticket":          pattern.FieldTypeString,
			"v.ticket.initiator":       pattern.FieldTypeString,
			"v.canceled":               pattern.FieldTypeReference,
			"v.canceled.initiator":     pattern.FieldTypeString,
			"v.snooze":                 pattern.FieldTypeReference,
			"v.snooze.a":               pattern.FieldTypeString,
			"v.snooze.initiator":       pattern.FieldTypeString,
			"v.state.initiator":        pattern.FieldTypeString,
			"v.last_comment.m":         pattern.FieldTypeString,
			"v.last_comment.a":         pattern.FieldTypeString,
			"v.last_comment.initiator": pattern.FieldTypeString,
			"v.change_state":           pattern.FieldTypeReference,
			"v.failed_ticket":          pattern.FieldTypeReference,
			"v.meta":                   pattern.FieldTypeString,
			"tags":                     pattern.FieldTypeTags,
		},
		PatternTypeEntity: {
			"_id":             pattern.FieldTypeString,
			"name":            pattern.FieldTypeString,
			"connector":       pattern.FieldTypeString,
			"component":       pattern.FieldTypeString,
			"category":        pattern.FieldTypeString,
			"type":            pattern.FieldTypeString,
			"impact_level":    pattern.FieldTypeInt,
			"last_event_date": pattern.FieldTypeTimestamp,
		},
		PatternTypePbehavior: {
			"pbehavior_info.id":             pattern.FieldTypeString,
			"pbehavior_info.type":           pattern.FieldTypeString,
			"pbehavior_info.reason":         pattern.FieldTypeString,
			"pbehavior_info.canonical_type": pattern.FieldTypeString,
		},
		PatternTypeEvent: {
			"connector":      pattern.FieldTypeString,
			"connector_name": pattern.FieldTypeString,
			"component":      pattern.FieldTypeString,
			"resource":       pattern.FieldTypeString,
			"output":         pattern.FieldTypeString,
			"long_output":    pattern.FieldTypeString,
			"author":         pattern.FieldTypeString,
			"event_type":     pattern.FieldTypeString,
			"source_type":    pattern.FieldTypeString,
			"state":          pattern.FieldTypeInt,
			"initiator":      pattern.FieldTypeString,
		},
		PatternTypeWeatherService: {
			"is_grey":        pattern.FieldTypeBool,
			"icon":           pattern.FieldTypeString,
			"secondary_icon": pattern.FieldTypeString,
			"state.val":      pattern.FieldTypeInt,
		},
	}

	mixedFields = map[string]map[string]map[string]string{
		PatternTypeAlarm: {
			"v.activation_date": {
				pattern.ConditionExist:        pattern.FieldTypeReference,
				pattern.ConditionTimeRelative: pattern.FieldTypeTimestamp,
				pattern.ConditionTimeAbsolute: pattern.FieldTypeTimestamp,
			},
		},
	}

	dynamicFieldPrefixes = map[string]map[string]string{
		PatternTypeAlarm: {
			"v.infos.*":            "v.infos.",
			DynamicFieldTicketData: "v.ticket.ticket_data.",
		},
		PatternTypeEntity: {
			"infos.*":           "infos.",
			"component_infos.*": "component_infos.",
		},
		PatternTypeEvent: {
			"extra.*": "extra.",
		},
	}

	initiatorEnumValues = []string{
		types.InitiatorUser,
		types.InitiatorSystem,
		types.InitiatorExternal,
	}

	eventTypeEnumValues = []string{
		types.EventTypeAck,
		types.EventTypeAckremove,
		types.EventTypeAssocTicket,
		types.EventTypeTicketRemove,
		types.EventTypeCancel,
		types.EventTypeCheck,
		types.EventTypeComment,
		types.EventTypeChangestate,
		types.EventTypeSnooze,
		types.EventTypeUnsnooze,
		types.EventTypeUncancel,
		types.EventTypeContextUpdate,
		types.EventTypeDeclareTicketWebhook,
		types.EventTypeWebhookStarted,
		types.EventTypeWebhookCompleted,
		types.EventTypeWebhookFailed,
		types.EventTypeAutoWebhookStarted,
		types.EventTypeAutoWebhookCompleted,
		types.EventTypeAutoWebhookFailed,
		types.EventTypePbhEnter,
		types.EventTypePbhLeaveAndEnter,
		types.EventTypePbhLeave,
		types.EventTypeResolveCancel,
		types.EventTypeResolveClose,
		types.EventTypeResolveDeleted,
		types.EventTypeUpdateStatus,
		types.EventTypeActivate,
		types.EventTypeRunDelayedScenario,
		types.EventTypeMetaAlarm,
		types.EventTypeMetaAlarmAttachChildren,
		types.EventTypeMetaAlarmDetachChildren,
		types.EventTypeInstructionStarted,
		types.EventTypeInstructionPaused,
		types.EventTypeInstructionResumed,
		types.EventTypeInstructionCompleted,
		types.EventTypeInstructionFailed,
		types.EventTypeInstructionAborted,
		types.EventTypeAutoInstructionStarted,
		types.EventTypeAutoInstructionCompleted,
		types.EventTypeAutoInstructionFailed,
		types.EventTypeInstructionJobStarted,
		types.EventTypeInstructionJobCompleted,
		types.EventTypeInstructionJobFailed,
		types.EventTypeRecomputeEntityService,
		types.EventTypeEntityUpdated,
		types.EventTypeEntityToggled,
		types.EventTypeUpdateCounters,
		types.EventTypeJunitTestSuiteUpdated,
		types.EventTypeJunitTestCaseUpdated,
		types.EventTypeNoEvents,
		types.EventTypeTrigger,
		types.EventTypeAutoInstructionActivate,
		types.EventTypeMetaAlarmChildActivate,
		types.EventTypeMetaAlarmChildDeactivate,
		types.EventTypeChangeTicketStatus,
	}

	alarmStateEnumValues = []int64{
		types.AlarmStateOK,
		types.AlarmStateMinor,
		types.AlarmStateMajor,
		types.AlarmStateCritical,
	}

	alarmStatusEnumValues = []int64{
		types.AlarmStatusOff,
		types.AlarmStatusOngoing,
		types.AlarmStatusStealthy,
		types.AlarmStatusFlapping,
		types.AlarmStatusCancelled,
		types.AlarmStatusNoEvents,
		types.AlarmStatusUnknown,
	}

	stringEnumValues = map[string][]string{
		"source_type": {
			types.SourceTypeResource,
			types.SourceTypeComponent,
			types.SourceTypeConnector,
			types.SourceTypeService,
		},
		"event_type":               eventTypeEnumValues,
		"initiator":                initiatorEnumValues,
		"v.ack.initiator":          initiatorEnumValues,
		"v.ticket.initiator":       initiatorEnumValues,
		"v.canceled.initiator":     initiatorEnumValues,
		"v.snooze.initiator":       initiatorEnumValues,
		"v.state.initiator":        initiatorEnumValues,
		"v.last_comment.initiator": initiatorEnumValues,
		"type": {
			types.EntityTypeConnector,
			types.EntityTypeComponent,
			types.EntityTypeResource,
			types.EntityTypeService,
		},
		"pbehavior_info.canonical_type": {
			pbehavior.TypeActive,
			pbehavior.TypePause,
			pbehavior.TypeMaintenance,
			pbehavior.TypeInactive,
		},
	}

	intEnumValues = map[string][]int64{
		"v.state.val":  alarmStateEnumValues,
		"state":        alarmStateEnumValues, // event_pattern
		"state.val":    alarmStateEnumValues, // weather_service_pattern
		"v.status.val": alarmStatusEnumValues,
	}
)
