package redis

import "time"

const (
	PeriodicalLockKey     = "pbehavior-periodical-lock-key"
	RecomputeLockKey      = "pbehavior-recompute-lock-key"
	RecomputeLockDuration = 10 * time.Second

	DelayedScenarioKey   = "delayed-scenario"
	ScenarioExecutionKey = "scenario-execution"

	PbehaviorSpanKey              = "pbehavior-span"
	PbehaviorTypesKey             = "pbehavior-types"
	PbehaviorDefaultActiveTypeKey = "pbehavior-default-active-type"
	PbehaviorComputedKey          = "pbehavior-computed-"
	PbehaviorEntityMatchKey       = "pbehavior-entity-match-"
)
