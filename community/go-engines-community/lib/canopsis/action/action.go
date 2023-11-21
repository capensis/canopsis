package action

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Scenario struct {
	ID                   string                     `bson:"_id,omitempty" json:"_id,omitempty"`
	Name                 string                     `bson:"name" json:"name"`
	Author               string                     `bson:"author" json:"author"`
	Enabled              bool                       `bson:"enabled" json:"enabled"`
	DisableDuringPeriods []string                   `bson:"disable_during_periods" json:"disable_during_periods"`
	Triggers             []string                   `bson:"triggers" json:"triggers"`
	Actions              []Action                   `bson:"actions" json:"actions"`
	Priority             int64                      `bson:"priority" json:"priority"`
	Delay                *datetime.DurationWithUnit `bson:"delay" json:"delay"`
	Created              datetime.CpsTime           `bson:"created,omitempty" json:"created,omitempty"`
	Updated              datetime.CpsTime           `bson:"updated,omitempty" json:"updated,omitempty"`
}

func (s Scenario) IsTriggered(triggers []string) string {
	for _, expectedTrigger := range s.Triggers {
		for _, trigger := range triggers {
			if expectedTrigger == trigger {
				return trigger
			}
		}
	}

	return ""
}

// Action represents a canopsis Action on alarms.
type Action struct {
	Type                     string     `bson:"type" json:"type"`
	Comment                  string     `bson:"comment" json:"comment"`
	Parameters               Parameters `bson:"parameters,omitempty" json:"parameters,omitempty"`
	DropScenarioIfNotMatched bool       `bson:"drop_scenario_if_not_matched" json:"drop_scenario_if_not_matched"`
	EmitTrigger              bool       `bson:"emit_trigger" json:"emit_trigger"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}

// TODO: change arguments to pointers, check any other places where changing event, alarm and entity structs to pointers gives a performance boost.
func (a Action) Match(entity types.Entity, alarm types.Alarm) (bool, error) {
	return match.Match(&entity, &alarm, a.EntityPattern, a.AlarmPattern)
}

type Parameters struct {
	Output string `json:"output,omitempty" bson:"output,omitempty" binding:"max=255"`

	ForwardAuthor *bool  `json:"forward_author,omitempty" bson:"forward_author,omitempty"`
	Author        string `json:"author,omitempty" bson:"author,omitempty"`

	// State is used in changestate action.
	//   * `0` - Info
	//   * `1` - Minor
	//   * `2` - Major
	//   * `3` - Critical
	State *types.CpsNumber `json:"state,omitempty" bson:"state,omitempty"`
	// Ticket is used in assocticket action.
	Ticket string `json:"ticket,omitempty" binding:"max=255" bson:"ticket,omitempty"`
	// TicketURL is used in assocticket action.
	TicketURL string `json:"ticket_url,omitempty" binding:"max=255" bson:"ticket_url,omitempty"`
	// TicketSystemName is used in assocticket and webhook action.
	TicketSystemName string `json:"ticket_system_name,omitempty" binding:"max=255" bson:"ticket_system_name,omitempty"`
	// TicketData is used in assocticket action.
	TicketData map[string]string `json:"ticket_data,omitempty" bson:"ticket_data,omitempty"`
	// Duration is used in snooze and pbehavior actions.
	Duration *datetime.DurationWithUnit `json:"duration,omitempty" bson:"duration,omitempty"`
	// Name is used in pbehavior action.
	Name string `json:"name,omitempty" binding:"max=255" bson:"name,omitempty"`
	// Reason is used in pbehavior action.
	Reason string `json:"reason,omitempty" bson:"reason,omitempty"`
	// Type is used in pbehavior action.
	Type string `json:"type,omitempty" bson:"type,omitempty"`
	// RRule is used in pbehavior action.
	RRule string `json:"rrule,omitempty" bson:"rrule,omitempty"`
	// Tstart is used in pbehavior action.
	Tstart *datetime.CpsTime `json:"tstart,omitempty" bson:"tstart,omitempty" swaggertype:"integer"`
	// Tstop is used in pbehavior action.
	Tstop *datetime.CpsTime `json:"tstop,omitempty" bson:"tstop,omitempty" swaggertype:"integer"`
	// StartOnTrigger is used in pbehavior action.
	StartOnTrigger *bool `json:"start_on_trigger,omitempty" bson:"start_on_trigger,omitempty"`
	// Request is used in webhook action.
	Request *request.Parameters `json:"request,omitempty" bson:"request,omitempty"`
	// SkipForChild is used in webhook action.
	SkipForChild *bool `json:"skip_for_child,omitempty" bson:"skip_for_child,omitempty"`
	// SkipForInstruction is used in webhook action.
	SkipForInstruction *bool `json:"skip_for_instruction,omitempty" bson:"skip_for_instruction,omitempty"`
	// DeclareTicket is used in webhook action.
	DeclareTicket *request.WebhookDeclareTicket `json:"declare_ticket,omitempty" bson:"declare_ticket,omitempty"`
}
