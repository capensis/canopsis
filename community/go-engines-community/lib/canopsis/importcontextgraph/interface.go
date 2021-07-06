package importcontextgraph

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	ActionDelete  = "delete"
	ActionCreate  = "create"
	ActionSet     = "set"
	ActionUpdate  = "update"
	ActionDisable = "disable"
	ActionEnable  = "enable"
)

var ErrNotImplemented = errors.New("import action not implemented")

type Worker interface {
	Work(ctx context.Context, filename, source string) (Stats, error)
}

type Stats struct {
	ExecTime time.Duration `bson:"-" json:"-"`
	Deleted  int64         `bson:"deleted" json:"deleted"`
	Updated  int64         `bson:"updated" json:"updated"`
}

type EventPublisher interface {
	SendUpdateEntityServiceEvent(serviceId string) error
}

type ConfigurationItem struct {
	ID               string                     `json:"_id" bson:"-"`
	Name             *string                    `json:"name" bson:"name,omitempty"`
	Depends          []string                   `json:"-" bson:"depends"`
	Impact           []string                   `json:"-" bson:"impact"`
	EnableHistory    []string                   `json:"-" bson:"enable_history"`
	Measurements     []interface{}              `json:"measurements" bson:"measurements"`
	EntityPatterns   *pattern.EntityPatternList `bson:"entity_patterns,omitempty" json:"entity_patterns"`
	OutputTemplate   *string                    `bson:"output_template,omitempty" json:"output_template"`
	Infos            map[string]interface{}     `json:"infos" bson:"infos"`
	Type             *string                    `json:"type" bson:"type,omitempty"`
	Category         *string                    `json:"category" bson:"category,omitempty"`
	ImpactLevel      *int64                     `json:"impact_level" bson:"impact_level,omitempty"`
	Enabled          *bool                      `json:"enabled" bson:"enabled,omitempty"`
	Action           string                     `json:"action" bson:"-"`
	ActionProperties interface{}                `json:"action_properties" bson:"-"`
	ImportSource     string                     `json:"-" bson:"import_source"`
	Imported         types.CpsTime              `json:"-" bson:"imported"`
}

type Link struct {
	ID               string                 `json:"_id"`
	To               string                 `json:"to"`
	From             []string               `json:"from"`
	Infos            map[string]interface{} `json:"infos"`
	Action           string                 `json:"action"`
	ActionProperties interface{}            `json:"action_properties" bson:"-"`
}
