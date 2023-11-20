package exdate

import libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"

type Request struct {
	Begin libtime.CpsTime `json:"begin" binding:"required" swaggertype:"integer"`
	End   libtime.CpsTime `json:"end" binding:"required" swaggertype:"integer"`
}
