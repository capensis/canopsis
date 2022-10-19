package techmetrics

import "time"

type EventMetric struct {
	Timestamp time.Time
	EventType string
	Interval  time.Duration
}

type CheEventMetric struct {
	EventMetric
	EntityType        string
	IsNewEntity       bool
	IsInfosUpdated    bool
	IsServicesUpdated bool
}

type AxeEventMetric struct {
	EventMetric
	EntityType      string
	AlarmChangeType string
}

type PeriodicalMetric struct {
	Timestamp time.Time
	Interval  time.Duration
}

type AxePeriodicalMetric struct {
	PeriodicalMetric
	Events int64
}

type PbehaviorPeriodicalMetric struct {
	PeriodicalMetric
	Events     int64
	Entities   int64
	Pbehaviors int64
}

type ApiRequestMetric struct {
	Timestamp time.Time
	Interval  time.Duration
	Method    string
	Url       string
}
