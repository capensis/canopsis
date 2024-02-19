package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Response struct {
	entity.Entity  `bson:",inline"`
	OutputTemplate string `bson:"output_template" json:"output_template"`
	SliAvailState  int64  `bson:"sli_avail_state" json:"sli_avail_state"`

	savedpattern.EntityPatternFields `bson:",inline"`
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

type ContextGraphRequest struct {
	pagination.Query
	entity.SortRequest
	ID        string `form:"_id" json:"_id"`
	Search    string `form:"search" json:"search"`
	Category  string `form:"category" json:"category"`
	WithFlags bool   `form:"with_flags" json:"with_flags"`
	// Show dependencies defining the state of the entity
	DefineState bool `form:"define_state" json:"define_state"`
}

type EditRequest struct {
	ID             string               `json:"-"`
	Name           string               `json:"name" binding:"required,max=255"`
	Enabled        *bool                `json:"enabled" binding:"required"`
	OutputTemplate string               `json:"output_template" binding:"required,max=500"`
	Category       string               `json:"category"`
	ImpactLevel    int64                `json:"impact_level" binding:"required,min=1,max=10"`
	Infos          []entity.InfoRequest `json:"infos" binding:"dive"`
	SliAvailState  *int64               `json:"sli_avail_state" binding:"required,min=0,max=3"`

	Coordinates *types.Coordinates `json:"coordinates"`

	common.EntityPatternFieldsRequest
}

type ContextGraphEntity struct {
	entity.Entity     `bson:",inline"`
	StateDependsCount *int                  `bson:"state_depends_count" json:"state_depends_count,omitempty"`
	StateSetting      *StateSettingResponse `bson:"state_setting" json:"state_setting,omitempty"`
}

type ContextGraphAggregationResult struct {
	Data       []ContextGraphEntity `bson:"data"`
	TotalCount int64                `bson:"total_count" json:"total_count"`
}

func (r *ContextGraphAggregationResult) GetData() interface{} {
	return r.Data
}

func (r *ContextGraphAggregationResult) GetTotal() int64 {
	return r.TotalCount
}

type StateSettingResponse struct {
	ID     string `bson:"_id" json:"_id"`
	Title  string `bson:"title" json:"title"`
	Method string `bson:"method" json:"method"`
}
