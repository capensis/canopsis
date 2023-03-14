package redis

import "time"

const (
	KeyDelimiter = "$$"

	RunInfoKey = "engine-run-info"

	PbehaviorPeriodicalLockKey      = "pbehavior-periodical-lock-key"
	PbehaviorCleanPeriodicalLockKey = "pbehavior-clean-periodical-lock-key"
	RecomputeLockKey                = "pbehavior-recompute-lock-key"
	RecomputeLockDuration           = 10 * time.Second

	ActionPeriodicalLockKey    = "action-periodical-lock-key"
	ActionDelayedScenarioKey   = "action-delayed-scenario"
	ActionScenarioExecutionKey = "action-scenario-execution"

	AxePeriodicalLockKey                 = "axe-periodical-lock-key"
	AxeResolvedArchiverPeriodicalLockKey = "axe-resolved-archiver-periodical-lock-key"

	FifoDeleteOutdatedRatesLockKey = "fifo-delete-outdated-rates-lock-key"

	ChePeriodicalLockKey           = "che-periodical-lock-key"
	CheSoftDeletePeriodicalLockKey = "che-soft-delete-periodical-lock-key"

	ServicePeriodicalLockKey          = "service-periodical-lock"
	ServiceIdleSincePeriodicalLockKey = "service-periodical-idle-since-lock"

	RemediationPeriodicalLockKey      = "remediation-periodical-lock-key"
	RemediationStatsPeriodicalLockKey = "remediation-stats-periodical-lock-key"

	PbehaviorSpanKey              = "pbehavior-span"
	PbehaviorTypesKey             = "pbehavior-types"
	PbehaviorDefaultActiveTypeKey = "pbehavior-default-active-type"
	PbehaviorComputedKey          = "pbehavior-computed-"

	CheEntityInfosDictionaryPeriodicalLockKey = "che-entity-infos-dictionary-periodical-lock-key"
	DynamicInfosDictionaryPeriodicalLockKey   = "dynamic-infos-dictionary-periodical-lock-key"

	CheEventFiltersIntervalsPeriodicalLockKey = "che-event-filters-intervals-periodical-lock-key"

	MetricsNotAckedPeriodicalLockKey = "metrics-not-acked-periodical-lock-key"
)
