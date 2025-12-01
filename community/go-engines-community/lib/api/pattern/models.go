package pattern

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

const (
	OptimizeStatusCreated = iota
	OptimizeStatusRunning
	OptimizeStatusSucceeded
	OptimizeStatusFailed
	OptimizeStatusAccepted
	OptimizeStatusRejected
)

const (
	EntityFieldComponent = "component"
	EntityFieldName      = "name"
	EntityInfosPrefix    = "infos."
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy    string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id title author.name author.display_name created updated"`
	Corporate *bool  `json:"corporate" form:"corporate"`
	Type      string `json:"type" form:"type"`
}

type EditRequest struct {
	ID                    string                        `json:"-"`
	Title                 string                        `json:"title" binding:"required,max=255"`
	Type                  string                        `json:"type" binding:"required,oneof=alarm entity pbehavior weather_service"`
	IsCorporate           *bool                         `json:"is_corporate" binding:"required"`
	AlarmPattern          pattern.Alarm                 `json:"alarm_pattern" binding:"alarm_pattern"`
	EntityPattern         pattern.Entity                `json:"entity_pattern" binding:"entity_pattern"`
	PbehaviorPattern      pattern.PbehaviorInfo         `json:"pbehavior_pattern" binding:"pbehavior_pattern"`
	WeatherServicePattern pattern.WeatherServicePattern `json:"weather_service_pattern" binding:"weather_service_pattern"`
	Author                string                        `json:"author" swaggerignore:"true"`
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}

type Response struct {
	ID                    string                        `bson:"_id" json:"_id"`
	Title                 string                        `bson:"title" json:"title"`
	Type                  string                        `bson:"type" json:"type"`
	IsCorporate           bool                          `bson:"is_corporate" json:"is_corporate"`
	AlarmPattern          pattern.Alarm                 `bson:"alarm_pattern" json:"alarm_pattern,omitempty"`
	EntityPattern         pattern.Entity                `bson:"entity_pattern" json:"entity_pattern,omitempty"`
	PbehaviorPattern      pattern.PbehaviorInfo         `bson:"pbehavior_pattern" json:"pbehavior_pattern,omitempty"`
	WeatherServicePattern pattern.WeatherServicePattern `bson:"weather_service_pattern" json:"weather_service_pattern,omitempty"`
	Author                *author.Author                `bson:"author" json:"author"`
	Created               datetime.CpsTime              `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated               datetime.CpsTime              `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
}

type AggregationResult struct {
	Data       []Response `bson:"data" json:"data"`
	TotalCount int64      `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

type CountRequest struct {
	AlarmPattern     pattern.Alarm         `json:"alarm_pattern" binding:"alarm_pattern"`
	EntityPattern    pattern.Entity        `json:"entity_pattern" binding:"entity_pattern"`
	PbehaviorPattern pattern.PbehaviorInfo `json:"pbehavior_pattern" binding:"pbehavior_pattern"`
}

type CountAlarmsResponse struct {
	AlarmPattern     CountResponse `json:"alarm_pattern"`
	EntityPattern    CountResponse `json:"entity_pattern"`
	PbehaviorPattern CountResponse `json:"pbehavior_pattern"`
	All              CountResponse `json:"all"`
	Entities         CountResponse `json:"entities"`
}

type CountEntitiesResponse struct {
	AlarmPattern     CountResponse `json:"alarm_pattern"`
	EntityPattern    CountResponse `json:"entity_pattern"`
	PbehaviorPattern CountResponse `json:"pbehavior_pattern"`
	All              CountResponse `json:"all"`
}

type CountResponse struct {
	Count     int64 `bson:"count" json:"count"`
	OverLimit bool  `bson:"-" json:"over_limit"`
	Millisecs int64 `json:"ms"`
}

type OptimizeRequest struct {
	EntityPattern pattern.Entity `json:"entity_pattern" binding:"entity_pattern"`
}

type Suggestion struct {
	EntityPattern pattern.Entity `json:"entity_pattern"`
	FoundEntities int            `json:"found_entities"`
	Difference    int            `json:"difference"`
}

type OptimizeJob struct {
	ID                 string            `bson:"_id" json:"_id"`
	Status             int               `bson:"status" json:"status"`
	EntityPattern      pattern.Entity    `bson:"entity_pattern" json:"-"`
	Created            datetime.CpsTime  `bson:"created" json:"-"`
	LastPing           *datetime.CpsTime `bson:"last_ping" json:"-"`
	Retries            int64             `bson:"retries" json:"-"`
	Suggestions        []Suggestion      `bson:"suggestions" json:"suggestions"`
	FailReason         string            `bson:"fail_reason,omitempty" json:"fail_reason,omitempty"`
	AcceptedSuggestion *int              `bson:"accepted_suggestion,omitempty" json:"accepted_suggestion,omitempty"`
}

type OptimizeAcceptRequest struct {
	ID     string `json:"-"`
	Index  *int   `json:"index" binding:"required"`
	Accept *bool  `json:"accept" binding:"required"`
}

type LiteralFieldStats struct {
	FieldName         string `bson:"k"`
	LiteralFoundTimes int    `bson:"v"`
}
