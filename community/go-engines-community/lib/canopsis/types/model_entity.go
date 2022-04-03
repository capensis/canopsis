package types

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

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
	ID          string   `bson:"_id" json:"_id"`
	Name        string   `bson:"name" json:"name"`
	Description string   `bson:"description" json:"description"`
	Impacts     []string `bson:"impact" json:"impact"`
	// impacted_services field is only for connectors, see entity service RecomputeIdleSince method
	ImpactedServices         []string        `bson:"impacted_services,omitempty" json:"impacted_services,omitempty"`
	ImpactedServicesToAdd    []string        `bson:"impacted_services_to_add,omitempty" json:"impacted_services_to_add"`
	ImpactedServicesToRemove []string        `bson:"impacted_services_to_remove,omitempty" json:"impacted_services_to_remove"`
	Depends                  []string        `bson:"depends,omitempty" json:"depends"`
	EnableHistory            []CpsTime       `bson:"enable_history" json:"enable_history"`
	Measurements             interface{}     `bson:"measurements" json:"measurements"` // unused collection ids
	Enabled                  bool            `bson:"enabled" json:"enabled"`
	Infos                    map[string]Info `bson:"infos" json:"infos"`
	ComponentInfos           map[string]Info `bson:"component_infos,omitempty" json:"component_infos,omitempty"`
	Type                     string          `bson:"type" json:"type"`
	Component                string          `bson:"component,omitempty" json:"component,omitempty"`
	Category                 string          `bson:"category" json:"category"`
	ImpactLevel              int64           `bson:"impact_level" json:"impact_level"`
	AlarmsCumulativeData     struct {
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
	// LastIdleRuleApply is used to mark entity if some idle rule was applied.
	LastIdleRuleApply string `bson:"last_idle_rule_apply,omitempty" json:"last_idle_rule_apply,omitempty"`

	// IdleSince represents since when entity didn't receive any events.
	IdleSince    *CpsTime `bson:"idle_since,omitempty" json:"idle_since,omitempty"`
	ImportSource string   `bson:"import_source,omitempty" json:"import_source"`

	Imported      *CpsTime      `bson:"imported,omitempty" json:"imported"`
	PbehaviorInfo PbehaviorInfo `bson:"pbehavior_info,omitempty" json:"pbehavior_info,omitempty"`

	SliAvailState int64 `bson:"sli_avail_state" json:"sli_avail_state"`

	IsNew bool `bson:"-" json:"-"`
}

func (e *Entity) GetUpsertMongoBson(newImpacts []string, newDepends []string) bson.M {
	set := bson.M{
		"_id":            e.ID,
		"name":           e.Name,
		"measurements":   e.Measurements,
		"enabled":        e.Enabled,
		"type":           e.Type,
		"enable_history": e.EnableHistory,
	}
	if e.Component != "" {
		set["component"] = e.Component
	}

	addToSet := bson.M{
		"impact": bson.M{
			"$each": newImpacts,
		},
		"depends": bson.M{
			"$each": newDepends,
		},
	}

	setOnInsert := bson.M{}

	if len(e.Infos) == 0 {
		setOnInsert["infos"] = bson.M{}
	} else {
		for k, v := range e.Infos {
			key := "infos." + k
			set[key] = v
		}
	}

	upsertBson := bson.M{
		"$set":      set,
		"$addToSet": addToSet,
	}

	if len(setOnInsert) > 0 {
		upsertBson["$setOnInsert"] = setOnInsert
	}

	return upsertBson
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

// NewInfo instanciate a new info struct [sic]
func NewInfo(name string, description string, value interface{}) Info {
	return Info{
		Name:        name,
		Description: description,
		Value:       value,
	}
}

// NewEntity instanciate a new entity struct [sic]
func NewEntity(id string, name string, entityType string, infos map[string]Info, impacts, depends []string) Entity {
	history := []CpsTime{CpsTime{time.Now()}}

	if infos == nil {
		infos = make(map[string]Info)
	}
	component := ""
	switch entityType {
	case EntityTypeComponent:
		component = id
	case EntityTypeResource, EntityTypeMetaAlarm:
		idParts := strings.Split(id, "/")
		component = idParts[len(idParts)-1]
	}

	return Entity{
		ID:            id,
		Name:          name,
		Impacts:       impacts,
		Depends:       depends,
		EnableHistory: history,
		Measurements:  nil,
		Enabled:       true,
		Infos:         infos,
		Type:          entityType,
		Component:     component,
		IsNew:         true,
	}
}
