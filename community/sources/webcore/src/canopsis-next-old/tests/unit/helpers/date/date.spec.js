import { convertDateToString } from '@/helpers/date/date';

describe('convertDateToString', () => {
  const year = '2021';
  const month = '09';
  const monthString = 'September';
  const date = '20';
  const weekDay = 'Monday';
  const hours = '21';
  const minutes = '15';
  const seconds = '40';
  const ms = '5';

  const dateInstance = new Date(+year, +month - 1, +date, +hours, +minutes, +seconds, +ms);

  const timestamp = dateInstance.getTime();
  const expectedString = '20/09/2021 21:15:40';

  it('Date as timestamp with timestamp format', () => {
    expect(convertDateToString(timestamp, 'x')).toEqual(String(timestamp));
  });

  it('Convert date to formats', () => {
    expect(convertDateToString(timestamp, 'D')).toEqual(date);
    expect(convertDateToString(timestamp, 'DD')).toEqual(date);
    expect(convertDateToString(timestamp, 'ddd')).toEqual(weekDay.slice(0, 3));
    expect(convertDateToString(timestamp, 'dddd')).toEqual(weekDay);

    expect(convertDateToString(timestamp, 'M')).toEqual(month.slice(1, 2));
    expect(convertDateToString(timestamp, 'MM')).toEqual(month);
    expect(convertDateToString(timestamp, 'MMM')).toEqual(monthString.slice(0, 3));
    expect(convertDateToString(timestamp, 'MMMM')).toEqual(monthString);

    expect(convertDateToString(timestamp, 'Y')).toEqual(year);
    expect(convertDateToString(timestamp, 'YY')).toEqual(year.slice(2, 4));
    expect(convertDateToString(timestamp, 'YYYY')).toEqual(year);

    expect(convertDateToString(timestamp, 'H')).toEqual(hours);
    expect(convertDateToString(timestamp, 'HH')).toEqual(hours);
    expect(convertDateToString(timestamp, 'h')).toEqual(String(hours % 12));
    expect(convertDateToString(timestamp, 'hh')).toEqual(`0${String(hours % 12)}`);

    expect(convertDateToString(timestamp, 'm')).toEqual(minutes);
    expect(convertDateToString(timestamp, 'mm')).toEqual(minutes);

    expect(convertDateToString(timestamp, 's')).toEqual(seconds);
    expect(convertDateToString(timestamp, 'ss')).toEqual(seconds);
  });

  it('Date as unix', () => {
    expect(convertDateToString(timestamp / 1000, 'DD/MM/YYYY H:mm:ss')).toEqual(expectedString);
  });

  it('Date as native Date object', () => {
    expect(convertDateToString(new Date(timestamp), 'DD/MM/YYYY H:mm:ss')).toEqual(expectedString);
  });

  it('Date as ISO string', () => {
    expect(convertDateToString(new Date(timestamp).toISOString(), 'DD/MM/YYYY H:mm:ss')).toEqual(expectedString);
  });

  it('Invalid date without default value', () => {
    expect(convertDateToString(null, 'DD/MM/YYYY H:mm:ss')).toEqual('');
  });

  it('Invalid date with default value', () => {
    const defaultValue = new Date().toISOString();

    expect(convertDateToString(NaN, 'DD/MM/YYYY H:mm:ss', defaultValue)).toEqual(defaultValue);
  });
});
