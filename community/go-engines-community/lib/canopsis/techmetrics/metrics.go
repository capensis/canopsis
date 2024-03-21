package techmetrics

import "time"

type EventMetric struct {
	Timestamp time.Time
	EventType string
	Interval  time.Duration
}

type CpsEventMetric struct {
	EventMetric
	IsOkState *bool
}

type FifoEventMetric struct {
	EventMetric
	ExternalRequests map[string]int64
}

type CheEventMetric struct {
	EventMetric
	EntityType            string
	IsNewEntity           bool
	IsInfosUpdated        bool
	IsServicesUpdated     bool
	IsStateSettingUpdated bool
	ExecutedEnrichRules   int64
	ExternalRequests      map[string]int64
}

type AxeEventMetric struct {
	EventMetric
	EntityType        string
	AlarmChangeType   string
	IsOkState         *bool
	IsCountersUpdated bool
}

type CorrelationEventMetric struct {
	EventMetric
	MatchedRuleCount int64
	MatchedRuleTypes []string
}

type DynamicInfoEventMetric struct {
	EventMetric
	ExecutedRules int64
}

type ActionEventMetric struct {
	EventMetric
	ExecutedRules    int64
	ExecutedWebhooks int64
}

type PeriodicalMetric struct {
	Timestamp time.Time
	Interval  time.Duration
}

type AxePeriodicalMetric struct {
	PeriodicalMetric
	Events     int64
	IdleEvents int64
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

type CorrelationRetriesMetric struct {
	Timestamp time.Time
	Type      string
	Retries   int
}
