package types

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"

type Exdate struct {
	Begin datetime.CpsTime `bson:"begin" json:"begin" swaggertype:"integer"`
	End   datetime.CpsTime `bson:"end" json:"end" swaggertype:"integer"`
}
