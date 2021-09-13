import { convertTimestampToMomentByTimezone, convertDateToString } from '@/helpers/date/date';

export default function (date, timezone, format, defaultValue = '') {
  return convertDateToString(
    convertTimestampToMomentByTimezone(date, timezone),
    format,
    true,
    defaultValue,
  );
}
