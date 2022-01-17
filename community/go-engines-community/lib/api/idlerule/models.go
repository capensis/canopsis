package idlerule

import (
	"encoding/json"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/scenario"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/mitchellh/mapstructure"
)

type CountByPatternRequest struct {
	EntityPatterns pattern.EntityPatternList `json:"entity_patterns"`
	AlarmPatterns  pattern.AlarmPatternList  `json:"alarm_patterns"`
}

type CountByPatternResult struct {
	OverLimit          bool  `json:"over_limit"`
	TotalCountEntities int64 `json:"total_count_entities"`
	TotalCountAlarms   int64 `json:"total_count_alarms"`
}

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id name author duration created updated type priority"`
}

type EditRequest struct {
	Name                 string                    `json:"name" binding:"required,max=255"`
	Description          string                    `json:"description" binding:"max=255"`
	Author               string                    `json:"author" swaggerignore:"true"`
	Enabled              *bool                     `json:"enabled" binding:"required"`
	Type                 string                    `json:"type" binding:"required"`
	Priority             *int64                    `json:"priority" binding:"required"`
	Duration             types.DurationWithUnit    `json:"duration" binding:"required"`
	DisableDuringPeriods []string                  `json:"disable_during_periods"`
	EntityPatterns       pattern.EntityPatternList `json:"entity_patterns"`
	AlarmPatterns        pattern.AlarmPatternList  `json:"alarm_patterns"`
	AlarmCondition       string                    `json:"alarm_condition"`
	Operation            *OperationRequest         `json:"operation,omitempty"`
}

type CreateRequest struct {
	EditRequest
	ID string `json:"_id" binding:"id"`
}

type UpdateRequest struct {
	EditRequest
	ID string `json:"-"`
}

type BulkUpdateRequestItem struct {
	EditRequest
	ID string `json:"_id" binding:"required"`
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}

// for swagger
type BulkCreateResponseItem struct {
	ID     string            `json:"id,omitempty"`
	Item   CreateRequest     `json:"item"`
	Status int               `json:"status"`
	Error  string            `json:"error,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

// for swagger
type BulkUpdateResponseItem struct {
	ID     string                `json:"id,omitempty"`
	Item   BulkUpdateRequestItem `json:"item"`
	Status int                   `json:"status"`
	Error  string                `json:"error,omitempty"`
	Errors map[string]string     `json:"errors,omitempty"`
}

// for swagger
type BulkDeleteResponseItem struct {
	ID     string                `json:"id,omitempty"`
	Item   BulkDeleteRequestItem `json:"item"`
	Status int                   `json:"status"`
	Error  string                `json:"error,omitempty"`
	Errors map[string]string     `json:"errors,omitempty"`
}

type OperationRequest struct {
	Type       string      `json:"type" binding:"required"`
	Parameters interface{} `json:"parameters,omitempty"`
}

func (r *OperationRequest) UnmarshalJSON(b []byte) error {
	type Alias OperationRequest
	tmp := struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	switch r.Type {
	case types.ActionTypeSnooze:
		var params scenario.SnoozeParametersRequest
		err := mapstructure.Decode(r.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Parameters = params
	case types.ActionTypeChangeState:
		var params scenario.ChangeStateParametersRequest
		err := mapstructure.Decode(r.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Parameters = params
	case types.ActionTypeAssocTicket:
		var params scenario.AssocTicketParametersRequest
		err := mapstructure.Decode(r.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Parameters = params
	case types.ActionTypePbehavior:
		var params scenario.PbehaviorParametersRequest
		err := mapstructure.Decode(r.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Parameters = params
	default:
		var params scenario.ParametersRequest
		err := mapstructure.Decode(r.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Parameters = params
	}

	return nil
}

type AggregationResult struct {
	Data       []idlerule.Rule `bson:"data" json:"data"`
	TotalCount int64           `bson:"total_count" json:"total_count"`
}

// GetTotal implementation PaginatedData interface
func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

// GetData implementation PaginatedData interface
func (r AggregationResult) GetData() interface{} {
	return r.Data
}
