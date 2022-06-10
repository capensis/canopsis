package action

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
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
	Type                     string                       `bson:"type" json:"type"`
	Comment                  string                       `bson:"comment" json:"comment"`
	Parameters               Parameters                   `bson:"parameters,omitempty" json:"parameters,omitempty"`
	OldAlarmPatterns         oldpattern.AlarmPatternList  `bson:"old_alarm_patterns,omitempty" json:"old_alarm_patterns,omitempty"`
	OldEntityPatterns        oldpattern.EntityPatternList `bson:"old_entity_patterns,omitempty" json:"old_entity_patterns,omitempty"`
	DropScenarioIfNotMatched bool                         `bson:"drop_scenario_if_not_matched" json:"drop_scenario_if_not_matched"`
	EmitTrigger              bool                         `bson:"emit_trigger" json:"emit_trigger"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}

func (a Action) Match(entity types.Entity, alarm types.Alarm) (bool, error) {
	if !a.OldAlarmPatterns.IsSet() && !a.OldEntityPatterns.IsSet() && len(a.EntityPattern) == 0 && len(a.AlarmPattern) == 0 {
		return false, nil
	}

	var matched bool
	var err error

	if a.OldAlarmPatterns.IsSet() {
		if !a.OldAlarmPatterns.IsValid() {
			return false, InvalidOldAlarmPattern
		}

		matched = a.OldAlarmPatterns.Matches(&alarm)
	} else {
		matched, err = a.AlarmPattern.Match(alarm)
		if err != nil {
			return false, AlarmPatternError
		}
	}

	if !matched {
		return false, nil
	}

	if a.OldEntityPatterns.IsSet() {
		if !a.OldEntityPatterns.IsValid() {
			return false, InvalidOldEntityPattern
		}

		matched = a.OldEntityPatterns.Matches(&entity)
	} else {
		matched, _, err = a.EntityPattern.Match(entity)
		if err != nil {
			return false, EntityPatternError
		}
	}

	return matched, nil
}

type Parameters struct {
	Output string `json:"output" bson:"output,omitempty" binding:"max=255"`

	ForwardAuthor *bool  `json:"forward_author" bson:"forward_author,omitempty"`
	Author        string `json:"author" bson:"author,omitempty"`

	// State is used in changestate action.
	//   * `0` - Info
	//   * `1` - Minor
	//   * `2` - Major
	//   * `3` - Critical
	State *types.CpsNumber `json:"state" bson:"state,omitempty"`
	// Ticket is used in assocticket action.
	Ticket string `json:"ticket" binding:"max=255" bson:"ticket,omitempty"`
	// Duration is used in snooze and pbehavior actions.
	Duration *types.DurationWithUnit `json:"duration" bson:"duration,omitempty"`
	// Name is used in pbehavior action.
	Name string `json:"name" binding:"max=255" bson:"name,omitempty"`
	// Reason is used in pbehavior action.
	Reason string `json:"reason" bson:"reason,omitempty"`
	// Type is used in pbehavior action.
	Type string `json:"type" bson:"type,omitempty"`
	// RRule is used in pbehavior action.
	RRule string `json:"rrule" bson:"rrule,omitempty"`
	// Tstart is used in pbehavior action.
	Tstart *types.CpsTime `json:"tstart" bson:"tstart,omitempty" swaggertype:"integer"`
	// Tstop is used in pbehavior action.
	Tstop *types.CpsTime `json:"tstop" bson:"tstop,omitempty" swaggertype:"integer"`
	// StartOnTrigger is used in pbehavior action.
	StartOnTrigger *bool `json:"start_on_trigger" bson:"start_on_trigger,omitempty"`
	// Request is used in webhook action.
	Request *types.WebhookRequest `json:"request" bson:"request,omitempty"`
	// DeclareTicket is used in webhook action.
	DeclareTicket *types.WebhookDeclareTicket `json:"declare_ticket" bson:"declare_ticket,omitempty"`
	// RetryCount is used in webhook action.
	RetryCount int64 `json:"retry_count" bson:"retry_count,omitempty" binding:"min=0"`
	// RetryDelay is used in webhook action.
	RetryDelay *types.DurationWithUnit `json:"retry_delay" bson:"retry_delay,omitempty"`
}
