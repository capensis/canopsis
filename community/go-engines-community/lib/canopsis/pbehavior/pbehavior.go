package pbehavior

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
)

type Comment struct {
	ID        string            `bson:"_id" json:"_id"`
	Author    string            `bson:"author" json:"author"`
	Origin    string            `bson:"origin,omitempty" json:"origin,omitempty"`
	Timestamp *datetime.CpsTime `bson:"ts" json:"ts" swaggertype:"integer"`
	Message   string            `bson:"message" json:"message"`
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
	Color       string `bson:"color" json:"color"`
}

type Comments []*Comment

// PBehavior represents a canopsis periodical behavior.
type PBehavior struct {
	ID            string            `bson:"_id,omitempty"`
	Author        string            `bson:"author"`
	Comments      Comments          `bson:"comments,omitempty"`
	Enabled       bool              `bson:"enabled"`
	Name          string            `bson:"name"`
	Reason        string            `bson:"reason"`
	Type          string            `bson:"type_"`
	Exdates       []Exdate          `bson:"exdates"`
	Exceptions    []string          `bson:"exceptions"`
	Color         string            `bson:"color"`
	Created       *datetime.CpsTime `bson:"created,omitempty"`
	Updated       *datetime.CpsTime `bson:"updated,omitempty"`
	LastAlarmDate *datetime.CpsTime `bson:"last_alarm_date,omitempty"`

	Start    *datetime.CpsTime `bson:"tstart"`
	Stop     *datetime.CpsTime `bson:"tstop,omitempty"`
	RRule    string            `bson:"rrule"`
	RRuleEnd *datetime.CpsTime `bson:"rrule_end,omitempty"`
	// RRuleComputedStart is an auxiliary start date to compute rrule faster.
	RRuleComputedStart *datetime.CpsTime `bson:"rrule_cstart,omitempty"`

	// Origin is used if a pbehavior is created for certain entities.
	// Origin can contain some feature name or external service name.
	Origin string `bson:"origin,omitempty"`
	// Entity is used if a pbehavior is created for one certain entity.
	Entity string `bson:"entity,omitempty"`
	// Entities is used if a pbehavior is created for multiple certain entities.
	Entities []string `bson:"entities,omitempty"`

	savedpattern.EntityPatternFields `bson:",inline"`
}
