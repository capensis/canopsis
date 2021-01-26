import { convertTimestampToMomentByTimezone } from '@/helpers/date/date';

import dateFilter from './date';

export default function (date, timezone, format, ignoreTodayChecker, defaultValue = '') {
  if (!date) {
    return defaultValue;
  }

  return dateFilter(convertTimestampToMomentByTimezone(date, timezone), format, ignoreTodayChecker, defaultValue);
}

