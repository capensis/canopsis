import { describe, expect, it } from 'vitest';

import { MIN_DURATION_VALUE, MIN_FREQUENCY_LIMIT, MIN_PRIORITY } from '@/constants';

import { prepareAlarmStatusRuleNumbers } from './form';

describe('prepareAlarmStatusRuleNumbers', () => {
  it('restricts negative priority, duration, and frequency limit values', () => {
    expect(
      prepareAlarmStatusRuleNumbers({
        priority: -1,
        duration: {
          value: -5,
          unit: 'm',
        },
        freq_limit: -10,
      }),
    ).toEqual({
      priority: MIN_PRIORITY,
      duration: {
        value: MIN_DURATION_VALUE,
        unit: 'm',
      },
      freq_limit: MIN_FREQUENCY_LIMIT,
    });
  });

  it('preserves empty values so field validation can report them', () => {
    expect(
      prepareAlarmStatusRuleNumbers({
        priority: '',
        duration: {
          value: undefined,
          unit: 'm',
        },
        freq_limit: null,
      }),
    ).toEqual({
      priority: '',
      duration: {
        value: undefined,
        unit: 'm',
      },
      freq_limit: null,
    });
  });
});
