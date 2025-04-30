package usernotification

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

const (
	TypeInstructionApprove = iota
	TypeInstructionRate
	TypeEventFilterFailure
)

type Notification struct {
	ID      string           `bson:"_id"`
	Type    int              `bson:"type"`
	User    string           `bson:"user,omitempty"`
	Roles   []string         `bson:"roles,omitempty"`
	Time    datetime.CpsTime `bson:"time" swaggertype:"integer"`
	Author  string           `bson:"author,omitempty"`
	Comment string           `bson:"comment"`
	Rule    *Rule            `bson:"rule,omitempty"`
}

type Rule struct {
	ID   string `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}
