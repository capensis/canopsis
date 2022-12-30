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
	TicketInfo `bson:",inline"`
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
