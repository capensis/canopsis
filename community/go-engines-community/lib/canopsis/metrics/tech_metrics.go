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

type AxeEventMetric struct {
	Timestamp       time.Time
	EventType       string
	Interval        int64
	EntityType      string
	AlarmChangeType string
}

type AxePeriodicalMetric struct {
	Timestamp time.Time
	Interval  int64
}

type ApiRequestMetric struct {
	Timestamp time.Time
	Interval  int64
}

type SimpleEventMetric struct {
	Timestamp time.Time
	EventType string
	Interval  int64
}
