package entityservice

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/entity"
	"git.canopsis.net/canopsis/go-engines/lib/api/entitybasic"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type EntityService struct {
	entity.Entity  `bson:",inline"`
	EntityPatterns pattern.EntityPatternList `bson:"entity_patterns" json:"entity_patterns"`
	OutputTemplate string                    `bson:"output_template" json:"output_template"`
}

type AlarmWithEntity struct {
	Entity          entity.Entity `bson:"entity" json:"entity"`
	Alarm           *types.Alarm  `bson:"alarm" json:"alarm"`
	ImpactState     int64         `bson:"impact_state" json:"impact_state"`
	HasDependencies *bool         `bson:"has_dependencies,omitempty" json:"has_dependencies,omitempty"`
	HasImpacts      *bool         `bson:"has_impacts,omitempty" json:"has_impacts,omitempty"`
}

type CreateRequest struct {
	EditRequest
	ID string `json:"_id" binding:"id"`
}

type UpdateRequest struct {
	EditRequest
	ID string `json:"-"`
}

type ContextGraphRequest struct {
	pagination.Query
	ID string `form:"_id"`
}

type EditRequest struct {
	Name           string                    `json:"name" binding:"required,max=255"`
	Enabled        *bool                     `json:"enabled" binding:"required"`
	OutputTemplate string                    `json:"output_template" binding:"required,max=500"`
	Category       string                    `json:"category"`
	ImpactLevel    int64                     `json:"impact_level" binding:"required,min=1,max=10"`
	EntityPatterns pattern.EntityPatternList `json:"entity_patterns"`
	Infos          []entitybasic.InfoRequest `json:"infos" binding:"dive"`
}

type ContextGraphAggregationResult struct {
	Data       []AlarmWithEntity `bson:"data"`
	TotalCount int64             `bson:"total_count" json:"total_count"`
}

func (r *ContextGraphAggregationResult) GetData() interface{} {
	return r.Data
}

func (r *ContextGraphAggregationResult) GetTotal() int64 {
	return r.TotalCount
}
