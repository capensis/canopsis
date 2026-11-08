import moment from 'moment-timezone';

import '@/i18n';

import { DEFAULT_USER_INTERFACE_DATE_FORMAT, USER_INTERFACE_DATE_FORMATS } from '@/constants';

import {
  convertDateToDateObject,
  convertDateToEndOfDayDateObject,
  convertDateToStartOfDayDateObject,
  convertDateToString,
  convertDateToTimestampByTimezone,
  setDateLocale,
} from '@/helpers/date/date';
import { registerUserInterfaceDateFormatProvider } from '@/helpers/date/date-format';
import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';
import { pbehaviorExdateItemToForm } from '@/helpers/entities/pbehavior/form';
import { pbehaviorToTimespanRequest } from '@/helpers/entities/pbehavior/timespans/form';

describe.each(Object.values(USER_INTERFACE_DATE_FORMATS))('date boundaries with $format', ({ value }) => {
  beforeEach(() => {
    registerUserInterfaceDateFormatProvider(() => value);
    vi.useFakeTimers();
    vi.setSystemTime(new Date(2026, 9, 29, 12, 30));
  });

  afterEach(() => {
    registerUserInterfaceDateFormatProvider(() => DEFAULT_USER_INTERFACE_DATE_FORMAT);
    setDateLocale('en');
    vi.useRealTimers();
  });

  test.each(['en', 'fr'])('initializes valid dates for today in %s', (locale) => {
    setDateLocale(locale);

    const { begin, end } = pbehaviorExdateItemToForm();

    expect(begin).toEqual(new Date(2026, 9, 29));
    expect(end).toEqual(new Date(2026, 9, 29, 23, 59, 59));
  });

  test('keeps wall time when removing timezone information', () => {
    const date = moment.tz('2026-10-29 12:30', 'Pacific/Kiritimati');

    expect(convertDateToStartOfDayDateObject(date)).toEqual(new Date(2026, 9, 29));
    expect(convertDateToEndOfDayDateObject(date)).toEqual(new Date(2026, 9, 29, 23, 59, 59));
    expect(date.format('YYYY-MM-DD HH:mm')).toBe('2026-10-29 12:30');
  });

  test('sends finite calendar view boundaries in the timespan request', () => {
    const timezone = 'Europe/Paris';
    const from = convertDateToTimestampByTimezone(
      convertDateToStartOfDayDateObject('2026-10-01'),
      timezone,
    );
    const to = convertDateToTimestampByTimezone(
      convertDateToEndOfDayDateObject('2026-10-31'),
      timezone,
    );
    const request = pbehaviorToTimespanRequest({
      pbehavior: { type: { _id: 'maintenance' }, tstart: 1792224000, tstop: 1792310400 },
      from,
      to,
    });

    expect(request.view_from).toBe(moment.tz('2026-10-01', timezone).unix());
    expect(request.view_to).toBe(moment.tz('2026-10-31 23:59:59', timezone).unix());
    expect(JSON.parse(JSON.stringify(request))).toMatchObject({ view_from: from, view_to: to });
  });

  test.each([
    { from: 'now-6d/d', to: 'now/d' },
    { from: '2026-10-23', to: '2026-10-29' },
  ])('sends day-start KPI timestamps for $from to $to', (interval) => {
    expect(convertMetricIntervalToTimestamp({ interval, timezone: 'Europe/Paris' })).toEqual({
      from: moment.tz('2026-10-23', 'Europe/Paris').unix(),
      to: moment.tz('2026-10-29', 'Europe/Paris').unix(),
    });
  });

  test('localizes French weekday names and parses the displayed value', () => {
    setDateLocale('fr');

    const date = new Date(2026, 9, 29, 12, 30);
    const expected = {
      [USER_INTERFACE_DATE_FORMATS.dateTime.value]: '29/10/2026 12:30:00',
      [USER_INTERFACE_DATE_FORMATS.dateTimeWithWeekday.value]: '29/10/2026 (jeu.) 12:30:00',
      [USER_INTERFACE_DATE_FORMATS.weekdayWithDateTime.value]: 'jeu. 29/10/2026 12:30:00',
    };

    expect(convertDateToString(date)).toBe(expected[value]);
    expect(convertDateToDateObject(expected[value], 'long')).toEqual(date);
  });
});
