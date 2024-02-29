package mongo

const (
	DB                                = "canopsis"
	ConfigurationMongoCollection      = "configuration"
	SessionMongoCollection            = "session"
	AlarmMongoCollection              = "periodical_alarm"
	EntityMongoCollection             = "default_entities"
	PbehaviorMongoCollection          = "pbehavior"
	PbehaviorTypeMongoCollection      = "pbehavior_type"
	PbehaviorReasonMongoCollection    = "pbehavior_reason"
	PbehaviorExceptionMongoCollection = "pbehavior_exception"
	FileMongoCollection               = "files"
	MetaAlarmRulesMongoCollection     = "meta_alarm_rules"
	IdleRuleMongoCollection           = "idle_rule"
	ExportTaskMongoCollection         = "export_task"
	ActionLogMongoCollection          = "action_log"
	DynamicInfosRulesMongoCollection  = "dynamic_infos"
	EntityCategoryMongoCollection     = "entity_category"
	ImportJobMongoCollection          = "default_importgraph"
	JunitTestSuiteMongoCollection     = "junit_test_suite"
	JunitTestCaseMediaMongoCollection = "junit_test_case_media"
	PlaylistMongoCollection           = "view_playlist"
	StateSettingsMongoCollection      = "state_settings"
	BroadcastMessageMongoCollection   = "broadcast_message"
	AssociativeTableCollection        = "default_associativetable"
	NotificationMongoCollection       = "notification"

	ViewMongoCollection           = "views"
	ViewTabMongoCollection        = "viewtabs"
	WidgetMongoCollection         = "widgets"
	WidgetFiltersMongoCollection  = "widget_filters"
	WidgetTemplateMongoCollection = "widget_templates"
	ViewGroupMongoCollection      = "viewgroups"

	// Following collections are used for event statistics.
	MessageRateStatsMinuteCollectionName = "message_rate_statistic_minute"
	MessageRateStatsHourCollectionName   = "message_rate_statistic_hour"

	// Collection for ok/ko event statistics
	EventStatistics = "event_statistics"

	// Remediation collections
	InstructionMongoCollection          = "instruction"
	InstructionExecutionMongoCollection = "instruction_execution"
	InstructionRatingMongoCollection    = "instruction_rating"
	JobConfigMongoCollection            = "job_config"
	JobMongoCollection                  = "job"
	JobHistoryMongoCollection           = "job_history"

	// Data storage alarm collections
	ResolvedAlarmMongoCollection = "resolved_alarms"
	ArchivedAlarmMongoCollection = "archived_alarms"
	// Data storage entity collections
	ArchivedEntitiesMongoCollection = "archived_entities"

	TokenMongoCollection               = "token"
	ShareTokenMongoCollection          = "share_token"
	WebsocketConnectionMongoCollection = "websocket_connection"

	ResolveRuleMongoCollection  = "resolve_rule"
	FlappingRuleMongoCollection = "flapping_rule"

	UserPreferencesMongoCollection = "userpreferences"

	KpiFilterMongoCollection = "kpi_filter"

	PatternMongoCollection = "pattern"

	EntityInfosDictionaryCollection  = "entity_infos_dictionary"
	DynamicInfosDictionaryCollection = "dynamic_infos_dictionary"

	MapMongoCollection = "map"

	AlarmTagCollection      = "alarm_tag"
	AlarmTagColorCollection = "alarm_tag_color"

	MibCollection       = "default_mibs"
	SnmpRulesCollection = "default_snmprules"

	ScenarioMongoCollection          = "action_scenario"
	DeclareTicketRuleMongoCollection = "declare_ticket_rule"
	WebhookHistoryMongoCollection    = "webhook_history"

	LinkRuleMongoCollection = "link_rule"

	UserCollection         = "user"
	RoleCollection         = "role"
	RoleTemplateCollection = "role_template"
	PermissionCollection   = "permission"

	EventFilterRuleCollection    = "eventfilter"
	EventFilterFailureCollection = "eventfilter_failure"

	EntityCountersCollection = "entity_counters"

	MetaAlarmStatesCollection = "meta_alarm_states"

	ColorThemeCollection = "color_theme"

	EngineNotificationCollection = "engine_notification"

	IconCollection = "icon"
)
