package widgetfilter

import (
	"encoding/json"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
)

type ListRequest struct {
	pagination.Query
	Widget  string `form:"widget" json:"widget" binding:"required"`
	Private *bool  `form:"private" json:"private"`
}

type EditRequest struct {
	BaseEditRequest
	IsUserPreference *bool  `json:"is_user_preference" binding:"required"`
	Author           string `json:"author" swaggerignore:"true"`
	IsPrivate        bool   `json:"-"`
}

type BaseEditRequest struct {
	Title string `json:"title" binding:"required,max=255"`

	common.AlarmPatternFieldsRequest
	common.EntityPatternFieldsRequest
	common.PbehaviorPatternFieldsRequest

	WeatherServicePattern view.WeatherServicePattern `json:"weather_service_pattern"`
}

type CreateRequest struct {
	EditRequest
	Widget string `json:"widget" binding:"required"`
}

type UpdateRequest struct {
	EditRequest
	ID string `json:"-"`
}

type Response struct {
	ID               string            `bson:"_id" json:"_id"`
	Widget           string            `bson:"widget" json:"-"`
	Title            string            `bson:"title" json:"title"`
	IsUserPreference bool              `bson:"is_user_preference" json:"is_user_preference"`
	Author           *author.Author    `bson:"author" json:"author,omitempty"`
	Created          *datetime.CpsTime `bson:"created" json:"created,omitempty" swaggertype:"integer"`
	Updated          *datetime.CpsTime `bson:"updated" json:"updated,omitempty" swaggertype:"integer"`

	savedpattern.AlarmPatternFields     `bson:",inline"`
	savedpattern.EntityPatternFields    `bson:",inline"`
	savedpattern.PbehaviorPatternFields `bson:",inline"`

	WeatherServicePattern view.WeatherServicePattern `bson:"weather_service_pattern" json:"weather_service_pattern,omitempty"`

	IsPrivate bool `bson:"is_private" json:"is_private"`
}

type AggregationResult struct {
	Data       []Response `bson:"data" json:"data"`
	TotalCount int64      `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

type EditPositionRequest struct {
	Items []string `json:"items" binding:"required,notblank,unique"`
}

func (r EditPositionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Items)
}

func (r *EditPositionRequest) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.Items)
}
