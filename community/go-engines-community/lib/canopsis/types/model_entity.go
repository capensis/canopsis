package types

import (
	"fmt"
	"slices"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

// Entity types
const (
	EntityTypeConnector = "connector"
	EntityTypeComponent = "component"
	EntityTypeResource  = "resource"
	EntityTypeService   = "service"
)

// Entity info types
const (
	EntityInfoTypeBoolean = iota
	EntityInfoTypeNumber
	EntityInfoTypeTimestamp
	EntityInfoTypeString
	EntityInfoTypeStringArray
)

const EntityDefaultImpactLevel = 1

// Info contain extra values for the entity
type Info struct {
	Name        string `bson:"name,omitempty" json:"name"`
	Description string `bson:"description,omitempty" json:"description"`
	Value       any    `bson:"value,omitempty" json:"value"`
}

// Entity ...
type Entity struct {
	ID             string             `bson:"_id" json:"_id"`
	Name           string             `bson:"name" json:"name"`
	Author         string             `bson:"author,omitempty" json:"author,omitempty"`
	Description    string             `bson:"description" json:"description"`
	EnableHistory  []datetime.CpsTime `bson:"enable_history" json:"enable_history"`
	Measurements   any                `bson:"measurements" json:"measurements"` // unused collection ids
	Enabled        bool               `bson:"enabled" json:"enabled"`
	Infos          map[string]Info    `bson:"infos" json:"infos"`
	ComponentInfos map[string]Info    `bson:"component_infos,omitempty" json:"component_infos,omitempty"`
	Type           string             `bson:"type" json:"type"`
	Category       string             `bson:"category" json:"category"`
	ImpactLevel    int64              `bson:"impact_level" json:"impact_level"`
	Created        datetime.CpsTime   `bson:"created" json:"created"`
	Updated        *datetime.CpsTime  `bson:"updated,omitempty" json:"updated,omitempty"`
	LastEventDate  *datetime.CpsTime  `bson:"last_event_date,omitempty" json:"last_event_date,omitempty"`

	Connector string `bson:"connector,omitempty" json:"connector,omitempty"`
	Component string `bson:"component,omitempty" json:"component,omitempty"`
	Upstream  string `bson:"upstream,omitempty" json:"upstream,omitempty"`

	Services          []string `bson:"services" json:"services,omitempty"`
	InheritedServices []string `bson:"inherited_services,omitempty" json:"inherited_services,omitempty"`
	// ServicesToAdd and ServicesToRemove are used to recompute dependencies counters.
	ServicesToAdd    []string `bson:"services_to_add,omitempty" json:"services_to_add,omitempty"`
	ServicesToRemove []string `bson:"services_to_remove,omitempty" json:"services_to_remove,omitempty"`
	// IsUpstreamChanged is used to update corresponding alarm status.
	IsUpstreamChanged bool `bson:"is_upstream_changed,omitempty" json:"is_upstream_changed,omitempty"`

	// ImpactedServices field is only for connectors, see entity service RecomputeIdleSince method.
	ImpactedServices []string `bson:"impacted_services" json:"-"`

	// LastIdleRuleApply is used to mark entity if some idle rule was applied.
	LastIdleRuleApply string `bson:"last_idle_rule_apply,omitempty" json:"last_idle_rule_apply,omitempty"`
	// IdleSince represents since when entity didn't receive any events.
	IdleSince *datetime.CpsTime `bson:"idle_since,omitempty" json:"idle_since,omitempty"`

	ImportSource string            `bson:"import_source,omitempty" json:"import_source"`
	Imported     *datetime.CpsTime `bson:"imported,omitempty" json:"imported"`
	ImportTags   []string          `bson:"imtags,omitempty" json:"imtags,omitempty"`

	PbehaviorInfo     PbehaviorInfo     `bson:"pbehavior_info,omitempty" json:"pbehavior_info"`
	LastPbehaviorDate *datetime.CpsTime `bson:"last_pbehavior_date,omitempty" json:"last_pbehavior_date,omitempty"`

	SliAvailState int64 `bson:"sli_avail_state" json:"sli_avail_state"`

	// Coordinates is used only in api, add json tag if it's required in an event.
	Coordinates Coordinates `bson:"coordinates,omitempty" json:"-"`

	// SoftDeleted is used to recompute service counters after dependency delete.
	SoftDeleted *datetime.CpsTime `bson:"soft_deleted,omitempty" json:"soft_deleted,omitempty"`
	// ResolveDeletedEventProcessed is set after processing service counters recalculation event on soft delete.
	ResolveDeletedEventProcessed *datetime.CpsTime `bson:"resolve_deleted_event_processed,omitempty" json:"resolve_deleted_event_processed,omitempty"`

	PerfData        []string          `bson:"perf_data,omitempty" json:"-"`
	PerfDataUpdated *datetime.CpsTime `bson:"perf_data_updated,omitempty" json:"-"`

	Healthcheck bool `bson:"healthcheck,omitempty" json:"-"`

	StateInfo *StateInfo `bson:"state_info" json:"state_info"`

	ComponentStateSettings         bool `bson:"component_state_settings,omitempty" json:"component_state_settings,omitempty"`
	ComponentStateSettingsToAdd    bool `bson:"component_state_settings_to_add,omitempty" json:"component_state_settings_to_add,omitempty"`
	ComponentStateSettingsToRemove bool `bson:"component_state_settings_to_remove,omitempty" json:"component_state_settings_to_remove,omitempty"`

	Comments    []EntityComment `bson:"comments,omitempty" json:"comments,omitempty"`
	LastComment *EntityComment  `bson:"last_comment,omitempty" json:"last_comment,omitempty"`

	// Pause is the active pause, nil if the entity is not paused.
	Pause *Pause `bson:"pause,omitempty" json:"pause,omitempty"`
	// LatestPause is the newest pause received, active or ended (End is set).
	// It is used to order pause/resume events and reject late ones, not to check if the entity is paused.
	LatestPause *LatestPause `bson:"latest_pause,omitempty" json:"latest_pause,omitempty"`
}

// MaxPausePathLength is the maximum number of the entities which a pause can be propagated through.
// A real hierarchy never reaches it, so the pause is not propagated further to stop a propagation loop.
const MaxPausePathLength = 100

// LatestPause is a pause with its end, Source is always set (entity ID for its own pause).
type LatestPause struct {
	Pause `bson:",inline"`
	End   *datetime.CpsTime `bson:"end,omitempty" json:"end,omitempty"`
}

// CanApplyOver reports whether the pause p can be applied to the entity whose latest pause is o.
// A pause inherited from another entity replaces the entity's own pause, whatever their times.
// A pause of an entity above the source of o, i.e. which has passed through it, replaces it,
// and a pause from below it is ignored, so all the dependencies of a paused entity keep its pause.
// Two other inherited pauses are compared by their times.
// The entity's own pause is applied only if the entity is not paused and its latest pause has ended before it.
func (p LatestPause) CanApplyOver(o *LatestPause, entityID string) bool {
	switch {
	case o == nil:
		return true
	case p.Source == entityID:
		return o.End != nil && p.Timestamp.After(*o.End)
	case o.Source == entityID:
		return true
	case slices.Contains(p.Path, o.Source):
		return true
	case slices.Contains(o.Path, p.Source):
		return false
	case !p.Timestamp.Equal(o.Timestamp):
		return p.Timestamp.After(o.Timestamp)
	default:
		// Source breaks ties so all entities pick the same pause.
		return p.Source > o.Source
	}
}

// InheritedPause returns the pause which the dependencies of the entity inherit from it,
// nil if the entity is not paused or if its pause is applied only to itself.
// The own pause of the entity becomes a pause with the entity as a source, and the entity is added
// to the pause path, so the dependencies know that the pause has passed through it.
func (e Entity) InheritedPause() (*LatestPause, error) {
	if e.Pause == nil || !e.Pause.Inherited {
		return nil, nil
	}

	pause := LatestPause{Pause: *e.Pause}
	if pause.Source != "" && pause.Source != e.ID {
		if len(pause.Path) >= MaxPausePathLength {
			return nil, fmt.Errorf("the path of the pause of %s is longer than %d entities", e.ID, MaxPausePathLength)
		}

		// The pause of another entity passes through this one.
		pause.Path = append(slices.Clone(pause.Path), e.ID)

		return &pause, nil
	}

	// The entity is the source of its own pause, so it's not in the path.
	pause.Source = e.ID
	pause.SourceType = e.Type
	switch e.Type {
	case EntityTypeComponent:
		pause.Message += "\n" + OutputComponentPrefix + e.Name
	case EntityTypeService:
		pause.Message += "\n" + OutputServicePrefix + e.Name
	}

	return &pause, nil
}

type StateInfo struct {
	ID               string         `bson:"_id" json:"_id"`
	InheritedPattern pattern.Entity `bson:"inherited_pattern,omitempty" json:"inherited_pattern,omitempty"`
}

type EntityComment struct {
	ID        string           `bson:"_id" json:"_id"`
	Timestamp datetime.CpsTime `bson:"t" json:"t"`
	Author    *Author          `bson:"a" json:"a"`
	Message   string           `bson:"m" json:"m"`
}

// Author in contrary to the author.Author struct without Name field
type Author struct {
	ID          string `bson:"_id" json:"_id"`
	DisplayName string `bson:"display_name" json:"display_name"`
}

type Coordinates struct {
	Lat float64 `bson:"lat" json:"lat" binding:"required,latitude"`
	Lng float64 `bson:"lng" json:"lng" binding:"required,longitude"`
}

func (c Coordinates) IsZero() bool {
	return c == Coordinates{}
}

type Pause struct {
	Timestamp  datetime.CpsTime `bson:"timestamp" json:"timestamp" swaggertype:"integer"`
	Message    string           `bson:"message" json:"message"`
	Author     string           `bson:"author" json:"author"`
	UserID     string           `bson:"user_id" json:"user_id"`
	Initiator  string           `bson:"initiator" json:"initiator"`
	Type       string           `bson:"type,omitempty" json:"type"`
	Source     string           `bson:"source,omitempty" json:"source"`
	SourceType string           `bson:"source_type,omitempty" json:"source_type"`
	Inherited  bool             `bson:"inherited,omitempty" json:"inherited,omitempty"`
	// Path is the entities the pause has passed through between its source and the paused entity,
	// its own source is not in it. It's used to compare a pause with a pause of an entity
	// above or below its source.
	Path []string `bson:"path,omitempty" json:"path,omitempty"`
}

// EnsureInitialized verifies that all complex structs are well initialized
func (e *Entity) EnsureInitialized() {
	if e.Infos == nil {
		e.Infos = make(map[string]Info)
	}
	if e.EnableHistory == nil {
		e.EnableHistory = make([]datetime.CpsTime, 0)
	}
}

// GetStringField is a magic getter for string fields for easier field retrieving when matching entity pattern
func (e *Entity) GetStringField(f string) (string, bool) {
	switch f {
	case "_id":
		return e.ID, true
	case "name":
		return e.Name, true
	case "category":
		return e.Category, true
	case "type":
		return e.Type, true
	case "connector":
		return e.Connector, true
	case "component":
		return e.Component, true
	case "pause.author":
		if e.Pause == nil {
			return "", false
		}

		return e.Pause.Author, true
	case "pause.initiator":
		if e.Pause == nil {
			return "", false
		}

		return e.Pause.Initiator, true
	default:
		return "", false
	}
}

// GetIntField is a magic getter for int fields for easier field retrieving when matching entity pattern
func (e *Entity) GetIntField(f string) (int64, bool) {
	switch f {
	case "impact_level":
		return e.ImpactLevel, true
	default:
		return 0, false
	}
}

// GetTimeField is a magic getter for time fields for easier field retrieving when matching entity pattern
func (e *Entity) GetTimeField(f string) (time.Time, bool) {
	switch f {
	case "last_event_date":
		if e.LastEventDate != nil {
			return e.LastEventDate.Time, true
		}

		return time.Time{}, true
	default:
		return time.Time{}, false
	}
}

// GetInfoVal is a magic getter for infos fields for easier field retrieving when matching entity pattern
func (e *Entity) GetInfoVal(f string) (any, bool) {
	if v, ok := e.Infos[f]; ok {
		return v.Value, true
	}

	return nil, false
}

// GetRefField is a magic getter for reference fields for easier field retrieving when matching entity pattern
func (e *Entity) GetRefField(f string) (any, bool) {
	switch f {
	case "pause":
		return e.Pause, true
	default:
		return nil, false
	}
}
