package alarmtag

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"
)

const (
	TypeExternal = iota
	TypeInternal
)

type AlarmTag struct {
	ID      string          `bson:"_id" json:"_id"`
	Type    int64           `bson:"type" json:"type"`
	Value   string          `bson:"value" json:"value"`
	Label   string          `bson:"label" json:"label"`
	Color   string          `bson:"color" json:"color"`
	Author  string          `bson:"author" json:"author"`
	Created libtime.CpsTime `bson:"created" json:"created"`
	Updated libtime.CpsTime `bson:"updated" json:"updated"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}
