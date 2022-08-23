package metrics

import "time"

type FifoEventMetric struct {
	Timestamp time.Time
	EventType string
	Interval  int64
}

type CheEventMetric struct {
	Timestamp         time.Time
	EventType         string
	Interval          int64
	EntityType        string
	IsNewEntity       bool
	IsInfosUpdated    bool
	IsServicesUpdated bool
}
