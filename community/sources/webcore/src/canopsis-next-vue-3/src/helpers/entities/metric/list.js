import {
  DATETIME_FORMATS,
  KPI_RATING_ENTITY_METRICS,
  KPI_RATING_USER_METRICS,
  SAMPLINGS,
  TIME_UNITS,
} from '@/constants'

import {
  addUnitToDate,
  convertDateToString,
  convertDateToTimestampByTimezone,
} from '@/helpers/date/date'

import { isUserCriteria } from './form'

/**
 * Convert metrics timestamps to local timezone.
 *
 * @param {Array} metrics
 * @param {string} timezone
 * @returns {Array}
 */
export const convertMetricsToTimezone = (metrics, timezone) =>
  metrics.map((metric) => ({
    ...metric,
    timestamp: convertDateToTimestampByTimezone(metric.timestamp, timezone),
  }))

/**
 * Return label by sampling.
 *
 * @param {number|string} value
 * @param {string} sampling
 * @returns {string}
 */
export const getDateLabelBySampling = (value, sampling) => {
  switch (sampling) {
    case SAMPLINGS.hour:
      return convertDateToString(value, DATETIME_FORMATS.dateTimePicker)
    case SAMPLINGS.day:
      return convertDateToString(value, DATETIME_FORMATS.short)
    case SAMPLINGS.week:
      return [
        convertDateToString(value, DATETIME_FORMATS.short),
        convertDateToString(addUnitToDate(value, 1, TIME_UNITS.week), DATETIME_FORMATS.short),
      ].join(' - \n')
  }

  return convertDateToString(value, DATETIME_FORMATS.yearWithMonth)
}

/**
 * Get all metrics available for a rating criteria.
 *
 * @param {string} criteria
 * @returns {string[]}
 */
export const getAvailableMetricsByCriteria = (criteria) =>
  isUserCriteria(criteria) ? KPI_RATING_USER_METRICS : KPI_RATING_ENTITY_METRICS

/**
 * Return selected metric if it is available for criteria, otherwise first available metric.
 *
 * @param {string} metric
 * @param {string} [criteria]
 * @returns {string}
 */
export const getAvailableMetricByCriteria = (metric, criteria) => {
  const metrics = getAvailableMetricsByCriteria(criteria)

  if (criteria && metrics.includes(metric)) {
    return metric
  }

  const [firstMetric] = metrics

  return firstMetric
}
