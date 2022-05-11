package action

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	PriorityField = "priority"
	IdField       = "_id"
)

type Scenario struct {
	ID                   string                  `bson:"_id,omitempty" json:"_id,omitempty"`
	Name                 string                  `bson:"name" json:"name"`
	Author               string                  `bson:"author" json:"author"`
	Enabled              bool                    `bson:"enabled" json:"enabled"`
	DisableDuringPeriods []string                `bson:"disable_during_periods" json:"disable_during_periods"`
	Triggers             []string                `bson:"triggers" json:"triggers"`
	Actions              []Action                `bson:"actions" json:"actions"`
	Priority             int                     `bson:"priority" json:"priority"`
	Delay                *types.DurationWithUnit `bson:"delay" json:"delay"`
	Created              types.CpsTime           `bson:"created,omitempty" json:"created,omitempty"`
	Updated              types.CpsTime           `bson:"updated,omitempty" json:"updated,omitempty"`
}

func (s Scenario) IsTriggered(triggers []string) bool {
	for _, expectedTrigger := range s.Triggers {
		for _, trigger := range triggers {
			if expectedTrigger == trigger {
				return true
			}
		}
	}

	return false
}

// Action represents a canopsis Action on alarms.
type Action struct {
	Type                     string                    `bson:"type" json:"type"`
	Comment                  string                    `bson:"comment" json:"comment"`
	Parameters               Parameters                `bson:"parameters,omitempty" json:"parameters,omitempty"`
	AlarmPatterns            pattern.AlarmPatternList  `bson:"alarm_patterns" json:"alarm_patterns"`
	EntityPatterns           pattern.EntityPatternList `bson:"entity_patterns" json:"entity_patterns"`
	DropScenarioIfNotMatched bool                      `bson:"drop_scenario_if_not_matched" json:"drop_scenario_if_not_matched"`
	EmitTrigger              bool                      `bson:"emit_trigger" json:"emit_trigger"`
}

type Parameters struct {
	Output string `json:"output" bson:"output,omitempty" binding:"max=255"`

	ForwardAuthor *bool  `json:"forward_author" bson:"forward_author,omitempty"`
	Author        string `json:"author" bson:"author,omitempty"`

	// ChangeState
	State *types.CpsNumber `json:"state" bson:"state,omitempty"`
	// AssocTicket
	Ticket string `json:"ticket" binding:"max=255" bson:"ticket,omitempty"`
	// Snooze and Pbehavior
	Duration *types.DurationWithUnit `json:"duration" bson:"duration,omitempty"`
	// Pbehavior
	Name           string         `json:"name" binding:"max=255" bson:"name,omitempty"`
	Reason         string         `json:"reason" bson:"reason,omitempty"`
	Type           string         `json:"type" bson:"type,omitempty"`
	RRule          string         `json:"rrule" bson:"rrule,omitempty"`
	Tstart         *types.CpsTime `json:"tstart" bson:"tstart,omitempty"`
	Tstop          *types.CpsTime `json:"tstop" bson:"tstop,omitempty"`
	StartOnTrigger *bool          `json:"start_on_trigger" bson:"start_on_trigger,omitempty"`
	// Webhook
	Request       *types.WebhookRequest       `json:"request" bson:"request,omitempty"`
	DeclareTicket *types.WebhookDeclareTicket `json:"declare_ticket" bson:"declare_ticket,omitempty"`
	RetryCount    int64                       `json:"retry_count" bson:"retry_count,omitempty" binding:"min=0"`
	RetryDelay    *types.DurationWithUnit     `json:"retry_delay" bson:"retry_delay,omitempty"`
}
