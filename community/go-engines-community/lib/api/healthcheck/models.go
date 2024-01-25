package healthcheck

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"

const (
	ServiceMongoDB     = "MongoDB"
	ServiceRedis       = "Redis"
	ServiceRabbitMQ    = "RabbitMQ"
	ServiceTimescaleDB = "TimescaleDB"
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
	Name      string `json:"name"`
	IsRunning bool   `json:"is_running"`
}

type Graph struct {
	Nodes []string `json:"nodes"`
	Edges []Edge   `json:"edges"`
}

type Engine struct {
	Name             string            `json:"name,omitempty"`
	Instances        *int              `json:"instances,omitempty"`
	MinInstances     *int              `json:"min_instances,omitempty"`
	OptimalInstances *int              `json:"optimal_instances,omitempty"`
	QueueLength      *int              `json:"queue_length,omitempty"`
	Time             *datetime.CpsTime `json:"time,omitempty" swaggertype:"integer"`

	IsRunning             bool `json:"is_running"`
	IsQueueOverflown      bool `json:"is_queue_overflown"`
	IsTooFewInstances     bool `json:"is_too_few_instances"`
	IsDiffInstancesConfig bool `json:"is_diff_instances_config"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type LiveResponse struct {
	Ok bool `json:"ok"`
}
