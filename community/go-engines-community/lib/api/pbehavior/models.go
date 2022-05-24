package pbehavior

import (
	"encoding/json"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorreason"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" json:"sort_by" binding:"oneoforempty=name author enabled tstart tstop type.name reason.name created updated rrule type.icon_name last_alarm_date"`
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
	Start      *types.CpsTime                     `json:"tstart" binding:"required" swaggertype:"integer"`
	Stop       *types.CpsTime                     `json:"tstop" swaggertype:"integer"`
	Type       string                             `json:"type" binding:"required"`
	Exdates    []pbehaviorexception.ExdateRequest `json:"exdates" binding:"dive"`
	Exceptions []string                           `json:"exceptions"`

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

// for swagger
type BulkCreateResponseItem struct {
	ID     string            `json:"id,omitempty"`
	Item   CreateRequest     `json:"item"`
	Status int               `json:"status"`
	Error  string            `json:"error,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

// for swagger
type BulkUpdateResponseItem struct {
	ID     string                `json:"id,omitempty"`
	Item   BulkUpdateRequestItem `json:"item"`
	Status int                   `json:"status"`
	Error  string                `json:"error,omitempty"`
	Errors map[string]string     `json:"errors,omitempty"`
}

// for swagger
type BulkDeleteResponseItem struct {
	ID     string                `json:"id,omitempty"`
	Item   BulkDeleteRequestItem `json:"item"`
	Status int                   `json:"status"`
	Error  string                `json:"error,omitempty"`
	Errors map[string]string     `json:"errors,omitempty"`
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

	EntityPattern          pattern.Entity             `json:"entity_pattern"`
	CorporateEntityPattern *string                    `json:"corporate_entity_pattern"`
	CorporatePattern       *savedpattern.SavedPattern `json:"-"`
}

type FilterRequest struct {
	EntityPattern pattern.Entity `json:"entity_pattern" binding:"entity_pattern" binding:"required"`
}

type FindByEntityIDRequest struct {
	ID string `form:"id" binding:"required"`
}

type Response struct {
	ID            string                         `bson:"_id" json:"_id"`
	Author        string                         `bson:"author" json:"author"`
	Comments      pbehavior.Comments             `bson:"comments" json:"comments"`
	Enabled       bool                           `bson:"enabled" json:"enabled"`
	Name          string                         `bson:"name" json:"name"`
	Reason        *pbehaviorreason.Reason        `bson:"reason" json:"reason"`
	RRule         string                         `bson:"rrule" json:"rrule"`
	Start         *types.CpsTime                 `bson:"tstart" json:"tstart" swaggertype:"integer"`
	Stop          *types.CpsTime                 `bson:"tstop" json:"tstop" swaggertype:"integer"`
	Created       *types.CpsTime                 `bson:"created" json:"created" swaggertype:"integer"`
	Updated       *types.CpsTime                 `bson:"updated" json:"updated" swaggertype:"integer"`
	Type          *pbehavior.Type                `bson:"type" json:"type"`
	Exdates       []pbehaviorexception.Exdate    `bson:"exdates" json:"exdates"`
	Exceptions    []pbehaviorexception.Exception `bson:"exceptions" json:"exceptions"`
	LastAlarmDate *types.CpsTime                 `bson:"last_alarm_date,omitempty" json:"last_alarm_date" swaggertype:"integer"`
	// IsActiveStatus represents if pbehavior is in action for current time.
	IsActiveStatus *bool `bson:"-" json:"is_active_status,omitempty"`

	OldMongoQuery map[string]interface{} `bson:"old_mongo_query" json:"old_mongo_query,omitempty"`

	savedpattern.EntityPatternFields `bson:",inline"`
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
