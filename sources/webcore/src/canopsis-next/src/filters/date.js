import convertTimestampToMoment from '@/helpers/date';

import { DATETIME_FORMATS } from '@/constants';

export default function (date, format, ignoreTodayChecker, defaultValue) {
  let momentFormat = DATETIME_FORMATS[format] || format;

  if (!date) {
    return defaultValue || date;
  }

  const dateObject = convertTimestampToMoment(date);

  if (!dateObject || !dateObject.isValid()) {
    console.warn('Could not build a valid `moment` object from input.');
    return date;
  }

  if (!ignoreTodayChecker && dateObject.isSame(new Date(), 'day')) {
    momentFormat = DATETIME_FORMATS.time;
  }

  return dateObject.format(momentFormat);
}

