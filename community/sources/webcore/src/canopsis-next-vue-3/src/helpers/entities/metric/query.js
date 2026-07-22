import { isOmitEqual } from '@/helpers/collection'
import { DATETIME_FORMATS, DATETIME_INTERVAL_TYPES } from '@/constants'
import { convertDateIntervalToTimestampByTimezone } from '@/helpers/date/date-intervals'

/**
 * Check query is changed with interval.
 *
 * @param {Object} query
 * @param {Object} oldQuery
 * @param {number} minDate
 * @returns {boolean}
 */
export const isMetricsQueryChanged = (query, oldQuery, minDate) => {
  const isFromChanged = query.interval?.from !== oldQuery.interval?.from
  const isFromEqualMinDate = query.interval?.from === minDate
  const isToChanged = query.interval?.to !== oldQuery.interval?.to
  const isQueryWithoutIntervalChanged = !isOmitEqual(query, oldQuery, ['interval'])

  return isQueryWithoutIntervalChanged || (isFromChanged && !isFromEqualMinDate) || isToChanged
}

/**
 * Convert KPI metrics interval to API timestamps.
 *
 * Metrics API validates both `from` and `to` as day starts in Canopsis timezone.
 *
 * @param {{ from: number|string, to: number|string }} interval
 * @param {string} timezone
 * @returns {{ from: number, to: number }}
 */
export const convertMetricsIntervalToTimestamp = (interval, timezone) => ({
  from: convertDateIntervalToTimestampByTimezone(
    interval.from,
    DATETIME_INTERVAL_TYPES.start,
    DATETIME_FORMATS.dateTimePicker,
    undefined,
    timezone,
  ),
  to: convertDateIntervalToTimestampByTimezone(
    interval.to,
    DATETIME_INTERVAL_TYPES.start,
    DATETIME_FORMATS.dateTimePicker,
    undefined,
    timezone,
  ),
})
