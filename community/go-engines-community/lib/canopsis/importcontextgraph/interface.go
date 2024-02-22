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
	ActionSet     = "set"
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
