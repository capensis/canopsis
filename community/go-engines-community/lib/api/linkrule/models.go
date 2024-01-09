package linkrule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id name author.name created updated"`
}

type EditRequest struct {
	ID         string            `json:"-"`
	Name       string            `json:"name" binding:"required,max=255"`
	Type       string            `json:"type" binding:"required,oneof=alarm entity"`
	Enabled    *bool             `json:"enabled" binding:"required"`
	Links      []link.Parameters `json:"links" binding:"dive"`
	SourceCode string            `json:"source_code"`
	Author     string            `json:"author" swaggerignore:"true"`

	ExternalData map[string]link.ExternalDataParameters `bson:"external_data" json:"external_data" binding:"dive"`

	common.AlarmPatternFieldsRequest
	common.EntityPatternFieldsRequest
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}

type Response struct {
	ID         string            `bson:"_id" json:"_id"`
	Name       string            `bson:"name" json:"name"`
	Type       string            `bson:"type" json:"type"`
	Enabled    bool              `bson:"enabled" json:"enabled"`
	Links      []link.Parameters `bson:"links" json:"links,omitempty"`
	SourceCode string            `bson:"source_code" json:"source_code,omitempty"`
	Author     *author.Author    `bson:"author" json:"author"`
	Created    datetime.CpsTime  `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated    datetime.CpsTime  `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`

	ExternalData map[string]link.ExternalDataParameters `bson:"external_data" json:"external_data"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
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

type CategoriesRequest struct {
	Limit  int64  `form:"limit" binding:"numeric,gte=0"`
	Type   string `form:"type" binding:"oneoforempty=alarm entity"`
	Search string `form:"search"`
}

type CategoryResponse struct {
	Categories []string `json:"categories"`
}
