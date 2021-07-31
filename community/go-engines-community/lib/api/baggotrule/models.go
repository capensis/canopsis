package baggotrule

import (
	"bytes"
	"encoding/json"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type PatternsRaw struct {
	AlarmPatterns  interface{} `json:"alarm_patterns"`
	EntityPatterns interface{} `json:"entity_patterns"`
}

type User struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"crecord_name" json:"name"`
}

type Payload struct {
	ID             string                     `bson:"_id" json:"_id" binding:"id"`
	Description    string                     `bson:"description" json:"description" binding:"required,max=255"`
	Duration       types.DurationWithUnit     `bson:"duration" json:"duration" binding:"required"`
	AlarmPatterns  *pattern.AlarmPatternList  `bson:"alarm_patterns" json:"alarm_patterns"`
	EntityPatterns *pattern.EntityPatternList `bson:"entity_patterns" json:"entity_patterns"`
	Created        *types.CpsTime             `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated        *types.CpsTime             `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
	Priority       *int                       `bson:"priority" json:"priority" binding:"required,min=0"`
	PRaw           PatternsRaw                `bson:"-" json:"-"`
}

type CreateRequest struct {
	Payload `bson:",inline"`
	Author  string `bson:"author" json:"author" swaggerignore:"true"`
}

type UpdateRequest struct {
	ID            string `bson:"-" json:"-"`
	CreateRequest `bson:",inline"`
}

type RuleResponse struct {
	Payload `bson:",inline"`
	Author  User `bson:"author" json:"author"`
}

type AggregationResult struct {
	Data       []RuleResponse `bson:"data" json:"data"`
	TotalCount int64          `bson:"total_count" json:"total_count"`
}

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id author created updated"`
}

func (mp *UpdateRequest) UnmarshalJSON(b []byte) error {
	var r CreateRequest
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	mp.CreateRequest = r
	mp.ID = ""
	mp.CreateRequest.ID = ""
	return nil
}

func (mp *CreateRequest) UnmarshalJSON(b []byte) error {
	var p PatternsRaw
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	type DRule CreateRequest
	var payload DRule
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		// patten unmarshall err
		// do not throw the error, the struct validator will take care the corresponding error
		if !errors.Is(err, pattern.UnexpectedFieldsError{}) {
			return err
		}
	}

	*mp = CreateRequest(payload)
	mp.PRaw = p
	return nil
}
