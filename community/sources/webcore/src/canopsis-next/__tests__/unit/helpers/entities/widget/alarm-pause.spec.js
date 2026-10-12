import { ALARM_LIST_PAUSE_METHODS } from '@/constants';

import {
  alarmListPauseParametersToForm,
  formToAlarmListWidgetParameters,
  getAlarmListPauseMode,
} from '@/helpers/entities/widget/forms/alarm';

describe('alarm list pause settings helpers', () => {
  test('uses the pause action for new widgets', () => {
    expect(getAlarmListPauseMode()).toBe(ALARM_LIST_PAUSE_METHODS.pause);
    expect(alarmListPauseParametersToForm()).toEqual({
      mode: ALARM_LIST_PAUSE_METHODS.pause,
      types: [],
    });
  });

  test('keeps the PBehavior mode for legacy widgets with fast PBehaviors', () => {
    expect(getAlarmListPauseMode({ fast_pbehaviors: [{ type: 'pause-type' }] })).toBe(
      ALARM_LIST_PAUSE_METHODS.pbehavior,
    );
  });

  test('an explicitly selected pause mode takes precedence over legacy settings', () => {
    const parameters = {
      pause: { mode: ALARM_LIST_PAUSE_METHODS.pause },
      fast_pbehaviors: [{ type: 'pause-type' }],
    };

    expect(getAlarmListPauseMode(parameters)).toBe(ALARM_LIST_PAUSE_METHODS.pause);
  });

  test('removes UI keys from saved pause types', () => {
    const parameters = formToAlarmListWidgetParameters({
      pause: {
        mode: ALARM_LIST_PAUSE_METHODS.pause,
        types: [{ key: 'ui-key', name: 'Maintenance', comment: 'Comment', type: 'maintenance' }],
      },
    });

    expect(parameters.pause).toEqual({
      mode: ALARM_LIST_PAUSE_METHODS.pause,
      types: [{ name: 'Maintenance', comment: 'Comment', type: 'maintenance' }],
    });
  });
});
