package pbehaviortimespan

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type TimespansRequest struct {
	StartAt    types.CpsTime   `json:"start_at" binding:"required" swaggertype:"integer"`
	EndAt      *types.CpsTime  `json:"end_at" swaggertype:"integer"`
	RRule      string          `json:"rrule"`
	Exdates    []ExdateRequest `json:"exdates" binding:"dive"`
	Exceptions []string        `json:"exceptions"`
	ViewFrom   types.CpsTime   `json:"view_from" binding:"required" swaggertype:"integer"`
	ViewTo     types.CpsTime   `json:"view_to" binding:"required" swaggertype:"integer"`
	ByDate     bool            `json:"by_date"`
}

type ExdateRequest struct {
	Begin types.CpsTime `json:"begin" binding:"required" swaggertype:"integer"`
	End   types.CpsTime `json:"end" binding:"required" swaggertype:"integer"`
	Type  string        `json:"type,omitempty"`
}

type timespansItemResponse struct {
	From types.CpsTime   `json:"from" swaggertype:"integer"`
	To   types.CpsTime   `json:"to" swaggertype:"integer"`
	Type *pbehavior.Type `json:"type,omitempty"`
}
