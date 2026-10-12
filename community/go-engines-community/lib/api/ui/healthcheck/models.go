package healthcheck

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"

const (
	ServiceMongoDB     = "MongoDB"
	ServiceRedis       = "Redis"
	ServiceRabbitMQ    = "RabbitMQ"
	ServiceTimescaleDB = "TimescaleDB"
)

const (
	StatusHealthy   = "healthy"
	StatusDegraded  = "degraded"
	StatusUnhealthy = "unhealthy"
)

const (
	EngineIssueQueueOverflowed         = "queue_overflowed"
	EngineIssueTooFewInstances         = "too_few_instances"
	EngineIssueInstancesConfigMismatch = "instances_config_mismatch"
)

const (
	ServiceIssueNodeDown       = "node_down"
	ServiceIssueNoQuorum       = "no_quorum"
	ServiceIssueQuorumDegraded = "quorum_degraded"
)

type Info struct {
	Services []Service `json:"services"`
	Engines  Engines   `json:"engines"`

	MaxQueueLength    int `json:"max_queue_length"`
	MaxMessagesLength int `json:"max_messages_length"`

	HasInvalidEnginesOrder bool `json:"has_invalid_engines_order"`
}

type Engines struct {
	Graph      Graph             `json:"graph"`
	Parameters map[string]Engine `json:"parameters"`
}

type Status struct {
	Services []Service `json:"services"`
	Engines  []Engine  `json:"engines"`

	HasInvalidEnginesOrder bool `json:"has_invalid_engines_order"`
}

type Service struct {
	Name         string   `json:"name"`
	Status       string   `json:"status"`
	Issues       []string `json:"issues,omitempty"`
	Nodes        *int     `json:"nodes,omitempty"`
	RunningNodes *int     `json:"running_nodes,omitempty"`
}

type Graph struct {
	Nodes []string `json:"nodes"`
	Edges []Edge   `json:"edges"`
}

type Engine struct {
	Name             string            `json:"name,omitempty"`
	Status           string            `json:"status"`
	Issues           []string          `json:"issues,omitempty"`
	Instances        *int              `json:"instances,omitempty"`
	MinInstances     *int              `json:"min_instances,omitempty"`
	OptimalInstances *int              `json:"optimal_instances,omitempty"`
	QueueLength      *int              `json:"queue_length,omitempty"`
	Time             *datetime.CpsTime `json:"time,omitempty" swaggertype:"integer"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type LiveResponse struct {
	Ok bool `json:"ok"`
}
