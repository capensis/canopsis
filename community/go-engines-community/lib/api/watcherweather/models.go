package watcherweather

import (
	alarmapi "git.canopsis.net/canopsis/go-engines/lib/api/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehaviortype"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" json:"sort_by" binding:"oneoforempty=name state infos.*"`
}

type Watcher struct {
	ID             string              `json:"_id" bson:"_id"`
	Name           string              `json:"name" bson:"name"`
	Infos          map[string]Info     `json:"infos" bson:"infos"`
	Connector      string              `json:"connector" bson:"connector"`
	ConnectorName  string              `json:"connector_name" bson:"connector_name"`
	Component      string              `json:"component" bson:"component"`
	Resource       string              `json:"resource" bson:"resource"`
	HasOpenAlarm   bool                `json:"is_action_required" bson:"has_open_alarm"`
	State          alarmapi.AlarmStep  `json:"state" bson:"state"`
	Status         alarmapi.AlarmStep  `json:"status" bson:"status"`
	Snooze         *alarmapi.AlarmStep `json:"snooze" bson:"snooze"`
	Ack            *alarmapi.AlarmStep `json:"ack" bson:"ack"`
	Icon           string              `json:"icon" bson:"icon"`
	SecondaryIcon  string              `json:"secondary_icon" bson:"secondary_icon"`
	Color          string              `json:"color" bson:"color"`
	Output         string              `json:"output" bson:"output"`
	LastUpdateDate types.CpsTime       `json:"last_update_date" bson:"last_update_date" swaggertype:"integer"`
	AlarmCounters  []AlarmCounter      `json:"alarm_counters" bson:"alarm_counters"`
	Links          []struct {
		Name  string      `json:"cat_name" bson:"cat_name"`
		Links interface{} `json:"links" bson:"links"`
	} `json:"linklist" bson:"links"`
	Pbehaviors  []pbehavior.PBehavior `json:"pbehaviors" bson:"-"`
	PbehaviorID string                `json:"-" bson:"pbehavior_id"`
	// Keep for compatibility.
	EntityID            string                `json:"entity_id" bson:"entity_id"`
	DisplayName         string                `json:"display_name" bson:"display_name"`
	TileIcon            string                `json:"tileIcon" bson:"tileIcon"`
	TileSecondaryIcon   string                `json:"tileSecondaryIcon" bson:"tileSecondaryIcon"`
	TileColor           string                `json:"tileColor" bson:"tileColor"`
	WatcherPbehaviors   []pbehavior.PBehavior `json:"watcher_pbehavior" bson:"-"`
	SlaTex              string                `json:"sla_tex" bson:"-"`
	MFilter             string                `json:"mfilter" bson:"mfilter"`
	IsActionRequired    bool                  `json:"isActionRequired" bson:"-"`
	IsAllEntitiesPaused bool                  `json:"isAllEntitiesPaused" bson:"is_all_watched_inactive"`
	IsWatcherPaused     bool                  `json:"isWatcherPaused" bson:"is_inactive"`
}

type Info struct {
	Description string      `bson:"description,omitempty" json:"description"`
	Value       interface{} `bson:"value" json:"value"`
}

type AlarmCounter struct {
	Count int64              `json:"count" bson:"count"`
	Type  pbehaviortype.Type `json:"type" bson:"type"`
}

type AggregationResult struct {
	Data       []Watcher `bson:"data" json:"data"`
	TotalCount int64     `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

type Entity struct {
	ID             string                `json:"_id" bson:"_id"`
	Name           string                `json:"name" bson:"name"`
	Infos          map[string]Info       `json:"infos" bson:"infos"`
	Type           string                `json:"source_type" bson:"type"`
	Connector      string                `json:"connector" bson:"connector"`
	ConnectorName  string                `json:"connector_name" bson:"connector_name"`
	Component      string                `json:"component" bson:"component"`
	Resource       string                `json:"resource" bson:"resource"`
	State          alarmapi.AlarmStep    `json:"state" bson:"state"`
	Status         alarmapi.AlarmStep    `json:"status" bson:"status"`
	Snooze         *alarmapi.AlarmStep   `json:"snooze" bson:"snooze"`
	Ack            *alarmapi.AlarmStep   `json:"ack" bson:"ack"`
	Ticket         *alarmapi.AlarmTicket `json:"ticket" bson:"ticket"`
	LastUpdateDate types.CpsTime         `json:"last_update_date" bson:"last_update_date" swaggertype:"integer"`
	CreationDate   types.CpsTime         `json:"alarm_creation_date" bson:"creation_date" swaggertype:"integer"`
	DisplayName    string                `json:"alarm_display_name" bson:"display_name"`
	Icon           string                `json:"icon" bson:"icon"`
	Color          string                `json:"color" bson:"color"`
	Pbehaviors     []pbehavior.PBehavior `json:"pbehaviors" bson:"-"`
	PbehaviorID    string                `json:"-" bson:"pbehavior_id"`
	Links          []WeatherLink         `json:"linklist" bson:"-"`
	// Keep for compatibility.
	EntityID         string                `json:"entity_id" bson:"entity_id"`
	SlaText          string                `json:"sla_text" bson:"-"`
	Org              string                `json:"org" bson:"org"`
	Stats            Stats                 `json:"stats" bson:"-"`
	IsInactive       bool                  `json:"-" bson:"is_inactive"`
	LastPbhLeaveDate *types.CpsTime        `json:"-" bson:"last_pbhleave_date"`
	EntityPbehaviors []pbehavior.PBehavior `json:"pbehavior" bson:"-"`
}

type WeatherLink struct {
	Category string      `json:"cat_name"`
	Links    interface{} `json:"links"`
}

type Stats struct {
	FailEventsCount int            `json:"ko" bson:"ko"`
	OKEventsCount   int            `json:"ok" bson:"ok"`
	LastFailEvent   *types.CpsTime `json:"last_ko" bson:"last_ko" swaggertype:"integer"`
	LastEvent       *types.CpsTime `json:"last_event" bson:"last_event" swaggertype:"integer"`
}

type EntityAggregationResult struct {
	Data       []Entity `bson:"data" json:"data"`
	TotalCount int64    `bson:"total_count" json:"total_count"`
}

func (r *EntityAggregationResult) GetData() interface{} {
	return r.Data
}

func (r *EntityAggregationResult) GetTotal() int64 {
	return r.TotalCount
}
