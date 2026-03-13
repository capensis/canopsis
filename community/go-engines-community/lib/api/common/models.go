package common

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const LimitLinkedRules = 11

type AlarmStep struct {
	Type      string            `bson:"_t" json:"_t"`
	Timestamp *datetime.CpsTime `bson:"t" json:"t" swaggertype:"integer"`
	Author    string            `bson:"a" json:"a"`
	UserID    string            `bson:"user_id,omitempty" json:"user_id"`
	Message   string            `bson:"m" json:"m"`
	Value     types.CpsNumber   `bson:"val" json:"val"`
	Initiator string            `bson:"initiator" json:"initiator"`
	Execution string            `bson:"exec,omitempty" json:"-"`
	IconName  string            `bson:"icon_name,omitempty" json:"icon_name,omitempty"`
	Color     string            `bson:"color,omitempty" json:"color,omitempty"`

	// Ticket related fields
	types.TicketInfo `bson:",inline"`

	InPbehaviorInterval bool `bson:"in_pbh,omitempty" json:"in_pbh,omitempty"`

	StructuredMessage []types.StructuredMessage `bson:"struct_m,omitempty" json:"struct_m,omitempty"`
}
