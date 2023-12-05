package idlerule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id name author.name created updated type priority"`
}

type EditRequest struct {
	Name                 string                    `json:"name" binding:"required,max=255"`
	Description          string                    `json:"description" binding:"max=255"`
	Author               string                    `json:"author" swaggerignore:"true"`
	Enabled              *bool                     `json:"enabled" binding:"required"`
	Type                 string                    `json:"type" binding:"required"`
	Priority             int64                     `json:"priority" binding:"min=0"`
	Duration             datetime.DurationWithUnit `json:"duration" binding:"required"`
	Comment              string                    `json:"comment"`
	DisableDuringPeriods []string                  `json:"disable_during_periods"`
	AlarmCondition       string                    `json:"alarm_condition"`
	Operation            *OperationRequest         `json:"operation,omitempty"`

	common.AlarmPatternFieldsRequest
	common.EntityPatternFieldsRequest
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

type OperationRequest struct {
	Type       string              `json:"type" binding:"required"`
	Parameters idlerule.Parameters `json:"parameters,omitempty"`
}

type Rule struct {
	ID          string                    `bson:"_id,omitempty" json:"_id"`
	Name        string                    `bson:"name" json:"name"`
	Description string                    `bson:"description" json:"description"`
	Author      *author.Author            `bson:"author" json:"author"`
	Enabled     bool                      `bson:"enabled" json:"enabled"`
	Type        string                    `bson:"type" json:"type"`
	Priority    int64                     `bson:"priority" json:"priority"`
	Duration    datetime.DurationWithUnit `bson:"duration" json:"duration"`
	Comment     string                    `bson:"comment" json:"comment"`
	// DisableDuringPeriods is an option that allows to disable the rule
	// when entity is in listed periods due pbehavior schedule.
	DisableDuringPeriods []string         `bson:"disable_during_periods" json:"disable_during_periods"`
	Created              datetime.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Updated              datetime.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`
	// Only for Alarm rules
	AlarmCondition string     `bson:"alarm_condition,omitempty" json:"alarm_condition,omitempty"`
	Operation      *Operation `bson:"operation,omitempty" json:"operation,omitempty"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}

type Operation struct {
	Type       string     `bson:"type" json:"type"`
	Parameters Parameters `bson:"parameters,omitempty" json:"parameters"`
}

type Parameters struct {
	Output string `json:"output,omitempty" bson:"output"`
	// ChangeState
	State *types.CpsNumber `json:"state,omitempty" bson:"state"`
	// AssocTicket
	Ticket           string            `json:"ticket,omitempty" bson:"ticket"`
	TicketURL        string            `json:"ticket_url,omitempty" bson:"ticket_url"`
	TicketSystemName string            `json:"ticket_system_name,omitempty" bson:"ticket_system_name"`
	TicketData       map[string]string `json:"ticket_data,omitempty" bson:"ticket_data"`
	// Snooze and Pbehavior
	Duration *datetime.DurationWithUnit `json:"duration,omitempty" bson:"duration"`
	// Pbehavior
	Name           string            `json:"name,omitempty" binding:"max=255" bson:"name"`
	Reason         *pbehavior.Reason `json:"reason,omitempty" bson:"reason"`
	Type           *pbehavior.Type   `json:"type,omitempty" bson:"type"`
	RRule          string            `json:"rrule,omitempty" bson:"rrule"`
	Tstart         *int64            `json:"tstart,omitempty" bson:"tstart"`
	Tstop          *int64            `json:"tstop,omitempty" bson:"tstop"`
	StartOnTrigger *bool             `json:"start_on_trigger,omitempty" bson:"start_on_trigger"`
}

type AggregationResult struct {
	Data       []Rule `bson:"data" json:"data"`
	TotalCount int64  `bson:"total_count" json:"total_count"`
}

// GetTotal implementation PaginatedData interface
func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

// GetData implementation PaginatedData interface
func (r AggregationResult) GetData() interface{} {
	return r.Data
}
