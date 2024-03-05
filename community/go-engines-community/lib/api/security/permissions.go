package security

const (
	ObjPbehavior          = "api_pbehavior"
	ObjPbehaviorType      = "api_pbehaviortype"
	ObjPbehaviorReason    = "api_pbehaviorreason"
	ObjPbehaviorException = "api_pbehaviorexception"

	ObjAction = "api_action"

	ObjEntity         = "api_entity"
	ObjEntityService  = "api_entityservice"
	ObjEntityCategory = "api_entitycategory"
	ObjContextGraph   = "api_contextgraph"

	ObjView           = "api_view"
	ObjWidgetTemplate = "api_widgettemplate"
	ObjViewGroup      = "api_viewgroup"
	ObjPlaylist       = "api_playlist"

	PermAlarmRead   = "api_alarm_read"
	PermAlarmUpdate = "api_alarm_update"

	PermAcl = "api_acl"

	PermShareToken = "api_share_token"

	ObjStateSettings = "api_state_settings"

	PermDataStorageRead   = "api_datastorage_read"
	PermDataStorageUpdate = "api_datastorage_update"

	ObjEventFilter = "api_eventfilter"

	ObjBroadcastMessage = "api_broadcast_message"

	ObjAssociativeTable = "api_associative_table"

	PermUserInterfaceUpdate = "api_user_interface_update"
	PermUserInterfaceDelete = "api_user_interface_delete"

	PermEvent = "api_event"

	ObjIdleRule = "api_idlerule"

	PermNotification = "api_notification"

	PermMessageRateStatsRead = "api_message_rate_stats_read"

	PermHealthcheck = "api_healthcheck"

	ObjFile = "api_file"

	ObjFlappingRule = "api_flapping_rule"
	ObjResolveRule  = "api_resolve_rule"

	PermCorporatePattern = "api_corporate_pattern"

	PermExportConfigurations = "api_export_configurations"

	PermTechMetrics = "api_techmetrics"

	ObjLinkRule = "api_link_rule"

	ObjAlarmTag = "api_alarm_tag"

	PermMaintenance = "api_maintenance"

	ObjColorTheme = "api_color_theme"

	PermPrivateViewGroups = "api_private_view_groups"

	PermIcon = "api_icon"
)

// PermCheck defines the permission check configuration where Obj is an object and Act is an action
type PermCheck struct {
	Obj string
	Act string
}
