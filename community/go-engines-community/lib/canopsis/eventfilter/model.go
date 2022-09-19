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

	//are used in api external data
	RequestParameters *request.Parameters `bson:"request,omitempty" json:"request,omitempty"`
}

type Rule struct {
	ID           string                            `bson:"_id" json:"_id" binding:"id"`
	Author       string                            `bson:"author" json:"author" swaggerignore:"true"`
	Description  string                            `bson:"description" json:"description" binding:"required,max=255"`
	Type         string                            `bson:"type" json:"type" binding:"required,oneof=break drop enrichment change_entity"`
	Priority     int                               `bson:"priority" json:"priority"`
	Enabled      bool                              `bson:"enabled" json:"enabled"`
	OldPatterns  oldpattern.EventPatternList       `bson:"old_patterns,omitempty" json:"old_patterns,omitempty"`
	Config       RuleConfig                        `bson:"config" json:"config"`
	ExternalData map[string]ExternalDataParameters `bson:"external_data" json:"external_data,omitempty"`
	Created      *types.CpsTime                    `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated      *types.CpsTime                    `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`

	EventPattern                     pattern.Event `json:"event_pattern" bson:"event_pattern"`
	savedpattern.EntityPatternFields `bson:",inline"`
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
