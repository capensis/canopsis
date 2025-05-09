package notification

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/usernotification"
)

type Response struct {
	ID      string                 `json:"_id" bson:"_id"`
	Type    int                    `json:"type" bson:"type"`
	Time    datetime.CpsTime       `json:"time" bson:"time" swaggertype:"integer"`
	Author  *author.Author         `json:"author" bson:"author"`
	Comment string                 `json:"comment" bson:"comment"`
	Rule    *usernotification.Rule `json:"rule" bson:"rule"`
}

type UpdateSettingsRequest struct {
	Instruction InstructionNotificationSettings `json:"instruction"`
	Author      string                          `json:"author" swaggerignore:"true"`
}

type SettingsResponse struct {
	Instruction InstructionNotificationSettings `json:"instruction" bson:"instruction"`
}

type Settings struct {
	ID          string                          `bson:"_id,omitempty"`
	Instruction InstructionNotificationSettings `bson:"instruction"`
	Author      string                          `bson:"author"`
	Updated     datetime.CpsTime                `bson:"updated"`
}

type InstructionNotificationSettings struct {
	Rate          *bool                     `json:"rate" bson:"rate" binding:"required"`
	RateFrequency datetime.DurationWithUnit `json:"rate_frequency" bson:"rate_frequency" binding:"required"`
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
