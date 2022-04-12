package pbehaviortimespan

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type TimespansRequest struct {
	StartAt types.CpsTime `json:"start_at" binding:"required" swaggertype:"integer"`
	EndAt   types.CpsTime `json:"end_at" swaggertype:"integer"`
	RRule   string        `json:"rrule"`
	Type    string        `json:"type" binding:"required"`

	Exdates    []pbehaviorexception.ExdateRequest `json:"exdates" binding:"dive"`
	Exceptions []string                           `json:"exceptions"`

	ViewFrom types.CpsTime `json:"view_from" binding:"required" swaggertype:"integer"`
	ViewTo   types.CpsTime `json:"view_to" binding:"required" swaggertype:"integer"`
}

type ItemResponse struct {
	From types.CpsTime  `json:"from" swaggertype:"integer"`
	To   types.CpsTime  `json:"to" swaggertype:"integer"`
	Type pbehavior.Type `json:"type"`
}
