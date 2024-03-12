package exdate

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"

type Request struct {
	Begin datetime.CpsTime `json:"begin" binding:"required" swaggertype:"integer"`
	End   datetime.CpsTime `json:"end" binding:"required" swaggertype:"integer"`
}
