import { QUICK_RANGES } from './common';

export const AVAILABILITY_SHOW_TYPE = {
  percent: 0,
  duration: 1,
};

export const AVAILABILITY_DISPLAY_PARAMETERS = {
  uptime: 0,
  downtime: 1,
};

export const AVAILABILITY_QUICK_RANGES = {
  [QUICK_RANGES.last1Hour.value]: QUICK_RANGES.last1Hour,
  [QUICK_RANGES.last24Hour.value]: QUICK_RANGES.last24Hour,
  [QUICK_RANGES.today.value]: QUICK_RANGES.today,
  [QUICK_RANGES.yesterday.value]: QUICK_RANGES.yesterday,
  [QUICK_RANGES.last7Days.value]: QUICK_RANGES.last7Days,
  [QUICK_RANGES.previousWeek.value]: QUICK_RANGES.previousWeek,
  [QUICK_RANGES.previousMonth.value]: QUICK_RANGES.previousMonth,
  [QUICK_RANGES.thisMonth.value]: QUICK_RANGES.thisMonth,
  [QUICK_RANGES.thisMonthSoFar.value]: QUICK_RANGES.thisMonthSoFar,
  [QUICK_RANGES.last3Months.value]: QUICK_RANGES.last3Months,
  [QUICK_RANGES.last6Months.value]: QUICK_RANGES.last6Months,
};
