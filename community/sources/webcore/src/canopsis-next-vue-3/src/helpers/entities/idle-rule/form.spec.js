import { describe, expect, it } from 'vitest';

import { ACTION_TYPES, MIN_DURATION_VALUE, MIN_PRIORITY, TIME_UNITS } from '@/constants';

import { formToIdleRule, idleRuleToForm, prepareIdleRuleNumbers } from './form';

describe('prepareIdleRuleNumbers', () => {
  it('restricts negative priority, time range awaiting, and snooze duration values', () => {
    expect(
      prepareIdleRuleNumbers({
        priority: -1,
        duration: {
          value: -5,
          unit: 'm',
        },
        operation: {
          type: ACTION_TYPES.snooze,
          parameters: {
            [ACTION_TYPES.snooze]: {
              duration: {
                value: -10,
                unit: 'h',
              },
            },
          },
        },
      }),
    ).toEqual({
      priority: MIN_PRIORITY,
      duration: {
        value: MIN_DURATION_VALUE,
        unit: 'm',
      },
      operation: {
        type: ACTION_TYPES.snooze,
        parameters: {
          [ACTION_TYPES.snooze]: {
            duration: {
              value: MIN_DURATION_VALUE,
              unit: 'h',
            },
          },
        },
      },
    });
  });

  it('sets valid default PBehavior duration for new idle rules', () => {
    const form = idleRuleToForm();

    expect(form.operation.parameters[ACTION_TYPES.pbehavior].start_on_trigger).toBe(false);
    expect(form.operation.parameters[ACTION_TYPES.pbehavior].duration).toEqual({
      value: MIN_DURATION_VALUE,
      unit: TIME_UNITS.hour,
    });
  });

  it('sends default PBehavior duration when default values are null', () => {
    const idleRule = formToIdleRule({
      ...idleRuleToForm(),
      operation: {
        type: ACTION_TYPES.pbehavior,
        parameters: {
          [ACTION_TYPES.pbehavior]: {
            ...idleRuleToForm().operation.parameters[ACTION_TYPES.pbehavior],
            start_on_trigger: true,
            duration: {
              value: null,
              unit: null,
            },
          },
        },
      },
    });

    expect(idleRule.operation.parameters).toMatchObject({
      start_on_trigger: true,
      duration: {
        value: MIN_DURATION_VALUE,
        unit: TIME_UNITS.hour,
      },
    });
  });
});
