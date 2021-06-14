package types

import (
	"fmt"
	"log"
	"strings"
	"time"

	mgobson "github.com/globalsign/mgo/bson"
	mongobson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Entity types
const (
	EntityTypeConnector = "connector"
	EntityTypeComponent = "component"
	EntityTypeResource  = "resource"
	EntityTypeWatcher   = "watcher"
	EntityTypeMetaAlarm = "metaalarm"
)

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
	Value       string      `bson:"-" json:"value"`
	RealValue   interface{} `bson:"value,omitempty" json:"-"`
}

// Entity ...
type Entity struct {
	ID                   string          `bson:"_id" json:"_id"`
	Name                 string          `bson:"name" json:"name"`
	Impacts              []string        `bson:"impact" json:"impact"`
	Depends              []string        `bson:"depends" json:"depends"`
	EnableHistory        []CpsTime       `bson:"enable_history" json:"enable_history"`
	Measurements         interface{}     `bson:"measurements" json:"measurements"` // unused collection ids
	Enabled              bool            `bson:"enabled" json:"enabled"`
	Infos                map[string]Info `bson:"infos" json:"infos"`
	ComponentInfos       map[string]Info `bson:"component_infos,omitempty" json:"component_infos,omitempty"`
	Type                 string          `bson:"type" json:"type"`
	Component            string          `bson:"component,omitempty" json:"component,omitempty"`
	IsNew                bool            `bson:"-" json:"-"`
	AlarmsCumulativeData struct {
		// Only for watcher.
		// WatchedCount is count of unresolved alarms.
		WatchedCount int64 `bson:"watched_count"`
		// WatchedPbheaviorCount contains counters of unresolved and in pbehavior alarms.
		WatchedPbheaviorCount map[string]int64 `bson:"watched_pbehavior_count"`
		// WatchedNotAckedCount is count of unresolved and not acked and active (by pbehavior) alarms.
		WatchedNotAckedCount int64 `bson:"watched_not_acked_count"`
	} `bson:"alarms_cumulative_data,omitempty" json:"-"`
}

func (info *Info) SetBSON(raw mgobson.Raw) error {
	m := make(map[string]interface{})
	err := raw.Unmarshal(&m)
	if err != nil {
		log.Printf("unable to parse alarm pattern list: %v", err)
		return nil
	}

	if v, ok := m["name"]; ok {
		info.Name = fmt.Sprintf("%v", v)
	}

	if v, ok := m["description"]; ok {
		info.Description = fmt.Sprintf("%v", v)
	}

	if v, ok := m["value"]; ok {
		info.Value = fmt.Sprintf("%v", v)
		info.RealValue = v
	}

	return nil
}

func (info *Info) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	m := make(map[string]interface{})
	err := mongobson.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	if v, ok := m["name"]; ok {
		info.Name = fmt.Sprintf("%v", v)
	}

	if v, ok := m["description"]; ok {
		info.Description = fmt.Sprintf("%v", v)
	}

	if v, ok := m["value"]; ok {
		info.Value = fmt.Sprintf("%v", v)
		info.RealValue = v
	}

	return nil
}

func (e *Entity) GetUpsertMongoBson(newImpacts []string, newDepends []string) mgobson.M {
	set := mgobson.M{
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

	addToSet := mgobson.M{
		"impact": mgobson.M{
			"$each": newImpacts,
		},
		"depends": mgobson.M{
			"$each": newDepends,
		},
	}

	setOnInsert := mgobson.M{}

	if len(e.Infos) == 0 {
		setOnInsert["infos"] = mgobson.M{}
	} else {
		for k, v := range e.Infos {
			key := "infos." + k
			set[key] = v
		}
	}

	upsertBson := mgobson.M{
		"$set":      set,
		"$addToSet": addToSet,
	}

	if len(setOnInsert) > 0 {
		upsertBson["$setOnInsert"] = setOnInsert
	}

	return upsertBson
}

func (e *Entity) getSelector() mgobson.M {
	return mgobson.M{}
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

// NewInfo instanciate a new info struct [sic]
func NewInfo(name string, description string, realValue interface{}) Info {
	return Info{
		Name:        name,
		Description: description,
		Value:       fmt.Sprintf("%v", realValue),
		RealValue:   realValue,
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
