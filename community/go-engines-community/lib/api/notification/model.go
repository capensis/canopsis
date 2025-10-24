package notification

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"

type Notification struct {
	Instruction InstructionNotification `json:"instruction" bson:"instruction"`
}

type InstructionNotification struct {
	Rate          *bool                     `json:"rate" bson:"rate" binding:"required"`
	RateFrequency datetime.DurationWithUnit `json:"rate_frequency" bson:"rate_frequency" binding:"required"`
}
