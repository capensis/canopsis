package exdate

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"

type Request struct {
	Begin types.CpsTime `json:"begin" binding:"required" swaggertype:"integer"`
	End   types.CpsTime `json:"end" binding:"required" swaggertype:"integer"`
}
