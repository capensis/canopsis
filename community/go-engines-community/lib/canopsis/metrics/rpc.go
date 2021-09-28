package metrics

import "github.com/prometheus/client_golang/prometheus"

const TestType = "test"

type RPCMetricEvent struct {
	Type   string            `json:"type"`
	Labels prometheus.Labels `json:"labels"`
}
