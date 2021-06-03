package pbehavior

import (
	"encoding/json"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorreason"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" json:"sort_by" binding:"oneoforempty=name author enabled tstart tstop type.name reason.name created updated rrule type.icon_name"`
}

type EIDsListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" json:"sort_by" binding:"oneoforempty=id"`
}

type EditRequest struct {
	Author     string                             `json:"author" swaggerignore:"true"`
	Enabled    *bool                              `json:"enabled" binding:"required"`
	Filter     interface{}                        `json:"filter" binding:"required"`
	Name       string                             `json:"name" binding:"required,max=255"`
	Reason     string                             `json:"reason" binding:"required"`
	RRule      string                             `json:"rrule"`
	Start      types.CpsTime                      `json:"tstart" binding:"required" swaggertype:"integer"`
	Stop       *types.CpsTime                     `json:"tstop" swaggertype:"integer"`
	Type       string                             `json:"type" binding:"required"`
	Exdates    []pbehaviorexception.ExdateRequest `json:"exdates" binding:"dive"`
	Exceptions []string                           `json:"exceptions"`
}

type CreateRequest struct {
	EditRequest
	ID string `json:"_id" binding:"id"`
}

type UpdateRequest struct {
	EditRequest
	ID string `json:"-"`
}

type FilterRequest struct {
	Filter interface{} `json:"filter" binding:"required"`
}

type FindByEntityIDRequest struct {
	ID string `form:"id" binding:"required"`
}

type Response struct {
	ID         string                         `bson:"_id" json:"_id"`
	Author     string                         `bson:"author" json:"author"`
	Comments   pbehavior.Comments             `bson:"comments" json:"comments"`
	Enabled    bool                           `bson:"enabled" json:"enabled"`
	Filter     Filter                         `bson:"filter" json:"filter"`
	Name       string                         `bson:"name" json:"name"`
	Reason     *pbehaviorreason.Reason        `bson:"reason" json:"reason"`
	RRule      string                         `bson:"rrule" json:"rrule"`
	Start      *types.CpsTime                 `bson:"tstart" json:"tstart" swaggertype:"integer"`
	Stop       *types.CpsTime                 `bson:"tstop" json:"tstop" swaggertype:"integer"`
	Created    *types.CpsTime                 `bson:"created" json:"created" swaggertype:"integer"`
	Updated    *types.CpsTime                 `bson:"updated" json:"updated" swaggertype:"integer"`
	Type       *pbehavior.Type                `bson:"type" json:"type"`
	Exdates    []pbehaviorexception.Exdate    `bson:"exdates" json:"exdates"`
	Exceptions []pbehaviorexception.Exception `bson:"exceptions" json:"exceptions"`
	// IsActiveStatus represents if pbehavior is in action for current time.
	IsActiveStatus *bool `bson:"-" json:"is_active_status,omitempty"`
}

type Filter struct {
	v interface{}
}

func NewFilter(v interface{}) Filter {
	return Filter{v: v}
}

func (f Filter) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.v)
}

func (f *Filter) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	v, _, ok := bsoncore.ReadString(b)
	if !ok {
		return errors.New("invalid value, expected string")
	}

	err := json.Unmarshal([]byte(v), &f.v)
	return err
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

type EID struct {
	ID string `bson:"id" json:"id"`
}

type AggregationEIDsResult struct {
	Data       []EID `bson:"data" json:"data"`
	TotalCount int64 `bson:"total_count" json:"total_count"`
}

func (r *AggregationEIDsResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationEIDsResult) GetTotal() int64 {
	return r.TotalCount
}

type CountFilterResult struct {
	OverLimit  bool  `bson:"-" json:"over_limit"`
	TotalCount int64 `bson:"total_count" json:"total_count"`
}

func (r *CountFilterResult) GetTotal() int64 {
	return r.TotalCount
}