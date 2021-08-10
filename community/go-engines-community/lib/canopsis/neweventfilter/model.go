package neweventfilter

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
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

type Rule struct {
	ID          string                   `bson:"_id"`
	Description string                   `bson:"description"`
	Type        string                   `bson:"type"`
	Patterns    pattern.EventPatternList `bson:"patterns"`
	Priority    int                      `bson:"priority"`
	Enabled     bool                     `bson:"enabled"`
	Config      RuleConfig               `bson:"config"`
	Created     *types.CpsTime           `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated     *types.CpsTime           `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
	Author      string                   `bson:"author"`

	//TODO: copy from eventfilter package, all mongo plugin feature should be refactored
	ExternalData map[string]eventfilter.DataSource `bson:"external_data" json:"external_data"`
}

type RuleConfig struct {
	Resource      string   `bson:"resource,omitempty" json:"resource,omitempty"`
	Component     string   `bson:"component,omitempty" json:"component,omitempty"`
	Connector     string   `bson:"connector,omitempty" json:"connector,omitempty"`
	ConnectorName string   `bson:"connector_name,omitempty" json:"connector_name,omitempty"`
	Actions       []Action `bson:"actions,omitempty" json:"actions,omitempty" binding:"required_if=Type enrichment"`
	OnSuccess     string   `bson:"on_success,omitempty" json:"on_success,omitempty" binding:"required_if=Type enrichment"`
	OnFailure     string   `bson:"on_failure,omitempty" json:"on_failure,omitempty" binding:"required_if=Type enrichment"`
}

type Action struct {
	Type        string      `bson:"type" json:"type"`
	Name        string      `bson:"name" json:"name"`
	Description string      `bson:"description" json:"description"`
	Value       interface{} `bson:"value" json:"value"`
}

type TemplateParameters struct {
	Event        types.Event
	RegexMatch   pattern.EventRegexMatches
	ExternalData map[string]interface{}
}
