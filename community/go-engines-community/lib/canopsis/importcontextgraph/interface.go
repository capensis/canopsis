package importcontextgraph

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
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
	WorkPartial(ctx context.Context, filename, source string) (Stats, error)
}

type Stats struct {
	ExecTime time.Duration `bson:"-" json:"-"`
	Deleted  int64         `bson:"deleted" json:"deleted"`
	Updated  int64         `bson:"updated" json:"updated"`
}

type EventPublisher interface {
	SendEvent(ctx context.Context, event types.Event) error
}

type EntityConfiguration struct {
	ID             string                `json:"-" bson:"_id"`
	Name           string                `json:"name" bson:"name" binding:"required"`
	Component      string                `json:"component" bson:"component,omitempty"`
	Services       []string              `json:"-" bson:"services,omitempty"`
	EnableHistory  []int64               `json:"-" bson:"enable_history"`
	EntityPattern  pattern.Entity        `json:"entity_pattern" bson:"entity_pattern,omitempty"`
	OutputTemplate string                `json:"output_template" bson:"output_template,omitempty"`
	Infos          map[string]types.Info `json:"infos" bson:"infos"`
	Type           string                `json:"type" bson:"type" binding:"required"`
	CategoryName   string                `json:"category" bson:"-"`
	CategoryID     string                `json:"-" bson:"category,omitempty"`
	ImpactLevel    int64                 `json:"impact_level" bson:"impact_level,omitempty"`
	Enabled        bool                  `json:"enabled" bson:"enabled,omitempty"`
	Action         string                `json:"action" bson:"-"`
	ImportSource   string                `json:"-" bson:"import_source"`
	Imported       datetime.CpsTime      `json:"-" bson:"imported"`
}

type ConfigurationItem struct {
	ID             string                 `json:"_id" bson:"-"`
	Name           *string                `json:"name" bson:"name,omitempty"`
	Component      string                 `json:"-" bson:"component,omitempty"`
	Connector      string                 `json:"-" bson:"connector,omitempty"`
	Services       []string               `json:"-" bson:"services,omitempty"`
	EnableHistory  []int64                `json:"-" bson:"enable_history"`
	Measurements   []interface{}          `json:"measurements" bson:"measurements"`
	EntityPattern  pattern.Entity         `json:"entity_pattern,omitempty" bson:"entity_pattern"`
	OutputTemplate *string                `json:"output_template,omitempty" bson:"output_template"`
	Infos          map[string]interface{} `json:"infos" bson:"infos"`
	Type           *string                `json:"type" bson:"type,omitempty" binding:"oneof=connector component resource service"`
	Category       *string                `json:"category" bson:"category,omitempty"`
	ImpactLevel    *int64                 `json:"impact_level" bson:"impact_level,omitempty"`
	Enabled        *bool                  `json:"enabled" bson:"enabled,omitempty"`
	Action         string                 `json:"action" bson:"-" binding:"oneof=set create update delete enable disable"`
	ImportSource   string                 `json:"-" bson:"import_source"`
	Imported       datetime.CpsTime       `json:"-" bson:"imported"`
}

type Link struct {
	ID               string                 `json:"_id"`
	To               string                 `json:"to"`
	From             []string               `json:"from"`
	Infos            map[string]interface{} `json:"infos"`
	Action           string                 `json:"action" binding:"oneof=create delete"`
	ActionProperties interface{}            `json:"action_properties" bson:"-"`
}
