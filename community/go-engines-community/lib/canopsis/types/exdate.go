package types

import libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"

type Exdate struct {
	Begin libtime.CpsTime `bson:"begin" json:"begin" swaggertype:"integer"`
	End   libtime.CpsTime `bson:"end" json:"end" swaggertype:"integer"`
}
