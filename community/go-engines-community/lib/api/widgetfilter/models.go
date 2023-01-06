package widgetfilter

import (
	"encoding/json"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type ListRequest struct {
	pagination.Query
	Widget  string `form:"widget" json:"widget" binding:"required"`
	Private *bool  `form:"private" json:"private"`
}

type EditRequest struct {
	BaseEditRequest
	ID        string `json:"-"`
	Widget    string `json:"widget" binding:"required"`
	IsPrivate *bool  `json:"is_private" binding:"required"`
	Author    string `json:"author" swaggerignore:"true"`
}

type BaseEditRequest struct {
	Title string `json:"title" binding:"required,max=255"`

	common.AlarmPatternFieldsRequest
	common.EntityPatternFieldsRequest
	common.PbehaviorPatternFieldsRequest

	WeatherServicePattern view.WeatherServicePattern `json:"weather_service_pattern"`
}

type Response struct {
	ID        string         `bson:"_id" json:"_id"`
	Widget    string         `bson:"widget" json:"-"`
	Title     string         `bson:"title" json:"title"`
	IsPrivate *bool          `bson:"is_private" json:"is_private,omitempty"`
	Author    *author.Author `bson:"author" json:"author,omitempty"`
	Created   *types.CpsTime `bson:"created" json:"created,omitempty" swaggertype:"integer"`
	Updated   *types.CpsTime `bson:"updated" json:"updated,omitempty" swaggertype:"integer"`

	OldMongoQuery OldMongoQuery `bson:"old_mongo_query" json:"old_mongo_query,omitempty"`

	savedpattern.AlarmPatternFields     `bson:",inline"`
	savedpattern.EntityPatternFields    `bson:",inline"`
	savedpattern.PbehaviorPatternFields `bson:",inline"`

	WeatherServicePattern view.WeatherServicePattern `bson:"weather_service_pattern" json:"weather_service_pattern,omitempty"`
}

type OldMongoQuery map[string]interface{}

func (q *OldMongoQuery) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	v, _, ok := bsoncore.ReadString(b)
	if !ok {
		return errors.New("invalid value, expected string")
	}

	err := json.Unmarshal([]byte(v), &q)
	return err
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
