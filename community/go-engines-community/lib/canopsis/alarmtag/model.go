package alarmtag

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
)

const (
	TypeExternal = iota
	TypeInternal
)

type AlarmTag struct {
	ID            string           `bson:"_id" json:"_id"`
	Type          int64            `bson:"type" json:"type"`
	Value         string           `bson:"value" json:"value"`
	Label         string           `bson:"label" json:"label"`
	Color         string           `bson:"color" json:"color"`
	Author        string           `bson:"author" json:"author"`
	Created       datetime.CpsTime `bson:"created" json:"created"`
	Updated       datetime.CpsTime `bson:"updated" json:"updated"`
	LastEventDate datetime.CpsTime `bson:"last_event_date,omitempty" json:"last_event_date,omitempty"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}
