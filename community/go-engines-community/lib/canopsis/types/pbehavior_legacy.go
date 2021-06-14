package types

import (
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"github.com/teambition/rrule-go"
	"os"
	"time"
)

// PBehavior represents a comment in a PBehavior.
type PBehaviorComment struct {
	ID        string   `bson:"_id" json:"_id"`
	Author    string   `bson:"author" json:"author"`
	Timestamp *CpsTime `bson:"ts" json:"ts"`
	Message   string   `bson:"message" json:"message"`
}

type PBehaviorComments []*PBehaviorComment

// PBehavior represents a canopsis periodical behavior.
type PBehaviorLegacy struct {
	ID            string            `bson:"_id,omitempty" json:"_id,omitempty"`
	Author        string            `bson:"author" json:"author"`
	Comments      PBehaviorComments `bson:"comments,omitempty" json:"comments,omitempty"`
	Connector     string            `bson:"connector" json:"connector"`
	ConnectorName string            `bson:"connector_name" json:"connector_name"`
	Eids          []string          `bson:"eids,omitempty" json:"eids,omitempty"`
	Enabled       bool              `bson:"enabled" json:"enabled"`
	Filter        string            `bson:"filter" json:"filter"`
	Name          string            `bson:"name" json:"name"`
	Reason        string            `bson:"reason" json:"reason"`
	RRule         string            `bson:"rrule" json:"rrule"`
	Start         *CpsTime          `bson:"tstart" json:"tstart"`
	Stop          *CpsTime          `bson:"tstop" json:"tstop"`
	Type          string            `bson:"type_" json:"type_"`
	Timezone      string            `bson:"timezone,omitempty" json:"timezone"`
	Exdate        []CpsTime         `bson:"exdate" json:"exdate"`
}

// NewPBehavior instanciate a new pbehavior
func NewPBehaviorLegacy(author, filter, name, reason, type_, rrule string, start, stop *CpsTime, timezone string, exdate []CpsTime) PBehaviorLegacy {
	comments := PBehaviorComments{}
	eids := []string{}

	if timezone == "" {
		timezoneName := os.Getenv("CPS_PBH_TIMEZONE")
		if timezoneName == "" {
			timezoneName = "Europe/Paris"
		}
		defaultLocation, _ := time.LoadLocation(timezoneName)
		timezone = defaultLocation.String()
	}

	return PBehaviorLegacy{
		ID:            utils.NewID(),
		Author:        author,
		Comments:      comments,
		Connector:     "",
		ConnectorName: "",
		Eids:          eids,
		Enabled:       true,
		Filter:        filter,
		Name:          name,
		Reason:        reason,
		RRule:         rrule,
		Start:         start,
		Stop:          stop,
		Type:          type_,
		Timezone:      timezone,
		Exdate:        exdate,
	}
}

// IsImpacted checks if an entity id is impacted by a pbehavior
func (pb PBehaviorLegacy) IsImpacted(eid string) bool {
	return utils.StringInSlice(eid, pb.Eids)
}

// isActiveRRule checks if a pbehavior with a rrule is active at the given time.
// Do no call this method directly, instead use IsActive.
func (pb PBehaviorLegacy) isActiveRRule(date time.Time) (bool, error) {
	pbhLocation, err := time.LoadLocation(pb.Timezone)
	if err != nil {
		return false, err
	}

	start := pb.Start.In(pbhLocation)
	stop := pb.Stop.In(pbhLocation)
	date = date.In(pbhLocation)
	duration := stop.Sub(start)

	recSet := rrule.Set{}

	rOption, err := rrule.StrToROption(pb.RRule)
	if err != nil {
		return false, err
	}
	rOption.Dtstart = start

	rr, err := rrule.NewRRule(*rOption)

	recSet.RRule(rr)
	recSet.DTStart(start)

	for _, date := range pb.Exdate {
		recSet.ExDate(date.In(pbhLocation))
	}

	recStart := recSet.Before(date, true)

	if recStart.IsZero() {
		return false, nil
	}

	recEnd := recStart.Add(duration)

	if recStart.Before(date) && recEnd.After(date) {
		return true, nil
	}

	return false, nil
}

// IsActive checks if a pbehavior is active at the given time.
func (pb PBehaviorLegacy) IsActive(date time.Time) (bool, error) {
	if pb.RRule == "" {
		return pb.Start.Before(date) && pb.Stop.After(date), nil
	} else {
		return pb.isActiveRRule(date)
	}

}
