package pbehavior

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/teambition/rrule-go"
)

const (
	PBehaviorCollectionName = mongo.PbehaviorMongoCollection
	TypeCollectionName      = mongo.PbehaviorTypeMongoCollection
)

type Comment struct {
	ID        string         `bson:"_id" json:"_id"`
	Author    string         `bson:"author" json:"author"`
	Timestamp *types.CpsTime `bson:"ts" json:"ts"`
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
	Filter     string         `bson:"filter"`
	Name       string         `bson:"name"`
	Reason     string         `bson:"reason"`
	RRule      string         `bson:"rrule"`
	Start      *types.CpsTime `bson:"tstart"`
	Stop       *types.CpsTime `bson:"tstop,omitempty"`
	Type       string         `bson:"type_"`
	Exdates    []Exdate       `bson:"exdates"`
	Exceptions []string       `bson:"exceptions"`
	Created    types.CpsTime  `bson:"created,omitempty"`
	Updated    types.CpsTime  `bson:"updated,omitempty"`
}

// isActiveRRule checks if a pbehavior with a rrule is active at the given time.
// Do no call this method directly, instead use IsActive.
func (pb PBehavior) isActiveRRule(date time.Time) (bool, error) {
	start := pb.Start.Time

	for _, exDate := range pb.Exdates {
		exDateBegin := exDate.Begin.Time
		exDateEnd := exDate.End.Time

		if date.After(exDateBegin) && date.Before(exDateEnd) {
			return false, nil
		}
	}

	recSet := rrule.Set{}

	rOption, err := rrule.StrToROption(pb.RRule)
	if err != nil {
		return false, err
	}
	rOption.Dtstart = start

	rr, err := rrule.NewRRule(*rOption)
	if err != nil {
		return false, err
	}

	recSet.RRule(rr)
	recSet.DTStart(start)

	recStart := recSet.Before(date, true)

	if recStart.IsZero() {
		return false, nil
	}

	if pb.Stop == nil {
		if recStart.Before(date) {
			return true, nil
		}
	} else {
		stop := pb.Stop.Time
		duration := stop.Sub(start)
		recEnd := recStart.Add(duration)

		if recStart.Before(date) && recEnd.After(date) {
			return true, nil
		}
	}

	return false, nil
}

// IsActive checks if a pbehavior is active at the given time.
func (pb PBehavior) IsActive(date time.Time) (bool, error) {
	if pb.RRule == "" {
		if pb.Stop == nil {
			return pb.Start.Before(date), nil
		}
		return pb.Start.Before(date) && pb.Stop.After(date), nil
	} else {
		return pb.isActiveRRule(date)
	}

}
