package eventfilter

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	RuleTypeChangeEntity = "change_entity"
	RuleTypeBreak        = "break"
	RuleTypeDrop         = "drop"
	RuleTypeEnrichment   = "enrichment"
)

const (
	ActionSetField = "set_field"

	// ActionSetFieldFromTemplate is a type of action that sets a string field of an
	// event using a template.
	ActionSetFieldFromTemplate = "set_field_from_template"

	// ActionSetEntityInfoFromTemplate is a type of action that sets an information
	// of an entity using a template.
	ActionSetEntityInfoFromTemplate = "set_entity_info_from_template"

	// ActionSetEntityInfoFromDictionary is a type of action that sets an information
	// of an entity using a dictionary.
	ActionSetEntityInfoFromDictionary = "set_entity_info_from_dictionary"

	// ActionSetEntityInfo is a type of action that sets an information
	// of an entity using a constant.
	ActionSetEntityInfo = "set_entity_info"

	// ActionCopyToEntityInfo is a type of action that copies a value from a field to
	// an information of an entity.
	ActionCopyToEntityInfo = "copy_to_entity_info"

	// ActionCopy is a type of action that copies a value from a field to another.
	ActionCopy = "copy"
)

type ExternalDataParameters struct {
	Type string `json:"type" bson:"type"`

	// are used in mongo external data
	Collection string            `json:"collection,omitempty" bson:"collection,omitempty"`
	Select     map[string]string `json:"select,omitempty" bson:"select,omitempty"`
	Regexp     map[string]string `json:"regexp,omitempty" bson:"regexp,omitempty"`
	SortBy     string            `json:"sort_by,omitempty" bson:"sort_by,omitempty"`
	Sort       string            `json:"sort,omitempty" bson:"sort,omitempty" binding:"oneoforempty=asc desc"`

	// are used in api external data
	RequestParameters *request.Parameters `bson:"request,omitempty" json:"request,omitempty"`
}

type Rule struct {
	ID           string                            `bson:"_id" json:"_id" binding:"id"`
	Author       string                            `bson:"author" json:"author" swaggerignore:"true"`
	Description  string                            `bson:"description" json:"description" binding:"required,max=255"`
	Type         string                            `bson:"type" json:"type" binding:"required,oneof=break drop enrichment change_entity"`
	Priority     int64                             `bson:"priority" json:"priority"`
	Enabled      bool                              `bson:"enabled" json:"enabled"`
	Config       RuleConfig                        `bson:"config" json:"config"`
	ExternalData map[string]ExternalDataParameters `bson:"external_data" json:"external_data,omitempty"`
	Created      *types.CpsTime                    `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated      *types.CpsTime                    `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
	EventsCount  int64                             `bson:"events_count,omitempty" json:"events_count,omitempty"`

	OldPatterns                      oldpattern.EventPatternList `bson:"old_patterns,omitempty" json:"old_patterns,omitempty"`
	EventPattern                     pattern.Event               `json:"event_pattern" bson:"event_pattern"`
	savedpattern.EntityPatternFields `bson:",inline"`

	RRule string         `json:"rrule" bson:"rrule"`
	Start *types.CpsTime `json:"start,omitempty" bson:"start,omitempty"`
	Stop  *types.CpsTime `json:"stop,omitempty" bson:"stop,omitempty"`

	//ResolvedStart and ResolvedStop shows the current or the next time interval, where eventfilter rule is enabled
	ResolvedStart *types.CpsTime `json:"-" bson:"resolved_start,omitempty"`
	ResolvedStop  *types.CpsTime `json:"-" bson:"resolved_stop,omitempty"`

	//NextResolvedStart and NextResolvedStop shows the next time interval after the one which is defined by ResolvedStart and ResolvedStop
	NextResolvedStart *types.CpsTime `json:"-" bson:"next_resolved_start,omitempty"`
	NextResolvedStop  *types.CpsTime `json:"-" bson:"next_resolved_stop,omitempty"`

	Exdates    []types.Exdate `json:"exdates" bson:"exdates"`
	Exceptions []string       `json:"exceptions" bson:"exceptions"`

	// ResolvedExdates shows exdates if their interval intersects with [now(); now() + 2 che periodical processes] interval
	ResolvedExdates []types.Exdate `json:"-" bson:"resolved_exdates"`
}

type RuleConfig struct {
	Resource      string `bson:"resource,omitempty" json:"resource,omitempty"`
	Component     string `bson:"component,omitempty" json:"component,omitempty"`
	Connector     string `bson:"connector,omitempty" json:"connector,omitempty"`
	ConnectorName string `bson:"connector_name,omitempty" json:"connector_name,omitempty"`

	// enrichment fields
	Actions   []Action `bson:"actions,omitempty" json:"actions,omitempty" binding:"dive,required_if=Type enrichment"`
	OnSuccess string   `bson:"on_success,omitempty" json:"on_success,omitempty"`
	OnFailure string   `bson:"on_failure,omitempty" json:"on_failure,omitempty"`
}

type Action struct {
	Type        string      `bson:"type" json:"type"`
	Name        string      `bson:"name" json:"name"`
	Description string      `bson:"description,omitempty" json:"description,omitempty"`
	Value       interface{} `bson:"value" json:"value" binding:"info_value"`
}

type RegexMatchWrapper struct {
	BackwardCompatibility bool
	OldRegexMatch         oldpattern.EventRegexMatches
	RegexMatch            RegexMatch
}

type RegexMatch struct {
	pattern.EventRegexMatches
	Entity pattern.EntityRegexMatches
}

type Template struct {
	Event             types.Event
	RegexMatchWrapper RegexMatchWrapper
	ExternalData      map[string]interface{}
}

func (t Template) GetTemplate() interface{} {
	if t.RegexMatchWrapper.BackwardCompatibility {
		return struct {
			Event        types.Event
			RegexMatch   oldpattern.EventRegexMatches
			ExternalData map[string]interface{}
		}{
			Event:        t.Event,
			RegexMatch:   t.RegexMatchWrapper.OldRegexMatch,
			ExternalData: t.ExternalData,
		}
	}

	return struct {
		Event        types.Event
		RegexMatch   RegexMatch
		ExternalData map[string]interface{}
	}{
		Event:        t.Event,
		RegexMatch:   t.RegexMatchWrapper.RegexMatch,
		ExternalData: t.ExternalData,
	}
}
