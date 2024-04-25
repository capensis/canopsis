package types

import (
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

const EntityDefaultImpactLevel = 1

// Info contain extra values for the entity
type Info struct {
	Name        string      `bson:"name,omitempty" json:"name"`
	Description string      `bson:"description,omitempty" json:"description"`
	Value       interface{} `bson:"value,omitempty" json:"value"`
}

// Entity ...
type Entity struct {
	ID             string             `bson:"_id" json:"_id"`
	Name           string             `bson:"name" json:"name"`
	Description    string             `bson:"description" json:"description"`
	EnableHistory  []datetime.CpsTime `bson:"enable_history" json:"enable_history"`
	Measurements   interface{}        `bson:"measurements" json:"measurements"` // unused collection ids
	Enabled        bool               `bson:"enabled" json:"enabled"`
	Infos          map[string]Info    `bson:"infos" json:"infos"`
	ComponentInfos map[string]Info    `bson:"component_infos,omitempty" json:"component_infos,omitempty"`
	Type           string             `bson:"type" json:"type"`
	Category       string             `bson:"category" json:"category"`
	ImpactLevel    int64              `bson:"impact_level" json:"impact_level"`
	Created        datetime.CpsTime   `bson:"created" json:"created"`
	LastEventDate  *datetime.CpsTime  `bson:"last_event_date,omitempty" json:"last_event_date,omitempty"`

	Connector string `bson:"connector,omitempty" json:"connector,omitempty"`
	Component string `bson:"component,omitempty" json:"component,omitempty"`
	// ImpactedServices field is only for connectors, see entity service RecomputeIdleSince method.
	ImpactedServices []string `bson:"impacted_services" json:"-"`

	// LastIdleRuleApply is used to mark entity if some idle rule was applied.
	LastIdleRuleApply string `bson:"last_idle_rule_apply,omitempty" json:"last_idle_rule_apply,omitempty"`
	// IdleSince represents since when entity didn't receive any events.
	IdleSince *datetime.CpsTime `bson:"idle_since,omitempty" json:"idle_since,omitempty"`

	ImportSource string            `bson:"import_source,omitempty" json:"import_source"`
	Imported     *datetime.CpsTime `bson:"imported,omitempty" json:"imported"`

	PbehaviorInfo     PbehaviorInfo     `bson:"pbehavior_info,omitempty" json:"pbehavior_info,omitempty"`
	LastPbehaviorDate *datetime.CpsTime `bson:"last_pbehavior_date,omitempty" json:"last_pbehavior_date,omitempty"`

	SliAvailState int64 `bson:"sli_avail_state" json:"sli_avail_state"`

	Services         []string `bson:"services" json:"services,omitempty"`
	ServicesToAdd    []string `bson:"services_to_add,omitempty" json:"services_to_add,omitempty"`
	ServicesToRemove []string `bson:"services_to_remove,omitempty" json:"services_to_remove,omitempty"`

	// Coordinates is used only in api, add json tag if it's required in an event.
	Coordinates Coordinates `bson:"coordinates,omitempty" json:"-"`

	SoftDeleted *datetime.CpsTime `bson:"soft_deleted,omitempty" json:"soft_deleted,omitempty"`

	PerfData        []string          `bson:"perf_data,omitempty" json:"-"`
	PerfDataUpdated *datetime.CpsTime `bson:"perf_data_updated,omitempty" json:"-"`

	// IsNew and IsUpdated used in engine che in entity creation and eventfilter
	IsNew     bool `bson:"-" json:"-"`
	IsUpdated bool `bson:"-" json:"-"`

	Healthcheck bool `bson:"healthcheck,omitempty" json:"-"`

	StateInfo *StateInfo `bson:"state_info" json:"state_info"`

	ComponentStateSettings         bool `bson:"component_state_settings,omitempty" json:"component_state_settings,omitempty"`
	ComponentStateSettingsToAdd    bool `bson:"component_state_settings_to_add,omitempty" json:"component_state_settings_to_add,omitempty"`
	ComponentStateSettingsToRemove bool `bson:"component_state_settings_to_remove,omitempty" json:"component_state_settings_to_remove,omitempty"`

	Comments    []EntityComment `bson:"comments,omitempty" json:"comments,omitempty"`
	LastComment *EntityComment  `bson:"last_comment,omitempty" json:"last_comment,omitempty"`
}

type StateInfo struct {
	ID               string          `bson:"_id" json:"_id"`
	InheritedPattern *pattern.Entity `bson:"inherited_pattern,omitempty" json:"inherited_pattern,omitempty"`
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
