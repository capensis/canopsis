package neweventfilter

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	RuleTypeChangeEntity = "change_entity"
)

type Rule struct {
	ID           string                            `bson:"_id"`
	Description  string                            `bson:"description"`
	Type         string                            `bson:"type"`
	Patterns     pattern.EventPatternList          `bson:"patterns"`
	Priority     int                               `bson:"priority"`
	Enabled      bool                              `bson:"enabled"`
	Config       RuleConfig                        `bson:"config"`
	Created      *types.CpsTime                    `bson:"created"`
	Updated      *types.CpsTime                    `bson:"updated"`
	Author       string                            `bson:"author"`
	ExternalData map[string]eventfilter.DataSource `bson:"external_data" json:"external_data"`
}

type RuleConfig struct {
	Resource      string `bson:"resource,omitempty" json:"resource,omitempty"`
	Component     string `bson:"component,omitempty" json:"component,omitempty"`
	Connector     string `bson:"connector,omitempty" json:"connector,omitempty"`
	ConnectorName string `bson:"connector_name,omitempty" json:"connector_name,omitempty"`
}

type ApplicatorParameters struct {
	Event        types.Event
	RegexMatch   pattern.EventRegexMatches
	ExternalData map[string]interface{}
}
