package alarm

import (
	"context"
	"fmt"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type service struct {
	cfg     config.CanopsisConf
	adapter Adapter
	logger  zerolog.Logger
}

// NewService gives the correct alarm adapter. Give nil to the redis
// client and it will create a new redis.Client with the dedicated redis
// database for alarms.
func NewService(alarmAdapter Adapter, logger zerolog.Logger, cfg config.CanopsisConf) Service {
	return &service{
		cfg:     cfg,
		adapter: alarmAdapter,
		logger:  logger,
	}
}

func (s *service) ResolveAlarms(ctx context.Context, baggotTime time.Duration) ([]types.Alarm, error) {
	defer trace.StartRegion(ctx, "alarm.ResolveAlarms").End()

	updatedAlarms := make([]types.Alarm, 0)

	unresolvedAlarms, err := s.adapter.GetUnresolved()
	if err != nil {
		return updatedAlarms, fmt.Errorf("unresolved alarms error: %v", err)
	}

	for _, alarm := range unresolvedAlarms {
		if alarm.Closable(baggotTime) {
			updatedAlarms = append(updatedAlarms, alarm)
		}
	}

	return updatedAlarms, nil
}

func (s *service) ResolveCancels(ctx context.Context, cancelAutosolveDelay time.Duration) ([]types.Alarm, error) {
	defer trace.StartRegion(ctx, "alarm.ResolveCancels").End()

	canceledAlarms := make([]types.Alarm, 0)

	alarms, err := s.adapter.GetAlarmsWithCancelMark()
	if err != nil {
		return canceledAlarms, fmt.Errorf("cancel alarms error: %v", err)
	}

	for _, alarm := range alarms {
		if time.Since(alarm.Value.Canceled.Timestamp.Time) >= cancelAutosolveDelay {
			canceledAlarms = append(canceledAlarms, alarm)
		}
	}

	return canceledAlarms, nil
}

func (s *service) ResolveSnoozes(ctx context.Context, disableActionSnoozeDelayOnPbh bool) ([]types.Alarm, error) {
	defer trace.StartRegion(ctx, "alarm.ResolveSnoozes").End()

	unsnoozedAlarms := make([]types.Alarm, 0)

	alarms, err := s.adapter.GetAlarmsWithSnoozeMark()
	if err != nil {
		return unsnoozedAlarms, fmt.Errorf("snooze alarms error: %v", err)
	}

	for _, alarm := range alarms {
		if !alarm.IsSnoozed() && (alarm.IsInActivePeriod() || disableActionSnoozeDelayOnPbh) {
			unsnoozedAlarms = append(unsnoozedAlarms, alarm)
		}
	}

	return unsnoozedAlarms, nil
}

func (s *service) UpdateFlappingAlarms(ctx context.Context) ([]types.Alarm, error) {
	defer trace.StartRegion(ctx, "alarm.UpdateFlappingAlarms").End()

	updatedAlarms := make([]types.Alarm, 0)

	flappingAlarms, err := s.adapter.GetAlarmsWithFlappingStatus()
	if err != nil {
		return updatedAlarms, fmt.Errorf("unable to get alarms with flapping status: %v", err)
	}

	for _, alarm := range flappingAlarms {
		currentAlarmStatus := alarm.CurrentStatus(s.cfg)
		newStatus := alarm.ComputeStatus(s.cfg)

		if newStatus != currentAlarmStatus || alarm.Value.Status == nil {
			updatedAlarms = append(updatedAlarms, alarm)
		}
	}

	return updatedAlarms, nil
}

func (s *service) ResolveDone(ctx context.Context) ([]types.Alarm, error) {
	defer trace.StartRegion(ctx, "alarm.ResolveDone").End()

	doneAlarms := make([]types.Alarm, 0)

	alarms, err := s.adapter.GetAlarmsWithDoneMark()
	if err != nil {
		return doneAlarms, fmt.Errorf("done alarms error: %v", err)
	}

	for _, alarm := range alarms {
		delta := time.Since(alarm.Value.Done.Timestamp.Time)
		if int(delta.Seconds()) >= canopsis.DoneAutosolveDelay {
			doneAlarms = append(doneAlarms, alarm)
		}
	}

	return doneAlarms, nil
}

// UpdateToWorstState updates meta-alarm's state from its worst children, return true when meta-alarm has updated
func UpdateToWorstState(metaAlarm *types.Alarm, updateChildren []*types.Alarm, a Adapter, cfg config.CanopsisConf) bool {
	if !metaAlarm.IsMetaAlarm() {
		return false
	}

	var (
		alarms     []types.Alarm
		stepTs, ts types.CpsTime
		worstState = types.CpsNumber(types.AlarmStateOK)
	)

	if metaAlarm.Value.Children != nil && len(metaAlarm.Value.Children) > 0 {
		if len(metaAlarm.Value.Children) == len(updateChildren) {
			// all children to update
			for _, child := range updateChildren {
				childState := child.CurrentState()
				if childState >= worstState {
					worstState = child.Value.State.Value
					stepTs = child.Value.LastUpdateDate
				}
			}
		} else {
			err := a.GetOpenedAlarmsByIDs(metaAlarm.Value.Children, &alarms)
			if err != nil {
				return false
			}
			for _, child := range alarms {
				childState := child.CurrentState()
				for _, uc := range updateChildren {
					ts = uc.Value.LastUpdateDate
					if uc.ID == child.ID {
						childState = uc.CurrentState()
						break
					}
				}
				if childState > worstState {
					worstState = childState
					stepTs = ts
				}
			}
			if stepTs.IsZero() {
				stepTs = ts
			}
		}
	}

	if worstState == metaAlarm.Value.State.Value {
		// state has not changed
		return false
	}

	stepType := types.AlarmStepStateIncrease
	if worstState < metaAlarm.Value.State.Value {
		stepType = types.AlarmStepStateDecrease
	}

	metaAlarm.Value.State.Value = worstState

	// Create new Step to keep track of the alarm history
	author := metaAlarm.Value.Connector + "." + metaAlarm.Value.ConnectorName
	newStep := types.AlarmStep{
		Type:      stepType,
		Timestamp: stepTs,
		Author:    author,
		Message:   metaAlarm.Value.Output,
		Value:     worstState,
	}

	if err := metaAlarm.Value.Steps.Add(newStep); err != nil {
		return false
	}
	metaAlarm.Value.State = &newStep

	// metaAlarm with state OK to status Off
	if metaAlarm.Value.State.Value == types.AlarmStateOK {
		metaAlarm.UpdateStatus(stepTs, author, metaAlarm.Value.Output, cfg)
	}

	metaAlarm.Value.StateChangesSinceStatusUpdate++
	metaAlarm.Value.TotalStateChanges++
	metaAlarm.Value.LastUpdateDate = stepTs
	return true
}
