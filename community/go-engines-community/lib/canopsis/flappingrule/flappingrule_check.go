package flappingrule

import (
	"context"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

// Service interface is used to implement alarms modification by idle baggotRules.
type FlappingCheck interface {
	IsFlapping(alarm *types.Alarm) bool
	GetRules() map[string]Rule
	LoadRules(ctx context.Context) error
	engine.PeriodicalWorker
}

type flapping struct {
	ruleAdapter    Adapter
	logger         zerolog.Logger
	changeChannel  chan bool
	reloadInterval time.Duration

	mu    sync.RWMutex
	rules map[string]Rule
}

func SetThenGetFlappingCheck(adapter Adapter, ctx context.Context, reloadInterval time.Duration, zLogger zerolog.Logger) FlappingCheck {
	if reloadInterval == 0 {
		reloadInterval = canopsis.PeriodicalWaitTime
	}

	flapping := &flapping{
		ruleAdapter:    adapter,
		logger:         zLogger,
		changeChannel:  make(chan bool),
		mu:             sync.RWMutex{},
		rules:          make(map[string]Rule),
		reloadInterval: reloadInterval,
	}
	if err := flapping.LoadRules(ctx); err != nil {
		zLogger.Fatal().Err(err)
	}
	types.SetFlapping(flapping.IsFlapping)
	return flapping
}

func (s *flapping) match(alr *types.Alarm) *Rule {
	for _, r := range s.rules {
		if r.Matches(alr) {
			return &r
		}
	}
	return nil
}

func (s *flapping) Work(ctx context.Context) error {
	return s.LoadRules(ctx)
}

func (s *flapping) LoadRules(ctx context.Context) error {
	report := LoadRulesReport{}

	rules, err := s.ruleAdapter.Get(ctx)
	if err != nil {
		return err
	}

	rulesMap := make(map[string]Rule)
	for _, r := range rules {
		rulesMap[r.ID] = r
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Get the rules that have been added or modified since the last time the
	// rules were loaded
	for _, rule := range rulesMap {
		previousValue, alreadyLoaded := s.rules[rule.ID]

		if !alreadyLoaded {
			report.Added = append(report.Added, rule)
		} else if !previousValue.Updated.Equal(rule.Updated.Time) {
			report.Modified = append(report.Modified, rule)
		} else {
			report.Unchanged = append(report.Unchanged, rule)
		}
	}

	// Get the rules that used to be in s.rules, but are not in the database
	// anymore or failed to parse template
	for _, rule := range s.rules {
		_, stillExists := rulesMap[rule.ID]
		if !stillExists {
			report.Removed = append(report.Removed, rule)
		}
	}

	s.rules = rulesMap

	s.logger.Info().
		Int("unchanged", len(report.Unchanged)).
		Int("added", len(report.Added)).
		Int("modified", len(report.Modified)).
		Int("removed", len(report.Removed)).
		Int("failed", len(report.FailParsed)).
		Msg("loaded flapping rules")

	return nil
}

func (s *flapping) GetRules() map[string]Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rules
}

func (s *flapping) IsFlapping(a *types.Alarm) bool {
	s.mu.RLock()
	rule := s.match(a)
	s.mu.RUnlock()
	if rule == nil {
		return false
	}

	return types.IsFlappingWithDurationAndStep(a, rule.FlappingInterval.Duration(), rule.FlappingFreqLimit)
}

func (s *flapping) GetInterval() time.Duration {
	return s.reloadInterval
}
