package messageratestats

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"

const (
	IntervalMinute = "minute"
	IntervalHour   = "hour"
)

type ListRequest struct {
	Interval string           `form:"interval" json:"interval" binding:"required,oneof=minute hour"`
	From     datetime.CpsTime `form:"from" json:"from" binding:"required" swaggertype:"integer"`
	To       datetime.CpsTime `form:"to" json:"to" binding:"required" swaggertype:"integer"`
}

type StatsResponse struct {
	ID   int64 `bson:"_id" json:"time"`
	Rate int64 `bson:"received" json:"rate"`
}

type StatsListResponse struct {
	Data []StatsResponse `json:"data"`
	Meta struct {
		DeletedBefore *datetime.CpsTime `json:"deleted_before,omitempty" swaggertype:"integer"`
	} `json:"meta"`
}
