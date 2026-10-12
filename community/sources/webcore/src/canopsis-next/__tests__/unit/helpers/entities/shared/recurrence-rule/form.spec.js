import { RRule } from 'rrule';
import { describe, expect, it } from 'vitest';

import {
  formOptionsToRecurrenceRuleOptions,
  recurrenceRuleToFormOptions,
  rruleBodyStringToRecurrenceRuleFormOptions,
} from '@/helpers/entities/shared/recurrence-rule/form';

describe('formOptionsToRecurrenceRuleOptions', () => {
  it('converts a recurrence rule after advanced fields are cleared', () => {
    const options = formOptionsToRecurrenceRuleOptions({
      freq: RRule.DAILY,
      interval: 1,
    });

    expect(options).toEqual({
      freq: RRule.DAILY,
      interval: 1,
    });
  });

  it('normalizes a single BYMONTH value for the multi-select month field', () => {
    const options = rruleBodyStringToRecurrenceRuleFormOptions('FREQ=YEARLY;BYMONTH=3');

    expect(options.bymonth).toEqual([3]);
    expect(options.advancedEnabled).toBe(true);
  });

  it('keeps BYMONTH as an array when serializing the recurrence rule', () => {
    const options = formOptionsToRecurrenceRuleOptions({
      freq: RRule.YEARLY,
      bymonth: [3],
    });

    expect(options.bymonth).toEqual([3]);
  });

  it('excludes advanced fields when advanced mode is disabled', () => {
    const options = formOptionsToRecurrenceRuleOptions({
      freq: RRule.YEARLY,
      bymonth: [3],
      advancedEnabled: false,
    });

    expect(options).toEqual({
      freq: RRule.YEARLY,
    });
  });

  it('defaults advanced mode to disabled for a new recurrence rule', () => {
    const options = recurrenceRuleToFormOptions(new RRule().origOptions);

    expect(options.advancedEnabled).toBe(false);
  });
});
