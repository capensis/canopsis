package entity

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitycategory"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const DefaultCategory = "default"

type ListRequestWithPagination struct {
	pagination.Query
	ListRequest
	WithFlags bool `form:"with_flags"`
	NoEvents  bool `form:"no_events"`
}

type ListRequest struct {
	Search   string   `form:"search" json:"search"`
	Filter   string   `form:"filter" json:"filter"`
	Sort     string   `form:"sort" json:"sort" binding:"oneoforempty=asc desc"`
	SortBy   string   `form:"sort_by" json:"sort_by" binding:"oneoforempty=name type category impact_level category.name idle_since"`
	SearchBy []string `form:"active_columns[]" json:"active_columns[]"`
	Category string   `form:"category" json:"category"`
}

type ExportRequest struct {
	ListRequest
	Separator string `form:"separator" json:"separator" binding:"oneoforempty=comma semicolon tab space"`
}

type CleanRequest struct {
	Archive             *bool `json:"archive" bson:"archive" binding:"required"`
	ArchiveDependencies bool  `json:"archive_dependencies" bson:"archive_dependencies"`
}

type CleanTask struct {
	Archive             *bool
	ArchiveDependencies bool
	UserID              string
}

type ExportResponse struct {
	ID     string `json:"_id"`
	Status int    `json:"status"`
}

type Entity struct {
	ID             string                   `bson:"_id" json:"_id"`
	Name           string                   `bson:"name" json:"name"`
	Impacts        []string                 `bson:"impact" json:"impact"`
	Depends        []string                 `bson:"depends" json:"depends"`
	EnableHistory  []types.CpsTime          `bson:"enable_history" json:"enable_history" swaggertype:"array,integer"`
	Measurements   interface{}              `bson:"measurements" json:"measurements"`
	Enabled        bool                     `bson:"enabled" json:"enabled"`
	Infos          Infos                    `bson:"infos" json:"infos"`
	ComponentInfos Infos                    `bson:"component_infos,omitempty" json:"component_infos,omitempty"`
	Type           string                   `bson:"type" json:"type"`
	Component      string                   `bson:"component,omitempty" json:"component,omitempty"`
	ImpactLevel    int64                    `bson:"impact_level" json:"impact_level"`
	Category       *entitycategory.Category `bson:"category" json:"category"`
	Deletable      *bool                    `bson:"deletable,omitempty" json:"deletable,omitempty"`
	IdleSince      *types.CpsTime           `bson:"idle_since,omitempty" json:"idle_since,omitempty" swaggertype:"integer"`
}

type Infos map[string]Info

type Info struct {
	Name        string      `bson:"name" json:"name"`
	Description string      `bson:"description" json:"description"`
	Value       interface{} `bson:"value" json:"value"`
}

func (i *Infos) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	var tmp map[string]Info
	err := bson.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	for k, info := range tmp {
		if info.Name == "" {
			info.Name = k
			tmp[k] = info
		}
	}

	*i = tmp

	return nil
}

type AggregationResult struct {
	Data       []Entity `bson:"data" json:"data"`
	TotalCount int64    `bson:"total_count" json:"total_count"`
}

func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

func (r AggregationResult) GetData() interface{} {
	return r.Data
}
