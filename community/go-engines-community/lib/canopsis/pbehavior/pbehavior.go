package pbehavior

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

const (
	PBehaviorCollectionName = mongo.PbehaviorMongoCollection
	TypeCollectionName      = mongo.PbehaviorTypeMongoCollection
)

type Comment struct {
	ID        string         `bson:"_id" json:"_id"`
	Author    string         `bson:"author" json:"author"`
	Timestamp *types.CpsTime `bson:"ts" json:"ts" swaggertype:"integer"`
	Message   string         `bson:"message" json:"message"`
}

const (
	TypeActive      = "active"
	TypeMaintenance = "maintenance"
	TypePause       = "pause"
	TypeInactive    = "inactive"
)

type Type struct {
	ID          string `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Type        string `bson:"type" json:"type"`
	Priority    int    `bson:"priority" json:"priority"`
	IconName    string `bson:"icon_name" json:"icon_name"`
	Color       string `bson:"color,omitempty" json:"color,omitempty"`
}

type Comments []*Comment

// PBehavior represents a canopsis periodical behavior.
type PBehavior struct {
	ID         string         `bson:"_id,omitempty"`
	Author     string         `bson:"author"`
	Comments   Comments       `bson:"comments,omitempty"`
	Enabled    bool           `bson:"enabled"`
	Name       string         `bson:"name"`
	Reason     string         `bson:"reason"`
	RRule      string         `bson:"rrule"`
	Start      *types.CpsTime `bson:"tstart"`
	Stop       *types.CpsTime `bson:"tstop,omitempty"`
	Type       string         `bson:"type_"`
	Exdates    []Exdate       `bson:"exdates"`
	Exceptions []string       `bson:"exceptions"`
	Created    *types.CpsTime `bson:"created,omitempty"`
	Updated    *types.CpsTime `bson:"updated,omitempty"`

	LastAlarmDate *types.CpsTime `bson:"last_alarm_date,omitempty"`

	savedpattern.EntityPatternFields `bson:",inline"`
	// Deprecated : contains old mongo query which cannot be migrated to pattern.
	OldMongoQuery map[string]interface{} `bson:"old_mongo_query,omitempty" json:"old_mongo_query,omitempty"`
}
