package logger

import (
	"time"
)

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
	ValueTypeEventRecord        = "eventrecord"
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

type ActionLog struct {
	OperationType     string
	ValueType         string
	ValueID           string
	Timestamp         time.Time
	CurDocument       map[string]any
	PrevDocument      map[string]any
	UpdateDescription map[string]any
}

func (l *ActionLog) GetCurAuthor() string {
	if rawAuthor, ok := l.CurDocument["author"]; ok {
		if strAuthor, ok := rawAuthor.(string); ok {
			return strAuthor
		}
	}

	return ""
}

func (l *ActionLog) GetPrevAuthor() string {
	if rawAuthor, ok := l.PrevDocument["author"]; ok {
		if strAuthor, ok := rawAuthor.(string); ok {
			return strAuthor
		}
	}

	return ""
}

func (l *ActionLog) GetPrevCreated() time.Time {
	if rawCreated, ok := l.PrevDocument["created"]; ok {
		if intCreated, ok := rawCreated.(int64); ok {
			return time.Unix(intCreated, 0).UTC()
		}
	}

	return time.Time{}
}
