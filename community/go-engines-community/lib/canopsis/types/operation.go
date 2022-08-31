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
	// AssocTicket and DeclareTicket
	Ticket string `bson:"ticket,omitempty" json:"ticket,omitempty"`
	// DeclareTicket
	Data map[string]string `bson:"data,omitempty" json:"data,omitempty"`
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
