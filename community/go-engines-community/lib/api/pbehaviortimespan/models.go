package pbehaviortimespan

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
)

type TimespansRequest struct {
	StartAt datetime.CpsTime `json:"start_at" binding:"required" swaggertype:"integer"`
	EndAt   datetime.CpsTime `json:"end_at" swaggertype:"integer"`
	RRule   string           `json:"rrule"`
	Type    string           `json:"type" binding:"required"`

	Exdates    []pbehaviorexception.ExdateRequest `json:"exdates" binding:"dive"`
	Exceptions []string                           `json:"exceptions"`

	ViewFrom datetime.CpsTime `json:"view_from" binding:"required" swaggertype:"integer"`
	ViewTo   datetime.CpsTime `json:"view_to" binding:"required" swaggertype:"integer"`
}

type ItemResponse struct {
	From datetime.CpsTime `json:"from" swaggertype:"integer"`
	To   datetime.CpsTime `json:"to" swaggertype:"integer"`
	Type pbehavior.Type   `json:"type"`
}
