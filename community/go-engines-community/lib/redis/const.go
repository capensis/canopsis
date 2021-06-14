package redis

import "time"

const (
	PeriodicalLockKey     = "pbehavior-periodical-lock-key"
	RecomputeLockKey      = "pbehavior-recompute-lock-key"
	RecomputeLockDuration = 10 * time.Second

	DelayedScenarioKey   = "delayed-scenario"
	ScenarioExecutionKey = "scenario-execution"
)
