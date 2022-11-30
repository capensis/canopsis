package types

// Entity types
const (
	EntityTypeConnector = "connector"
	EntityTypeComponent = "component"
	EntityTypeResource  = "resource"
	EntityTypeService   = "service"
	EntityTypeMetaAlarm = "metaalarm"
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
	ID                   string          `bson:"_id" json:"_id"`
	Name                 string          `bson:"name" json:"name"`
	Description          string          `bson:"description" json:"description"`
	EnableHistory        []CpsTime       `bson:"enable_history" json:"enable_history"`
	Measurements         interface{}     `bson:"measurements" json:"measurements"` // unused collection ids
	Enabled              bool            `bson:"enabled" json:"enabled"`
	Infos                map[string]Info `bson:"infos" json:"infos"`
	ComponentInfos       map[string]Info `bson:"component_infos,omitempty" json:"component_infos,omitempty"`
	Type                 string          `bson:"type" json:"type"`
	Category             string          `bson:"category" json:"category"`
	ImpactLevel          int64           `bson:"impact_level" json:"impact_level"`
	IsNew                bool            `bson:"-" json:"-"`
	AlarmsCumulativeData struct {
		// Only for Service.
		// WatchedCount is count of unresolved alarms.
		WatchedCount int64 `bson:"watched_count"`
		// WatchedPbehaviorCount contains counters of unresolved and in pbehavior alarms.
		WatchedPbehaviorCount map[string]int64 `bson:"watched_pbehavior_count"`
		// WatchedNotAckedCount is count of unresolved and not acked and active (by pbehavior) alarms.
		WatchedNotAckedCount int64 `bson:"watched_not_acked_count"`
	} `bson:"alarms_cumulative_data,omitempty" json:"-"`
	Created       CpsTime  `bson:"created" json:"created"`
	LastEventDate *CpsTime `bson:"last_event_date,omitempty" json:"last_event_date,omitempty"`

	Impacts []string `bson:"impact" json:"impact"`
	// impacted_services field is only for connectors, see entity service RecomputeIdleSince method
	ImpactedServices []string `bson:"impacted_services" json:"-"`
	Depends          []string `bson:"depends" json:"depends"`
	Connector        string   `bson:"connector,omitempty" json:"connector,omitempty"`
	Component        string   `bson:"component,omitempty" json:"component,omitempty"`

	// LastIdleRuleApply is used to mark entity if some idle rule was applied.
	LastIdleRuleApply string `bson:"last_idle_rule_apply,omitempty" json:"last_idle_rule_apply,omitempty"`
	// IdleSince represents since when entity didn't receive any events.
	IdleSince *CpsTime `bson:"idle_since,omitempty" json:"idle_since,omitempty"`

	ImportSource string   `bson:"import_source,omitempty" json:"import_source"`
	Imported     *CpsTime `bson:"imported,omitempty" json:"imported"`

	PbehaviorInfo PbehaviorInfo `bson:"pbehavior_info,omitempty" json:"pbehavior_info,omitempty"`

	SliAvailState int64 `bson:"sli_avail_state" json:"sli_avail_state"`
	// Coordinates is used only in api, add json tag if it's required in an event.
	Coordinates Coordinates `bson:"coordinates,omitempty" json:"-"`

	SoftDeleted bool `bson:"soft_deleted" json:"soft_deleted"`
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

func (e Entity) HasImpact(impact string) bool {
	for _, v := range e.Impacts {
		if v == impact {
			return true
		}
	}

	return false
}

func (e Entity) HasDepend(depend string) bool {
	for _, v := range e.Depends {
		if v == depend {
			return true
		}
	}

	return false
}
