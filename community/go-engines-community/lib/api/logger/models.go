package logger

import "time"

const (
	ApprovalDecisionApprove = "approve"
	ApprovalDecisionDismiss = "dismiss"
)

const (
	ValueTypeUser               = "user"
	ValueTypeRole               = "role"
	ValueTypePlayList           = "playlist"
	ValueTypeEventFilter        = "eventfilter"
	ValueTypeScenario           = "scenario"
	ValueTypeMetaalarmRule      = "metaalarmrule"
	ValueTypeDynamicInfo        = "dynamicinfo"
	ValueTypeEntity             = "entity"
	ValueTypeEntityService      = "entityservice"
	ValueTypeEntityCategory     = "entitycategory"
	ValueTypePbehaviorType      = "pbehaviortype"
	ValueTypePbehaviorReason    = "pbehaviorreason"
	ValueTypePbehaviorException = "pbehaviorexception"
	ValueTypePbehavior          = "pbehavior"
	ValueTypeJobConfig          = "jobconfig"
	ValueTypeJob                = "job"
	ValueTypeInstruction        = "instruction"
	ValueTypeStateSetting       = "statesetting"
	ValueTypeBroadcastMessage   = "broadcastmessage"
	ValueTypeIdleRule           = "idlerule"
	ValueTypeView               = "view"
	ValueTypeViewTab            = "viewtab"
	ValueTypeWidget             = "widget"
	ValueTypeWidgetFilter       = "widgetfilter"
	ValueTypeWidgetTemplate     = "widgettemplate"
	ValueTypeViewGroup          = "viewgroup"
	ValueTypeResolveRule        = "resolverule"
	ValueTypeFlappingRule       = "flappingrule"
	ValueTypeKpiFilter          = "kpi_filter"
	ValueTypePattern            = "pattern"
	ValueTypeMap                = "map"
	ValueTypeSnmpRule           = "snmprule"
	ValueTypeDeclareTicketRule  = "declareticketrule"
	ValueTypeLinkRule           = "linkrule"
	ValueTypeAlarmTag           = "alarmtag"
	ValueTypeColorTheme         = "colortheme"
	ValueTypeIcon               = "icon"
)

type ActionLogEvent struct {
	DocumentID        string         `bson:"document_id"`
	Collection        string         `bson:"collection"`
	OperationType     string         `bson:"operation_type"`
	Document          map[string]any `bson:"document"`
	DocumentBefore    map[string]any `bson:"document_before"`
	UpdateDescription map[string]any `bson:"update_description"`
	ClusterTime       time.Time      `bson:"cluster_time"`
}

type ActionCreateLog struct {
	ValueType    string
	ValueID      string
	InitialValue map[string]any
	Author       string
	Timestamp    time.Time
}

type ActionUpdateLog struct {
	ValueType string
	ValueID   string
	Author    string
	Timestamp time.Time

	PrevValue         map[string]any
	UpdateDescription map[string]any
}

type ActionDeleteLog struct {
	ValueType string
	ValueID   string
	PrevValue map[string]any
	Author    string
	Timestamp time.Time
}
