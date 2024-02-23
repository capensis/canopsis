package statesettings

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

type EditRequest struct {
	ID           string `json:"-" bson:"_id"`
	StateSetting `bson:"inline"`
}

type Response struct {
	ID           string `json:"_id" bson:"_id"`
	StateSetting `bson:"inline"`

	Editable  bool `bson:"editable" json:"editable"`
	Deletable bool `bson:"deletable" json:"deletable"`
}

type StateSetting struct {
	// Possible method values.
	//   * `worst` - take worst state.
	//   * `worst_of_share` - take worst of share defined in junit_thresholds.
	//   * `inherited` - take worst of subset of dependencies defined by pattern in inherited_entity_pattern.
	//   * `dependencies` - calculate state by rules defined in state_thresholds.
	Method string `json:"method" bson:"method" binding:"required"`

	// Service and component state setting only fields
	Title                  *string          `json:"title,omitempty" bson:"title,omitempty"`
	Enabled                *bool            `json:"enabled,omitempty" bson:"enabled,omitempty"`
	Priority               int64            `json:"priority" bson:"priority" binding:"min=0"`
	EntityPattern          *pattern.Entity  `json:"entity_pattern,omitempty" bson:"entity_pattern,omitempty"`
	InheritedEntityPattern *pattern.Entity  `json:"inherited_entity_pattern,omitempty" bson:"inherited_entity_pattern,omitempty"`
	StateThresholds        *StateThresholds `json:"state_thresholds,omitempty" bson:"state_thresholds,omitempty"`
	Type                   *string          `json:"type,omitempty" bson:"type,omitempty" binding:"required_if=Method inherited,required_if=Method dependencies,omitempty,oneof=component service"`

	// JUnit state setting only field
	JUnitThresholds *JUnitThresholds `json:"junit_thresholds,omitempty" bson:"junit_thresholds,omitempty"`
}

type StateThresholds struct {
	Critical *StateThreshold `json:"critical,omitempty" bson:"critical,omitempty"`
	Major    *StateThreshold `json:"major,omitempty" bson:"major,omitempty"`
	Minor    *StateThreshold `json:"minor,omitempty" bson:"minor,omitempty"`
	OK       *StateThreshold `json:"ok,omitempty" bson:"ok,omitempty"`
}

type StateThreshold struct {
	// Possible method values.
	//   * `number` - calculate by number of entities.
	//   * `share` - calculate by share of entities.
	Method string `json:"method" bson:"method" binding:"oneof=number share"`
	// Possible state values.
	//   * `critical` - calculate critical state.
	//   * `major` - calculate major state.
	//   * `minor` - calculate minor state.
	//   * `ok` - calculate ok state.
	State string `json:"state" bson:"state" binding:"oneof=critical major minor ok"`
	// Possible cond values.
	//   * `gt` - greater than.
	//   * `lt` - less than.
	Cond  string `json:"cond" bson:"cond" binding:"oneof=gt lt"`
	Value int    `json:"value" bson:"value"`
}

type JUnitThreshold struct {
	Minor    *float64 `json:"minor" bson:"minor" binding:"required,numeric,gte=0,lte=100,ltefield=Major,ltefield=Critical"`
	Major    *float64 `json:"major" bson:"major" binding:"required,numeric,gte=0,lte=100,ltefield=Critical"`
	Critical *float64 `json:"critical" bson:"critical" binding:"required,numeric,gte=0,lte=100"`
	Type     *int     `json:"type" bson:"type" binding:"required"`
}

type JUnitThresholds struct {
	Skipped  JUnitThreshold `json:"skipped" bson:"skipped" binding:"required"`
	Errors   JUnitThreshold `json:"errors" bson:"errors" binding:"required"`
	Failures JUnitThreshold `json:"failures" bson:"failures" binding:"required"`
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

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=title enabled priority method"`
}
