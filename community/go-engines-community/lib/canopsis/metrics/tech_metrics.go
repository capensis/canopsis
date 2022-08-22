package metrics

import "time"

type FifoEventMetric struct {
	Timestamp time.Time
	EventType string
	Interval  int64
}
