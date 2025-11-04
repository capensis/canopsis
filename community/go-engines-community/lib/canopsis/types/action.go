package types

const (
	ActionTypeAck             = EventTypeAck
	ActionTypeAckRemove       = EventTypeAckremove
	ActionTypeAssocTicket     = EventTypeAssocTicket
	ActionTypeCancel          = EventTypeCancel
	ActionTypeChangeState     = EventTypeChangestate
	ActionTypeSnooze          = EventTypeSnooze
	ActionTypeUnsnooze        = EventTypeUnsnooze
	ActionTypePbehavior       = "pbehavior"
	ActionTypePbehaviorRemove = "pbehaviorremove"
	ActionTypeWebhook         = "webhook"
)

type AdditionalData struct {
	Trigger   string `json:"trigger"`
	Author    string `json:"author"`
	User      string `json:"user"`
	Initiator string `json:"initiator"`
	Output    string `json:"event_output"`
	RuleName  string `json:"rule_name"`

	// Deprecated: use Trigger instead of AlarmChangeType
	AlarmChangeType string `json:"alarm_change_type"`
}
