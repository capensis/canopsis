package messageratestats

import libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"

const (
	IntervalMinute = "minute"
	IntervalHour   = "hour"
)

type ListRequest struct {
	Interval string          `form:"interval" json:"interval" binding:"required,oneof=minute hour"`
	From     libtime.CpsTime `form:"from" json:"from" binding:"required" swaggertype:"integer"`
	To       libtime.CpsTime `form:"to" json:"to" binding:"required" swaggertype:"integer"`
}

type StatsResponse struct {
	ID       int64 `bson:"_id" json:"time"`
	Received int64 `bson:"received" json:"rate"`
}

type StatsListResponse struct {
	Data []StatsResponse `json:"data"`
	Meta struct {
		DeletedBefore *libtime.CpsTime `json:"deleted_before,omitempty" swaggertype:"integer"`
	} `json:"meta"`
}
