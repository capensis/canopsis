import { TIME_UNITS, DAYS_IN_YEAR, MONTHS_IN_YEAR, DAYS_IN_MONTH, DAYS_IN_WEEK } from '@/constants';

import { convertUnit } from '@/helpers/date/duration';

const SECONDS_IN_MINUTE = 60;
const MINUTES_IN_HOUR = 60;
const HOURS_IN_DAY = 24;

const SECONDS_IN_HOUR = SECONDS_IN_MINUTE * MINUTES_IN_HOUR;
const SECONDS_IN_DAY = SECONDS_IN_HOUR * HOURS_IN_DAY;
const SECONDS_IN_WEEK = SECONDS_IN_DAY * DAYS_IN_WEEK;
const SECONDS_IN_MONTH = SECONDS_IN_DAY * DAYS_IN_MONTH;
const SECONDS_IN_YEAR = SECONDS_IN_DAY * DAYS_IN_YEAR;

const MINUTES_IN_DAY = MINUTES_IN_HOUR * HOURS_IN_DAY;
const MINUTES_IN_WEEK = MINUTES_IN_DAY * DAYS_IN_WEEK;
const MINUTES_IN_MONTH = MINUTES_IN_DAY * DAYS_IN_MONTH;
const MINUTES_IN_YEAR = MINUTES_IN_DAY * DAYS_IN_YEAR;

const HOURS_IN_WEEK = HOURS_IN_DAY * DAYS_IN_WEEK;
const HOURS_IN_MONTH = HOURS_IN_DAY * DAYS_IN_MONTH;
const HOURS_IN_YEAR = HOURS_IN_DAY * DAYS_IN_YEAR;

const rangeForEach = (from, to, callback) => {
  // eslint-disable-next-line no-plusplus
  for (let number = from; number <= to; number++) {
    callback(number);
  }
};

describe('convertUnit', () => {
  it('convertUnit from seconds', async () => {
    rangeForEach(0, 61, (seconds) => {
      expect(convertUnit(seconds, TIME_UNITS.second, TIME_UNITS.second)).toEqual(seconds);
      expect(convertUnit(seconds * SECONDS_IN_MINUTE, TIME_UNITS.second, TIME_UNITS.minute)).toEqual(seconds);
      expect(convertUnit(seconds * SECONDS_IN_HOUR, TIME_UNITS.second, TIME_UNITS.hour)).toEqual(seconds);
      expect(convertUnit(seconds * SECONDS_IN_DAY, TIME_UNITS.second, TIME_UNITS.day)).toEqual(seconds);
      expect(convertUnit(seconds * SECONDS_IN_WEEK, TIME_UNITS.second, TIME_UNITS.week)).toEqual(seconds);
      expect(convertUnit(seconds * SECONDS_IN_MONTH, TIME_UNITS.second, TIME_UNITS.month)).toEqual(seconds);
      expect(convertUnit(seconds * SECONDS_IN_YEAR, TIME_UNITS.second, TIME_UNITS.year)).toEqual(seconds);
    });
  });

  it('convertUnit from minute', async () => {
    rangeForEach(0, 61, (minutes) => {
      expect(convertUnit(minutes, TIME_UNITS.minute, TIME_UNITS.second)).toEqual(minutes * SECONDS_IN_MINUTE);
      expect(convertUnit(minutes, TIME_UNITS.minute, TIME_UNITS.minute)).toEqual(minutes);
      expect(convertUnit(minutes * MINUTES_IN_HOUR, TIME_UNITS.minute, TIME_UNITS.hour)).toEqual(minutes);
      expect(convertUnit(minutes * MINUTES_IN_DAY, TIME_UNITS.minute, TIME_UNITS.day)).toEqual(minutes);
      expect(convertUnit(minutes * MINUTES_IN_WEEK, TIME_UNITS.minute, TIME_UNITS.week)).toEqual(minutes);
      expect(convertUnit(minutes * MINUTES_IN_MONTH, TIME_UNITS.minute, TIME_UNITS.month)).toEqual(minutes);
      expect(convertUnit(minutes * MINUTES_IN_YEAR, TIME_UNITS.minute, TIME_UNITS.year)).toEqual(minutes);
    });
  });

  it('convertUnit from hour', async () => {
    rangeForEach(0, 25, (hours) => {
      expect(convertUnit(hours, TIME_UNITS.hour, TIME_UNITS.second)).toEqual(hours * SECONDS_IN_HOUR);
      expect(convertUnit(hours, TIME_UNITS.hour, TIME_UNITS.minute)).toEqual(hours * MINUTES_IN_HOUR);
      expect(convertUnit(hours, TIME_UNITS.hour, TIME_UNITS.hour)).toEqual(hours);
      expect(convertUnit(hours * HOURS_IN_DAY, TIME_UNITS.hour, TIME_UNITS.day)).toEqual(hours);
      expect(convertUnit(hours * HOURS_IN_WEEK, TIME_UNITS.hour, TIME_UNITS.week)).toEqual(hours);
      expect(convertUnit(hours * HOURS_IN_MONTH, TIME_UNITS.hour, TIME_UNITS.month)).toEqual(hours);
      expect(convertUnit(hours * HOURS_IN_YEAR, TIME_UNITS.hour, TIME_UNITS.year)).toEqual(hours);
    });
  });

  it('convertUnit from day', async () => {
    rangeForEach(0, 32, (days) => {
      expect(convertUnit(days, TIME_UNITS.day, TIME_UNITS.second)).toEqual(days * SECONDS_IN_DAY);
      expect(convertUnit(days, TIME_UNITS.day, TIME_UNITS.minute)).toEqual(days * MINUTES_IN_DAY);
      expect(convertUnit(days, TIME_UNITS.day, TIME_UNITS.hour)).toEqual(days * HOURS_IN_DAY);
      expect(convertUnit(days, TIME_UNITS.day, TIME_UNITS.day)).toEqual(days);
      expect(convertUnit(days * DAYS_IN_WEEK, TIME_UNITS.day, TIME_UNITS.week)).toEqual(days);
      expect(convertUnit(days * DAYS_IN_MONTH, TIME_UNITS.day, TIME_UNITS.month)).toEqual(days);
      expect(convertUnit(days * DAYS_IN_YEAR, TIME_UNITS.day, TIME_UNITS.year)).toEqual(days);
    });
  });

  it('convertUnit from week', async () => {
    rangeForEach(0, 53, (weeks) => {
      expect(convertUnit(weeks, TIME_UNITS.week, TIME_UNITS.second)).toEqual(weeks * SECONDS_IN_WEEK);
      expect(convertUnit(weeks, TIME_UNITS.week, TIME_UNITS.minute)).toEqual(weeks * MINUTES_IN_WEEK);
      expect(convertUnit(weeks, TIME_UNITS.week, TIME_UNITS.hour)).toEqual(weeks * HOURS_IN_WEEK);
      expect(convertUnit(weeks, TIME_UNITS.week, TIME_UNITS.day)).toEqual(weeks * DAYS_IN_WEEK);
      expect(convertUnit(weeks, TIME_UNITS.week, TIME_UNITS.week)).toEqual(weeks);
    });
  });

  it('convertUnit from month', async () => {
    rangeForEach(0, 13, (months) => {
      expect(convertUnit(months, TIME_UNITS.month, TIME_UNITS.second)).toEqual(months * SECONDS_IN_MONTH);
      expect(convertUnit(months, TIME_UNITS.month, TIME_UNITS.minute)).toEqual(months * MINUTES_IN_MONTH);
      expect(convertUnit(months, TIME_UNITS.month, TIME_UNITS.hour)).toEqual(months * HOURS_IN_MONTH);
      expect(convertUnit(months, TIME_UNITS.month, TIME_UNITS.day)).toEqual(months * DAYS_IN_MONTH);
      expect(convertUnit(months, TIME_UNITS.month, TIME_UNITS.month)).toEqual(months);
      expect(convertUnit(months * MONTHS_IN_YEAR, TIME_UNITS.month, TIME_UNITS.year)).toEqual(months);
    });
  });

  it('convertUnit from year', async () => {
    rangeForEach(1, 13, (years) => {
      expect(convertUnit(years, TIME_UNITS.year, TIME_UNITS.second)).toEqual(years * SECONDS_IN_YEAR);
      expect(convertUnit(years, TIME_UNITS.year, TIME_UNITS.minute)).toEqual(years * MINUTES_IN_YEAR);
      expect(convertUnit(years, TIME_UNITS.year, TIME_UNITS.hour)).toEqual(years * HOURS_IN_YEAR);
      expect(convertUnit(years, TIME_UNITS.year, TIME_UNITS.day)).toEqual(years * DAYS_IN_YEAR);
      expect(convertUnit(years, TIME_UNITS.year, TIME_UNITS.month)).toEqual(years * MONTHS_IN_YEAR);
      expect(convertUnit(years, TIME_UNITS.year, TIME_UNITS.year)).toEqual(years);
    });
  });
});
