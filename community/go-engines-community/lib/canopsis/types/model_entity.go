package types

// Entity types
const (
	EntityTypeConnector = "connector"
	EntityTypeComponent = "component"
	EntityTypeResource  = "resource"
	EntityTypeService   = "service"
)

const EntityDefaultImpactLevel = 1

// Measurement ...
/*
type Measurement struct {
	ID   string   `bson:"_id" json:"_id"`
	tags []string `bson:"tags" json:"tags"`
}
*/

// Info contain extra values for the entity
type Info struct {
	Name        string      `bson:"name,omitempty" json:"name"`
	Description string      `bson:"description,omitempty" json:"description"`
	Value       interface{} `bson:"value,omitempty" json:"value"`
}

// Entity ...
type Entity struct {
	ID             string          `bson:"_id" json:"_id"`
	Name           string          `bson:"name" json:"name"`
	Description    string          `bson:"description" json:"description"`
	EnableHistory  []CpsTime       `bson:"enable_history" json:"enable_history"`
	Measurements   interface{}     `bson:"measurements" json:"measurements"` // unused collection ids
	Enabled        bool            `bson:"enabled" json:"enabled"`
	Infos          map[string]Info `bson:"infos" json:"infos"`
	ComponentInfos map[string]Info `bson:"component_infos,omitempty" json:"component_infos,omitempty"`
	Type           string          `bson:"type" json:"type"`
	Category       string          `bson:"category" json:"category"`
	ImpactLevel    int64           `bson:"impact_level" json:"impact_level"`
	IsNew          bool            `bson:"-" json:"-"`
	Created        CpsTime         `bson:"created" json:"created"`
	LastEventDate  *CpsTime        `bson:"last_event_date,omitempty" json:"last_event_date,omitempty"`

	Connector string   `bson:"connector,omitempty" json:"connector,omitempty"`
	Component string   `bson:"component,omitempty" json:"component,omitempty"`
	Services  []string `bson:"services," json:"services,omitempty"`
	// ImpactedServices field is only for connectors, see entity service RecomputeIdleSince method.
	ImpactedServices []string `bson:"impacted_services" json:"-"`

	// LastIdleRuleApply is used to mark entity if some idle rule was applied.
	LastIdleRuleApply string `bson:"last_idle_rule_apply,omitempty" json:"last_idle_rule_apply,omitempty"`
	// IdleSince represents since when entity didn't receive any events.
	IdleSince *CpsTime `bson:"idle_since,omitempty" json:"idle_since,omitempty"`

	ImportSource string   `bson:"import_source,omitempty" json:"import_source"`
	Imported     *CpsTime `bson:"imported,omitempty" json:"imported"`

	PbehaviorInfo     PbehaviorInfo `bson:"pbehavior_info,omitempty" json:"pbehavior_info,omitempty"`
	LastPbehaviorDate *CpsTime      `bson:"last_pbehavior_date,omitempty" json:"last_pbehavior_date,omitempty"`

	SliAvailState int64 `bson:"sli_avail_state" json:"sli_avail_state"`
	// Coordinates is used only in api, add json tag if it's required in an event.
	Coordinates Coordinates `bson:"coordinates,omitempty" json:"-"`

	// SoftDeleted is used to recompute service counters after dependency delete.
	SoftDeleted *CpsTime `bson:"soft_deleted,omitempty" json:"soft_deleted,omitempty"`
	// ResolveDeletedEventProcessed is set after processing service counters recalculation event on soft delete.
	ResolveDeletedEventProcessed *CpsTime `bson:"resolve_deleted_event_processed,omitempty" json:"resolve_deleted_event_processed,omitempty"`
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
		e.EnableHistory = make([]CpsTime, 0)
	}
}

// CacheID implements cache.Cache interface
func (e Entity) CacheID() string {
	return e.ID
}
