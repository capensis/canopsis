package types

// Operation represents alarm modification operation.
type Operation struct {
	Type       string      `bson:"type"`
	Parameters interface{} `bson:"parameters,omitempty"`
}

// OperationParameters represents default operation parameters.
type OperationParameters struct {
	Output string `bson:"output" json:"output"`
	Author string `bson:"author" json:"author"`
	User   string `bson:"user" json:"user"`
}

type OperationAssocTicketParameters struct {
	Ticket string `bson:"ticket" json:"ticket"`
	Output string `bson:"output" json:"output"`
	Author string `bson:"author" json:"author"`
	User   string `bson:"user" json:"user"`
}

type OperationSnoozeParameters struct {
	Duration DurationWithUnit `bson:"duration" json:"duration"`
	Output   string           `bson:"output" json:"output"`
	Author   string           `bson:"author" json:"author"`
	User     string           `bson:"user" json:"user"`
}

type OperationChangeStateParameters struct {
	State  CpsNumber `bson:"state" json:"state"`
	Output string    `bson:"output" json:"output"`
	Author string    `bson:"author" json:"author"`
	User   string    `bson:"user" json:"user"`
}

type OperationDeclareTicketParameters struct {
	Ticket string            `bson:"ticket" json:"ticket"`
	Data   map[string]string `bson:"data" json:"data"`
	Output string            `bson:"output" json:"output"`
	Author string            `bson:"author" json:"author"`
	User   string            `bson:"user" json:"user"`
}

type OperationPbhParameters struct {
	PbehaviorInfo PbehaviorInfo `json:"pbehavior_info"`
	Output        string        `json:"output"`
	Author        string        `json:"author"`
	User          string        `bson:"user" json:"user"`
}

type OperationInstructionParameters struct {
	Execution string `bson:"execution" json:"execution"`
	Output    string `bson:"output" json:"output"`
	Author    string `bson:"author" json:"author"`
	User      string `bson:"user" json:"user"`
}
