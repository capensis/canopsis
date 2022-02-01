package types

const (
	ActionTypeAck         = EventTypeAck
	ActionTypeAckRemove   = EventTypeAckremove
	ActionTypeAssocTicket = EventTypeAssocTicket
	ActionTypeCancel      = EventTypeCancel
	ActionTypeChangeState = EventTypeChangestate
	ActionTypeSnooze      = EventTypeSnooze
	ActionTypePbehavior   = "pbehavior"
	ActionTypeWebhook     = "webhook"
)

type ActionPBehaviorParameters struct {
	Author         string            `bson:"author" json:"author"`
	UserID         string            `bson:"user" json:"user"`
	Name           string            `bson:"name" json:"name"`
	Reason         string            `bson:"reason" json:"reason"`
	Type           string            `bson:"type" json:"type"`
	RRule          string            `bson:"rrule" json:"rrule"`
	Tstart         *int64            `bson:"tstart,omitempty" json:"tstart,omitempty"`
	Tstop          *int64            `bson:"tstop,omitempty" json:"tstop,omitempty"`
	StartOnTrigger *bool             `bson:"start_on_trigger,omitempty" json:"start_on_trigger,omitempty" mapstructure:"start_on_trigger,omitempty"`
	Duration       *DurationWithUnit `bson:"duration,omitempty" json:"duration,omitempty"`
}
