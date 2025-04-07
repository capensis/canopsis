package redis

import "time"

const (
	KeyDelimiter = "$$"

	RunInfoKey = "engine-run-info"

	PbehaviorPeriodicalLockKey = "pbehavior-periodical-lock-key"
	RecomputeLockKey           = "pbehavior-recompute-lock-key"
	RecomputeInheritedLockKey  = "pbehavior-recompute-inherited-lock-key"
	RecomputeLockDuration      = 10 * time.Second

	ActionPeriodicalLockKey    = "action-periodical-lock-key"
	ActionDelayedScenarioKey   = "action-delayed-scenario"
	ActionScenarioExecutionKey = "action-scenario-execution"

	AxePeriodicalLockKey                      = "axe-periodical-lock-key"
	AxeInternalTagsPeriodicalLockKey          = "axe-internal-tags-periodical-lock-key"
	AxeEntityServiceStateLockKey              = "axe-entity-service-state-lock-key"
	AxeIdleSincePeriodicalLockKey             = "axe-idle-since-periodical-lock-key"
	AxeNotAckedMetricsPeriodicalLockKey       = "axe-not-acked-metrics-periodical-lock-key"
	AxeSliMetricsPeriodicalLockKey            = "axe-sli-metrics-periodical-lock-key"
	AxeSoftDeleteCorrelationPeriodicalLockKey = "axe-soft-delete-correlation-periodical-lock-key"

	FifoDataStorageLockKey = "fifo-datastorage-lock-key"

	ChePeriodicalLockKey                      = "che-periodical-lock-key"
	CheSoftDeletePeriodicalLockKey            = "che-soft-delete-periodical-lock-key"
	CheEntityInfosDictionaryPeriodicalLockKey = "che-entity-infos-dictionary-periodical-lock-key"
	CheEventFiltersIntervalsPeriodicalLockKey = "che-event-filters-intervals-periodical-lock-key"

	RecorderLockKey = "recorder-lock-key"

	RemediationPeriodicalLockKey        = "remediation-periodical-lock-key"
	RemediationPostponedJobTasksLockKey = "remediation-postponed-job-tasks-lock-key"

	PbehaviorSpanKey              = "pbehavior-span"
	PbehaviorTypesKey             = "pbehavior-types"
	PbehaviorDefaultActiveTypeKey = "pbehavior-default-active-type"
	PbehaviorComputedKey          = "pbehavior-computed-"

	//todo: we should rename this key and do not use "pbehavior-computed-" prefix since it's not a computed pbehavior.
	PbehaviorComputedInheritedServices = "pbehavior-computed-inherited-services"

	DynamicInfosDictionaryPeriodicalLockKey = "dynamic-infos-dictionary-periodical-lock-key"

	ApiCleanEntitiesLockKey      = "api-clean-entities-lock-key"
	ApiUserActivityMetricLockKey = "api-user-activity-metric"
	ApiCacheRequestKey           = "api-cache-request"
	ApiActionLogWatchLockKey     = "api-action-log-watch-lock-key"

	CorrelationInactiveDelayPeriodicalLockKey = "correlation-inactive-delay-periodical-lock-key"
)
