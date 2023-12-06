package action

import (
	"context"
	"errors"
	"fmt"
	"time"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type DelayedScenarioManager interface {
	AddDelayedScenario(context.Context, types.Alarm, Scenario, AdditionalData) error
	PauseDelayedScenarios(context.Context, types.Alarm) error
	ResumeDelayedScenarios(context.Context, types.Alarm) error
	Run(context.Context) (<-chan DelayedScenarioTask, error)
}

func NewDelayedScenarioManager(
	adapter Adapter,
	alarmAdapter libalarm.Adapter,
	storage DelayedScenarioStorage,
	periodicalTimeout time.Duration,
	logger zerolog.Logger,
) DelayedScenarioManager {
	return &delayedScenarioManager{
		scenarioAdapter:   adapter,
		alarmAdapter:      alarmAdapter,
		storage:           storage,
		periodicalTimeout: periodicalTimeout,
		logger:            logger,
	}
}

type delayedScenarioManager struct {
	scenarioAdapter   Adapter
	alarmAdapter      libalarm.Adapter
	storage           DelayedScenarioStorage
	periodicalTimeout time.Duration
	runCh             chan<- DelayedScenarioTask
	logger            zerolog.Logger
}

type DelayedScenarioTask struct {
	Alarm    types.Alarm
	Scenario Scenario

	AdditionalData AdditionalData
}

func (m *delayedScenarioManager) AddDelayedScenario(ctx context.Context, alarm types.Alarm, scenario Scenario, additionalData AdditionalData) error {
	if scenario.Delay == nil || scenario.Delay.Value == 0 {
		return errors.New("scenario is not delayed")
	}

	now := datetime.NewCpsTime()
	delayedScenario := DelayedScenario{
		ScenarioID:     scenario.ID,
		ExecutionTime:  scenario.Delay.AddTo(now),
		AlarmID:        alarm.ID,
		AdditionalData: additionalData,
	}
	id, err := m.storage.Add(ctx, delayedScenario)
	if err != nil {
		return err
	}

	delayedScenario.ID = id
	m.logger.Debug().Str("scenario", scenario.ID).Str("alarm", alarm.ID).Time("timeout_expiration", delayedScenario.ExecutionTime.Time).Msg("start timeout of delayed scenario")

	if m.canStartWaitGoroutine(delayedScenario.ExecutionTime.Time) {
		go m.waitAlmostExpiredTimeoutScenario(ctx, delayedScenario)
	}

	return nil
}

func (m *delayedScenarioManager) PauseDelayedScenarios(ctx context.Context, alarm types.Alarm) error {
	delayedScenarios, err := m.storage.GetAll(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, delayedScenario := range delayedScenarios {
		if delayedScenario.AlarmID == alarm.ID && !delayedScenario.Paused {
			delayedScenario.Paused = true
			delayedScenario.TimeLeft = delayedScenario.ExecutionTime.Sub(now)
			_, err := m.storage.Update(ctx, delayedScenario)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *delayedScenarioManager) ResumeDelayedScenarios(ctx context.Context, alarm types.Alarm) error {
	delayedScenarios, err := m.storage.GetAll(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, delayedScenario := range delayedScenarios {
		if delayedScenario.AlarmID == alarm.ID && delayedScenario.Paused {
			delayedScenario.Paused = false
			delayedScenario.ExecutionTime = datetime.CpsTime{Time: now.Add(delayedScenario.TimeLeft)}
			delayedScenario.TimeLeft = 0
			_, err := m.storage.Update(ctx, delayedScenario)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *delayedScenarioManager) Run(parentCtx context.Context) (<-chan DelayedScenarioTask, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	ch := make(chan DelayedScenarioTask)
	m.runCh = ch

	go func() {
		defer func() {
			cancel()
			close(ch)
		}()

		m.checkExpiredTimeoutScenario(ctx, ch)

		ticker := time.NewTicker(m.periodicalTimeout)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.checkExpiredTimeoutScenario(ctx, ch)
			}
		}
	}()

	return ch, nil
}

func (m *delayedScenarioManager) checkExpiredTimeoutScenario(ctx context.Context, ch chan<- DelayedScenarioTask) {
	expired, almostExpired, err := m.getDelayedScenarios(ctx)
	if err != nil {
		m.logger.Err(err).Msg("couldn't get delayed scenarios")
		return
	}

	tasks, err := m.getExpiredTimeoutScenarios(ctx, expired)
	if err != nil {
		m.logger.Err(err).Msg("couldn't resolve expired delayed scenarios")
		return
	}

	for _, t := range tasks {
		ch <- t
	}

	for _, delayedScenario := range almostExpired {
		go m.waitAlmostExpiredTimeoutScenario(ctx, delayedScenario)
	}

	m.logger.Debug().Int("expired", len(expired)).Int("almost_expired", len(almostExpired)).Msg("checked expired timeout delayed actions")
}

func (m *delayedScenarioManager) waitAlmostExpiredTimeoutScenario(ctx context.Context, scenario DelayedScenario) {
	if m.runCh == nil {
		return
	}

	m.logger.Debug().Str("execution", scenario.ID).Str("scenario", scenario.ScenarioID).Time("timeout_expiration", scenario.ExecutionTime.Time).Msg("wait almost expired timeout")
	defer m.logger.Debug().Str("execution", scenario.ID).Str("scenario", scenario.ScenarioID).Time("timeout_expiration", scenario.ExecutionTime.Time).Msg("stop wait almost expired timeout")

	now := time.Now()
	select {
	case <-ctx.Done():
	case <-time.After(scenario.ExecutionTime.Sub(now)):
		updatedScenario, err := m.storage.Get(ctx, scenario.ID)
		if err != nil {
			m.logger.Err(err).Msg("failed to get delayed scenario")
			return
		}

		if updatedScenario == nil || updatedScenario.Paused || time.Now().Unix() < updatedScenario.ExecutionTime.Unix() {
			m.logger.Debug().Str("execution", scenario.ID).Time("timeout_expiration", scenario.ExecutionTime.Time).Msg("delayed scenario was deleted or paused")
			return
		}

		ok, err := m.storage.Delete(ctx, updatedScenario.ID)
		if err != nil {
			m.logger.Err(err).Msg("failed to delete delayed scenario")
			return
		}

		if !ok {
			return
		}

		tasks, err := m.getExpiredTimeoutScenarios(ctx, []DelayedScenario{*updatedScenario})
		if err != nil {
			m.logger.Err(err).Msg("failed to load delayed scenario")
			return
		}

		for _, t := range tasks {
			m.runCh <- t
			m.logger.Debug().Str("scenario", t.Scenario.ID).Str("alarm", t.Alarm.ID).Msg("expire timeout delayed scenario")
		}
	}
}

func (m *delayedScenarioManager) getDelayedScenarios(ctx context.Context) (
	[]DelayedScenario,
	[]DelayedScenario,
	error,
) {
	now := time.Now()
	expired := make([]DelayedScenario, 0)
	almostExpired := make([]DelayedScenario, 0)
	delayedScenarios, err := m.storage.GetAll(ctx)
	if err != nil {
		return nil, nil, err
	}

	for _, delayedScenario := range delayedScenarios {
		if delayedScenario.Paused {
			continue
		}

		if now.Unix() >= delayedScenario.ExecutionTime.Unix() {
			ok, err := m.storage.Delete(ctx, delayedScenario.ID)
			if err != nil {
				return nil, nil, err
			}

			if ok {
				expired = append(expired, delayedScenario)
			}
		} else if m.canStartWaitGoroutine(delayedScenario.ExecutionTime.Time) {
			almostExpired = append(almostExpired, delayedScenario)
		}
	}

	return expired, almostExpired, nil
}

func (m *delayedScenarioManager) getExpiredTimeoutScenarios(
	ctx context.Context,
	delayedScenarios []DelayedScenario,
) (
	[]DelayedScenarioTask,
	error,
) {
	if len(delayedScenarios) == 0 {
		return nil, nil
	}

	scenarioIDs := make([]string, len(delayedScenarios))
	alarmIDs := make([]string, len(delayedScenarios))

	for i, delayedScenario := range delayedScenarios {
		scenarioIDs[i] = delayedScenario.ScenarioID
		alarmIDs[i] = delayedScenario.AlarmID
	}

	scenariosByID, err := m.loadScenarios(ctx, scenarioIDs)
	if err != nil {
		return nil, err
	}

	alarmsByID, err := m.loadAlarms(ctx, alarmIDs)
	if err != nil {
		return nil, err
	}

	tasks := make([]DelayedScenarioTask, len(delayedScenarios))
	for i, delayedScenario := range delayedScenarios {
		scenario, ok := scenariosByID[delayedScenario.ScenarioID]
		if !ok {
			return nil, fmt.Errorf("cannot find enabled scenario id=%s", delayedScenario.ScenarioID)
		}

		alarm, ok := alarmsByID[delayedScenario.AlarmID]
		if !ok {
			return nil, fmt.Errorf("cannot find opened alarm id=%s", delayedScenario.AlarmID)
		}

		tasks[i] = DelayedScenarioTask{
			Alarm:    alarm,
			Scenario: scenario,

			AdditionalData: delayedScenario.AdditionalData,
		}
	}

	return tasks, nil
}

func (m *delayedScenarioManager) loadScenarios(ctx context.Context, ids []string) (map[string]Scenario, error) {
	scenarios, err := m.scenarioAdapter.GetEnabledByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	scenariosByID := make(map[string]Scenario)
	for _, scenario := range scenarios {
		scenariosByID[scenario.ID] = scenario
	}

	return scenariosByID, nil
}

func (m *delayedScenarioManager) loadAlarms(ctx context.Context, ids []string) (map[string]types.Alarm, error) {
	alarms := make([]types.Alarm, 0)
	err := m.alarmAdapter.GetOpenedAlarmsByAlarmIDs(ctx, ids, &alarms)
	if err != nil {
		return nil, err
	}

	alarmsByID := make(map[string]types.Alarm)
	for _, alarm := range alarms {
		alarmsByID[alarm.ID] = alarm
	}

	return alarmsByID, nil
}

func (m *delayedScenarioManager) canStartWaitGoroutine(executionTime time.Time) bool {
	now := time.Now()
	sub := executionTime.Sub(now)

	return sub > time.Second && sub < m.periodicalTimeout
}
