package pbehavior

import (
	"encoding/json"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorcomment"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorreason"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviortype"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" json:"sort_by" binding:"oneoforempty=name author.name author.display_name enabled tstart tstop type.name reason.name created updated rrule type.icon_name last_alarm_date"`
}

type EntitiesListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" json:"sort_by" binding:"oneoforempty=_id name type"`
}

type EditRequest struct {
	Author     string                             `json:"author" swaggerignore:"true"`
	Enabled    *bool                              `json:"enabled" binding:"required"`
	Name       string                             `json:"name" binding:"required,max=255"`
	Reason     string                             `json:"reason" binding:"required"`
	RRule      string                             `json:"rrule"`
	Start      *datetime.CpsTime                  `json:"tstart" binding:"required" swaggertype:"integer"`
	Stop       *datetime.CpsTime                  `json:"tstop" swaggertype:"integer"`
	Type       string                             `json:"type" binding:"required"`
	Exdates    []pbehaviorexception.ExdateRequest `json:"exdates" binding:"dive"`
	Exceptions []string                           `json:"exceptions"`
	Color      string                             `json:"color" binding:"iscolororempty"`

	common.EntityPatternFieldsRequest
}

type CreateRequest struct {
	EditRequest
	ID string `json:"_id" binding:"id"`
}

type UpdateRequest struct {
	EditRequest
	ID string `json:"-"`
}

type BulkUpdateRequestItem struct {
	EditRequest
	ID string `json:"_id" binding:"required"`
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}

type PatchRequest struct {
	ID         string                             `json:"-"`
	Author     string                             `json:"author" swaggerignore:"true"`
	Name       *string                            `json:"name"`
	Enabled    *bool                              `json:"enabled"`
	Reason     *string                            `json:"reason"`
	Type       *string                            `json:"type"`
	Start      *int64                             `json:"tstart" swaggertype:"integer"`
	Stop       NullableTime                       `json:"tstop" swaggertype:"integer"`
	RRule      *string                            `json:"rrule"`
	Exdates    []pbehaviorexception.ExdateRequest `json:"exdates" binding:"dive"`
	Exceptions []string                           `json:"exceptions"`
	Color      *string                            `json:"color"`

	EntityPattern          pattern.Entity             `json:"entity_pattern"`
	CorporateEntityPattern *string                    `json:"corporate_entity_pattern"`
	CorporatePattern       *savedpattern.SavedPattern `json:"-"`
}

type FindByEntityIDRequest struct {
	ID        string `form:"_id" binding:"required"`
	WithFlags bool   `form:"with_flags"`
}

type Response struct {
	ID            string                        `bson:"_id" json:"_id"`
	Author        *author.Author                `bson:"author" json:"author"`
	Comments      []pbehaviorcomment.Response   `bson:"comments" json:"comments"`
	Enabled       bool                          `bson:"enabled" json:"enabled"`
	Name          string                        `bson:"name" json:"name"`
	Reason        *pbehaviorreason.Response     `bson:"reason" json:"reason"`
	RRule         string                        `bson:"rrule" json:"rrule"`
	RRuleEnd      *datetime.CpsTime             `bson:"rrule_end" json:"rrule_end" swaggertype:"integer"`
	Start         *datetime.CpsTime             `bson:"tstart" json:"tstart" swaggertype:"integer"`
	Stop          *datetime.CpsTime             `bson:"tstop" json:"tstop" swaggertype:"integer"`
	Created       *datetime.CpsTime             `bson:"created" json:"created" swaggertype:"integer"`
	Updated       *datetime.CpsTime             `bson:"updated" json:"updated" swaggertype:"integer"`
	Type          *pbehaviortype.Response       `bson:"type" json:"type"`
	Color         string                        `bson:"color" json:"color"`
	Exdates       []pbehaviorexception.Exdate   `bson:"exdates" json:"exdates"`
	Exceptions    []pbehaviorexception.Response `bson:"exceptions" json:"exceptions"`
	LastAlarmDate *datetime.CpsTime             `bson:"last_alarm_date,omitempty" json:"last_alarm_date" swaggertype:"integer"`
	AlarmCount    int64                         `bson:"alarm_count" json:"alarm_count"`
	// IsActiveStatus represents if pbehavior is in action for current time.
	IsActiveStatus *bool `bson:"-" json:"is_active_status,omitempty"`

	Origin   string `bson:"origin" json:"origin"`
	Editable *bool  `bson:"editable,omitempty" json:"editable,omitempty"`

	savedpattern.EntityPatternFields `bson:",inline"`

	RRuleComputedStart *datetime.CpsTime `bson:"rrule_cstart" json:"-"`
}

type NullableTime struct {
	val   *int64
	isSet bool
}

func (t *NullableTime) UnmarshalJSON(data []byte) error {
	t.isSet = true
	if string(data) == "null" {
		return nil
	}
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	t.val = &i
	return nil
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

type AggregationEntitiesResult struct {
	Data       []entity.Entity `bson:"data" json:"data"`
	TotalCount int64           `bson:"total_count" json:"total_count"`
}

func (r *AggregationEntitiesResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationEntitiesResult) GetTotal() int64 {
	return r.TotalCount
}

type DeleteByNameRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type CalendarByEntityIDRequest struct {
	ID   string           `form:"_id" json:"_id" binding:"required"`
	From datetime.CpsTime `form:"from" json:"from" binding:"required" swaggertype:"integer"`
	To   datetime.CpsTime `form:"to" json:"to" binding:"required" swaggertype:"integer"`
}

type CalendarResponse struct {
	ID    string           `json:"_id"`
	Title string           `json:"title"`
	Color string           `json:"color"`
	From  datetime.CpsTime `json:"from" swaggertype:"integer"`
	To    datetime.CpsTime `json:"to" swaggertype:"integer"`
	Type  pbehavior.Type   `json:"type"`
}

type BulkEntityCreateRequestItem struct {
	Author  string            `json:"author" swaggerignore:"true"`
	Entity  string            `json:"entity" binding:"required"`
	Origin  string            `json:"origin" binding:"required,max=255"`
	Name    string            `json:"name" binding:"required,max=255"`
	Reason  string            `json:"reason" binding:"required"`
	Start   *datetime.CpsTime `json:"tstart" binding:"required" swaggertype:"integer"`
	Stop    *datetime.CpsTime `json:"tstop" swaggertype:"integer"`
	Type    string            `json:"type" binding:"required"`
	Color   string            `json:"color" binding:"iscolororempty"`
	Comment string            `json:"comment"`
}

type BulkEntityDeleteRequestItem struct {
	Entity string `json:"entity" binding:"required"`
	Origin string `json:"origin" binding:"required"`
}

type BulkConnectorCreateRequestItem struct {
	Author   string            `json:"author" swaggerignore:"true"`
	Entities []string          `json:"entities" binding:"required,notblank"`
	Origin   string            `json:"origin" binding:"required,max=255"`
	Start    *datetime.CpsTime `json:"tstart" binding:"required" swaggertype:"integer"`
	Stop     *datetime.CpsTime `json:"tstop" binding:"required" swaggertype:"integer"`
	Comment  string            `json:"comment" binding:"required"`
	Name     string            `json:"name" binding:"required,max=255"`
	Reason   string            `json:"reason" binding:"required"`
	Type     string            `json:"type" binding:"required"`
	Color    string            `json:"color" binding:"iscolororempty"`
}

type BulkConnectorDeleteRequestItem struct {
	Author   string            `json:"author" swaggerignore:"true"`
	Entities []string          `json:"entities" binding:"required,notblank"`
	Origin   string            `json:"origin" binding:"required"`
	Start    *datetime.CpsTime `json:"tstart" binding:"required" swaggertype:"integer"`
	Stop     *datetime.CpsTime `json:"tstop" binding:"required" swaggertype:"integer"`
	Comment  string            `json:"comment" binding:"required"`
}
