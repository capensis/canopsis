import { convertTimestampToMomentByTimezone } from '@/helpers/date';
import dateFilter from './date';

export default function (date, timezone, format, ignoreTodayChecker) {
  return dateFilter(convertTimestampToMomentByTimezone(date, timezone), format, ignoreTodayChecker);
}

