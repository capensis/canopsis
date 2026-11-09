import { createPinia, setActivePinia } from 'pinia';

import { DEFAULT_USER_INTERFACE_DATE_FORMAT, USER_INTERFACE_DATE_FORMATS } from '@/constants';

import { setDateLocale } from '@/helpers/date/date';
import { registerUserInterfaceDateFormatProvider } from '@/helpers/date/date-format';

import { usePbehaviorDateFormat } from '@/components/other/pbehavior/pbehaviors/hooks/pbehavior-date-format';

describe('pbehavior interval date formatting', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    setDateLocale('en');
  });

  afterEach(() => {
    registerUserInterfaceDateFormatProvider(() => DEFAULT_USER_INTERFACE_DATE_FORMAT);
  });

  it.each([
    [USER_INTERFACE_DATE_FORMATS.dateTime.value, '03/11/2026 01:00:00'],
    [USER_INTERFACE_DATE_FORMATS.dateTimeWithWeekday.value, '03/11/2026 (Tue) 01:00:00'],
    [USER_INTERFACE_DATE_FORMATS.weekdayWithDateTime.value, 'Tue 03/11/2026 01:00:00'],
  ])('respects %s for the next occurrence and regular intervals', (format, expected) => {
    registerUserInterfaceDateFormatProvider(() => format);

    const { timezone, formatIntervalDate } = usePbehaviorDateFormat();
    timezone.value = 'Europe/Paris';

    for (const rrule of ['', 'FREQ=DAILY']) {
      const item = { rrule, tstart: new Date('2026-11-03T00:00:00Z') };

      expect(formatIntervalDate(item, 'tstart')).toBe(expected);
    }
  });
});
