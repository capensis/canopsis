package statesettings

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
)

type EditRequest struct {
	ID string `json:"-"`
	StateSetting

	patternfields.EntityRequest
	InheritedEntityPatternRequest

	Author string `json:"author,omitempty" swaggerignore:"true"`
}

type Response struct {
	ID           string `json:"_id" bson:"_id"`
	StateSetting `bson:"inline"`

	savedpattern.EntityPatternFields          `bson:",inline"`
	statesetting.InheritedEntityPatternFields `bson:",inline"`

	Author    *author.Author    `json:"author,omitempty" bson:"author,omitempty"`
	Created   *datetime.CpsTime `json:"created,omitempty" bson:"created,omitempty" swaggertype:"integer"`
	Updated   *datetime.CpsTime `json:"updated,omitempty" bson:"updated,omitempty" swaggertype:"integer"`
	Editable  bool              `bson:"editable" json:"editable"`
	Deletable bool              `bson:"deletable" json:"deletable"`
}

type StateSetting struct {
	// Possible method values.
	//   * `worst` - take worst state.
	//   * `worst_of_share` - take worst of share defined in junit_thresholds.
	//   * `inherited` - take worst of subset of dependencies defined by pattern in inherited_entity_pattern.
	//   * `dependencies` - calculate state by rules defined in state_thresholds.
	Method string `json:"method" bson:"method" binding:"required"`

	// Service and component state setting only fields
	Title    *string `json:"title,omitempty" bson:"title,omitempty"`
	Enabled  *bool   `json:"enabled,omitempty" bson:"enabled,omitempty"`
	Priority int64   `json:"priority" bson:"priority" binding:"min=0"`

	StateThresholds *StateThresholds `json:"state_thresholds,omitempty" bson:"state_thresholds,omitempty"`
	Type            *string          `json:"type,omitempty" bson:"type,omitempty" binding:"required_if=Method inherited,required_if=Method dependencies,omitempty,oneof=component service"`

	// JUnit state setting only field
	JUnitThresholds *JUnitThresholds `json:"junit_thresholds,omitempty" bson:"junit_thresholds,omitempty"`
}

type StateThresholds struct {
	Critical *StateThreshold `json:"critical,omitempty" bson:"critical,omitempty"`
	Major    *StateThreshold `json:"major,omitempty" bson:"major,omitempty"`
	Minor    *StateThreshold `json:"minor,omitempty" bson:"minor,omitempty"`
	OK       *StateThreshold `json:"ok,omitempty" bson:"ok,omitempty"`
}

func (t *StateThresholds) ToModel() *statesetting.StateThresholds {
	if t == nil {
		return nil
	}

	return &statesetting.StateThresholds{
		Critical: t.Critical.ToModel(),
		Major:    t.Major.ToModel(),
		Minor:    t.Minor.ToModel(),
		OK:       t.OK.ToModel(),
	}
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

func (t *StateThreshold) ToModel() *statesetting.StateThreshold {
	if t == nil {
		return nil
	}

	return &statesetting.StateThreshold{
		Method: t.Method,
		State:  t.State,
		Cond:   t.Cond,
		Value:  t.Value,
	}
}

type JUnitThreshold struct {
	Minor    *float64 `json:"minor" bson:"minor" binding:"required,numeric,gte=0,lte=100,ltefield=Major,ltefield=Critical"`
	Major    *float64 `json:"major" bson:"major" binding:"required,numeric,gte=0,lte=100,ltefield=Critical"`
	Critical *float64 `json:"critical" bson:"critical" binding:"required,numeric,gte=0,lte=100"`
	Type     *int     `json:"type" bson:"type" binding:"required"`
}

func (t *JUnitThreshold) ToModel() *statesetting.JUnitThreshold {
	if t == nil {
		return nil
	}

	r := statesetting.JUnitThreshold{}
	if t.Minor != nil {
		r.Minor = *t.Minor
	}

	if t.Major != nil {
		r.Major = *t.Major
	}

	if t.Critical != nil {
		r.Critical = *t.Critical
	}

	if t.Type != nil {
		r.Type = *t.Type
	}

	return &r
}

type JUnitThresholds struct {
	Skipped  *JUnitThreshold `json:"skipped" bson:"skipped" binding:"required"`
	Errors   *JUnitThreshold `json:"errors" bson:"errors" binding:"required"`
	Failures *JUnitThreshold `json:"failures" bson:"failures" binding:"required"`
}

func (t *JUnitThresholds) ToModel() *statesetting.JUnitThresholds {
	if t == nil {
		return nil
	}

	r := statesetting.JUnitThresholds{}
	s := t.Skipped.ToModel()
	if s != nil {
		r.Skipped = *s
	}

	e := t.Errors.ToModel()
	if e != nil {
		r.Errors = *e
	}

	f := t.Failures.ToModel()
	if f != nil {
		r.Failures = *f
	}

	return &r
}

type InheritedEntityPatternRequest struct {
	InheritedEntityPattern          pattern.Entity `json:"inherited_entity_pattern" binding:"entity_pattern"`
	CorporateInheritedEntityPattern string         `json:"corporate_inherited_entity_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
}

func (r InheritedEntityPatternRequest) ToModelWithoutFields(collectionName string) statesetting.InheritedEntityPatternFields {
	if r.CorporatePattern.ID == "" {
		return statesetting.InheritedEntityPatternFields{
			InheritedEntityPattern: r.InheritedEntityPattern,
		}
	}

	forbiddenFields := patternfields.GetForbiddenFieldsInEntityPattern(collectionName)

	return statesetting.InheritedEntityPatternFields{
		InheritedEntityPattern:               r.CorporatePattern.EntityPattern.RemoveFields(forbiddenFields),
		CorporateInheritedEntityPattern:      r.CorporatePattern.ID,
		CorporateInheritedEntityPatternTitle: r.CorporatePattern.Title,
	}
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
