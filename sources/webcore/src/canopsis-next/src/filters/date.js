import moment from 'moment';

import { DATETIME_FORMATS } from '@/constants';

export default function (date, format, ignoreTodayChecker, defaultValue) {
  let momentFormat = format;
  let dateObject;

  if (DATETIME_FORMATS[format]) {
    momentFormat = DATETIME_FORMATS[format];
  }

  if (!date) {
    return defaultValue || date;
  }

  // If it's unix timestamp in seconds
  if (typeof date === 'number' && date < 100000000000) {
    dateObject = moment.unix(date);
  } else {
    dateObject = moment(date);
  }

  if (!date || !dateObject || !dateObject.isValid()) {
    console.warn('Could not build a valid `moment` object from input.');
    return date;
  }

  if (!ignoreTodayChecker && dateObject.isSame(new Date(), 'day')) {
    momentFormat = DATETIME_FORMATS.time;
  }

  return dateObject.format(momentFormat);
}
