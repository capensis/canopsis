package pbehaviortimespan

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"
)

type TimespansRequest struct {
	StartAt libtime.CpsTime `json:"start_at" binding:"required" swaggertype:"integer"`
	EndAt   libtime.CpsTime `json:"end_at" swaggertype:"integer"`
	RRule   string          `json:"rrule"`
	Type    string          `json:"type" binding:"required"`

	Exdates    []pbehaviorexception.ExdateRequest `json:"exdates" binding:"dive"`
	Exceptions []string                           `json:"exceptions"`

	ViewFrom libtime.CpsTime `json:"view_from" binding:"required" swaggertype:"integer"`
	ViewTo   libtime.CpsTime `json:"view_to" binding:"required" swaggertype:"integer"`
}

type ItemResponse struct {
	From libtime.CpsTime `json:"from" swaggertype:"integer"`
	To   libtime.CpsTime `json:"to" swaggertype:"integer"`
	Type pbehavior.Type  `json:"type"`
}
