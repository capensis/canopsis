package importcontextgraph

import (
	"context"
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
	Name           string                `json:"name" binding:"required"`
	Component      string                `json:"component"`
	Upstream       *string               `json:"upstream"`
	EntityPattern  pattern.Entity        `json:"entity_pattern"`
	OutputTemplate string                `json:"output_template"`
	Infos          map[string]types.Info `json:"infos"`
	Type           string                `json:"type" binding:"required"`
	CategoryName   string                `json:"category"`
	ImpactLevel    int64                 `json:"impact_level"`
	Enabled        bool                  `json:"enabled"`
	Tags           map[string]string     `json:"tags"`
	Action         string                `json:"action"`

	// Aliases is used to ease find by entity info property api.
	Aliases []string `json:"-"`
}

type Entity struct {
	ID                string                `bson:"_id"`
	Name              string                `bson:"name"`
	Component         string                `bson:"component,omitempty"`
	Services          []string              `bson:"services,omitempty"`
	Upstream          *string               `bson:"upstream,omitempty"`
	IsUpstreamChanged bool                  `bson:"is_upstream_changed,omitempty"`
	EnableHistory     []int64               `bson:"enable_history"`
	EntityPattern     pattern.Entity        `bson:"entity_pattern,omitempty"`
	OutputTemplate    string                `bson:"output_template,omitempty"`
	Infos             map[string]types.Info `bson:"infos"`
	Type              string                `bson:"type"`
	CategoryID        string                `bson:"category,omitempty"`
	ImpactLevel       int64                 `bson:"impact_level,omitempty"`
	Enabled           bool                  `bson:"enabled,omitempty"`
	ImportTags        []string              `bson:"imtags"`
	ImportSource      string                `bson:"import_source"`
	Imported          datetime.CpsTime      `bson:"imported"`

	// Aliases is used to ease find by entity info property api.
	Aliases []string `bson:"aliases"`
}
