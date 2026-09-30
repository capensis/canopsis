import { describe, expect, it } from 'vitest';

import { EVENT_FILTER_TYPES } from '@/constants';

import { formToEventFilter } from '@/helpers/entities/event-filter/rule/form';

describe('formToEventFilter', () => {
  it('keeps the recurrence rule in the event-filter request', () => {
    const eventFilter = formToEventFilter({
      type: EVENT_FILTER_TYPES.drop,
      rrule: 'FREQ=DAILY',
      config: {},
      patterns: {},
      exdates: [],
      exceptions: [],
      external_data: {},
    });

    expect(eventFilter.rrule).toBe('FREQ=DAILY');
  });
});
