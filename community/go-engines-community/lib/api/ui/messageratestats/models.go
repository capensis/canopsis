package messageratestats

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"

const (
	IntervalMinute = "minute"
	IntervalHour   = "hour"
)

type SearchRequest struct {
	EventTypes     []string `form:"event_types[]" json:"event_types"`
	ConnectorNames []string `form:"connector_names[]" json:"connector_names"`
}

type ListRequest struct {
	Interval string           `form:"interval" json:"interval" binding:"required,oneof=minute hour"`
	From     datetime.CpsTime `form:"from" json:"from" swaggertype:"integer"`
	To       datetime.CpsTime `form:"to" json:"to" swaggertype:"integer"`

	SearchRequest
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
