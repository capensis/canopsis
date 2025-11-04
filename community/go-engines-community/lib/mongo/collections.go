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
	DynamicInfosRulesMongoCollection  = "dynamic_infos"
	EntityCategoryMongoCollection     = "entity_category"
	ImportJobMongoCollection          = "default_importgraph"
	JunitTestSuiteMongoCollection     = "junit_test_suite"
	JunitTestCaseMediaMongoCollection = "junit_test_case_media"
	PlaylistMongoCollection           = "view_playlist"
	StateSettingsMongoCollection      = "state_settings"
	BroadcastMessageCollection        = "broadcast_message"
	BroadcastMessageReadCollection    = "broadcast_message_read"
	AssociativeTableCollection        = "default_associativetable"

	ViewMongoCollection           = "views"
	ViewTabMongoCollection        = "viewtabs"
	WidgetMongoCollection         = "widgets"
	WidgetFiltersMongoCollection  = "widget_filters"
	WidgetTemplateMongoCollection = "widget_templates"
	ViewGroupMongoCollection      = "viewgroups"

	// Collection for ok/ko event statistics
	EventStatistics = "event_statistics"

	EventRecordsMongoCollection = "event_records"

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

	ScenarioCollection            = "action_scenario"
	DeclareTicketRuleCollection   = "declare_ticket_rule"
	WebhookTokenRuleCollection    = "webhook_token_rule" //nolint:gosec
	WebhookHistoryCollection      = "webhook_history"
	WebhookTokenHistoryCollection = "webhook_token_history" //nolint:gosec

	LinkRuleMongoCollection = "link_rule"

	UserCollection            = "user"
	RoleCollection            = "role"
	RoleTemplateCollection    = "role_template"
	PermissionCollection      = "permission"
	PermissionGroupCollection = "permission_group"

	EventFilterRuleCollection    = "eventfilter"
	EventFilterFailureCollection = "eventfilter_failure"

	EntityCountersCollection = "entity_counters"

	MetaAlarmStatesCollection = "meta_alarm_states"

	ColorThemeCollection = "color_theme"

	EngineNotificationCollection = "engine_notification"

	IconCollection = "icon"

	CloseDelayJobCollection = "close_delay_job"

	ExternalDataTableCollection        = "external_data_table"
	ExternalDataImportWorkerCollection = "external_data_import_worker"

	EntityInfosPropertyCollection = "entity_infos_property"

	UserNotificationSettingsCollection = "user_notification_settings"
	UserNotificationCollection         = "user_notification"

	TemplateTestDataCollection = "template_data"
	TemplateTestCollection     = "template_test"

	CommentTemplateMongoCollection = "comment_template"
)
