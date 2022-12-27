package types

// Operation represents alarm modification operation.
type Operation struct {
	Type       string              `bson:"type" json:"type"`
	Parameters OperationParameters `bson:"parameters,omitempty" json:"parameters"`
}

// OperationParameters represents default operation parameters.
type OperationParameters struct {
	Output string `bson:"output,omitempty" json:"output,omitempty"`
	Author string `bson:"author,omitempty" json:"author,omitempty"`
	User   string `bson:"user,omitempty" json:"user,omitempty"`
	// AssocTicket and Webhook
	Ticket            string            `bson:"ticket,omitempty" json:"ticket,omitempty"`
	TicketUrl         string            `bson:"ticket_url,omitempty" json:"ticket_url,omitempty"`
	TicketComment     string            `bson:"ticket_comment,omitempty" json:"ticket_comment,omitempty"`
	TicketSystemName  string            `bson:"ticket_system_name,omitempty" json:"ticket_system_name"`
	TicketMetaAlarmID string            `bson:"ticket_meta_alarm_id,omitempty" json:"ticket_meta_alarm_id"`
	TicketRuleName    string            `bson:"ticket_rule_name,omitempty" json:"ticket_rule_name,omitempty"`
	TicketData        map[string]string `bson:"ticket_data,omitempty" json:"ticket_data,omitempty"`
	// Webhook
	DeclareTicket        bool `bson:"declare_ticket,omitempty" json:"declare_ticket,omitempty"`
	DeclareTicketRequest bool `bson:"declare_ticket_request,omitempty" json:"declare_ticket_request,omitempty"`
	// Snooze
	Duration *DurationWithUnit `bson:"duration,omitempty" json:"duration,omitempty"`
	// ChangeState
	State *CpsNumber `bson:"state,omitempty" json:"state,omitempty"`
	// Pbehavior
	PbehaviorInfo *PbehaviorInfo `bson:"pbehavior_info,omitempty" json:"pbehavior_info,omitempty"`
	// Instruction
	Execution string `bson:"execution,omitempty" json:"execution,omitempty"`
	// Instruction is used only for manual instructions kpi metrics
	Instruction string `bson:"instruction,omitempty" json:"instruction,omitempty"`
}
