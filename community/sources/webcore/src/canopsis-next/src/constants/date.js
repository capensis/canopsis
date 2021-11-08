export const TIME_UNITS = {
  second: 's',
  minute: 'm',
  hour: 'h',
  day: 'd',
  week: 'w',
  month: 'M',
  year: 'y',
};

export const AVAILABLE_SORTED_TIME_UNITS = [
  TIME_UNITS.year,
  TIME_UNITS.month,
  TIME_UNITS.day,
  TIME_UNITS.hour,
  TIME_UNITS.minute,
  TIME_UNITS.second,
];

export const DEFAULT_DURATION_FORMAT = 'D __ H _ m _ s _';

export const AVAILABLE_TIME_UNITS = {
  second: {
    text: 'common.times.second',
    value: TIME_UNITS.second,
  },
  minute: {
    text: 'common.times.minute',
    value: TIME_UNITS.minute,
  },
  hour: {
    text: 'common.times.hour',
    value: TIME_UNITS.hour,
  },
  day: {
    text: 'common.times.day',
    value: TIME_UNITS.day,
  },
  week: {
    text: 'common.times.week',
    value: TIME_UNITS.week,
  },
  month: {
    text: 'common.times.month',
    value: TIME_UNITS.month,
  },
  year: {
    text: 'common.times.year',
    value: TIME_UNITS.year,
  },
};

export const DURATION_UNITS = {
  minute: AVAILABLE_TIME_UNITS.minute,
  hour: AVAILABLE_TIME_UNITS.hour,
  day: AVAILABLE_TIME_UNITS.day,
  week: AVAILABLE_TIME_UNITS.week,
  month: AVAILABLE_TIME_UNITS.month,
  year: AVAILABLE_TIME_UNITS.month,
};

export const SNOOZE_DURATION_UNITS = {
  second: AVAILABLE_TIME_UNITS.second,
  ...DURATION_UNITS,
};

export const PERIODIC_REFRESH_UNITS = {
  second: AVAILABLE_TIME_UNITS.second,
  minute: AVAILABLE_TIME_UNITS.minute,
  hour: AVAILABLE_TIME_UNITS.hour,
};

export const DEFAULT_PERIODIC_REFRESH = {
  seconds: 60,
  unit: TIME_UNITS.second,
  enabled: false,
};

export const DEFAULT_TIME_INTERVAL = {
  interval: 60,
  unit: TIME_UNITS.second,
};

export const DATETIME_FORMATS = {
  long: 'DD/MM/YYYY H:mm:ss',
  longWithDayOfWeek: 'ddd DD/MM/YYYY H:mm:ss',
  medium: 'DD/MM H:mm',
  short: 'DD/MM/YYYY',
  shortWithDayOfWeek: 'ddd DD/MM/YYYY',
  time: 'H:mm:ss',
  dateTimePicker: 'DD/MM/YYYY HH:mm',
  dateTimePickerWithSeconds: 'DD/MM/YYYY HH:mm:ss',
  datePicker: 'YYYY-MM-DD',
  timePicker: 'HH:mm',
  timePickerWithSeconds: 'HH:mm:ss',
  veeValidateDateTimeFormat: 'dd/MM/yyyy HH:mm',
  refreshFieldFormat: 'Y __ D __ H _ m _ s _',
  testSuiteFormat: 'DD MMM YYYY [\n]H:mm:ss [GMT]',
};

export const DATETIME_INTERVAL_TYPES = {
  start: 'start',
  stop: 'stop',
};

export const SAMPLINGS = {
  hour: 'hour',
  day: 'day',
  week: 'week',
  month: 'month',
};
