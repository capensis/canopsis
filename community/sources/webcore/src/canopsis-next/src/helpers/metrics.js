import {
  AGGREGATE_FUNCTIONS,
  ALARM_METRIC_PARAMETERS,
  DATETIME_FORMATS,
  KPI_RATING_ENTITY_METRICS,
  KPI_RATING_USER_CRITERIA,
  KPI_RATING_USER_METRICS,
  SAMPLINGS,
  TIME_UNITS,
  USER_METRIC_PARAMETERS,
} from '@/constants';

import { addUnitToDate, convertDateToString, convertDateToTimestampByTimezone } from '@/helpers/date/date';
import { isOmitEqual } from '@/helpers/equal';

/**
 * @typedef { 'hour' | 'day' | 'week' | 'month' } Sampling
 */

/**
 * @typedef {Object} Metric
 * @property {number} timestamp
 */

/**
 * Check metric is time
 *
 * @param {string} metric
 * @returns {boolean}
 */
export const isTimeMetric = metric => [
  USER_METRIC_PARAMETERS.totalUserActivity,
  ALARM_METRIC_PARAMETERS.averageAck,
  ALARM_METRIC_PARAMETERS.averageResolve,
  ALARM_METRIC_PARAMETERS.timeToAck,
  ALARM_METRIC_PARAMETERS.timeToResolve,
].includes(metric);

/**
 * Check metric is ratio
 *
 * @param {string} metric
 * @returns {boolean}
 */
export const isRatioMetric = metric => [
  ALARM_METRIC_PARAMETERS.ratioCorrelation,
  ALARM_METRIC_PARAMETERS.ratioInstructions,
  ALARM_METRIC_PARAMETERS.ratioTickets,
  ALARM_METRIC_PARAMETERS.ratioNonDisplayed,
  ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms,
].includes(metric);

/**
 * Return label by sampling
 *
 * @param {number | string} value
 * @param {Sampling} sampling
 * @returns {string}
 *
 * @example
 * getDateLabelBySampling(1636523087405, 'hour') // 10/11/2021 12:44
 * getDateLabelBySampling(1636523087405, 'day') // 10/11/2021
 * getDateLabelBySampling(1636523087405, 'week') // 10/11/2021 - \n17/11/2021
 * getDateLabelBySampling(1636523087405, 'month') // November 2021
 */
export const getDateLabelBySampling = (value, sampling) => {
  switch (sampling) {
    case SAMPLINGS.hour:
      return convertDateToString(value, DATETIME_FORMATS.dateTimePicker);
    case SAMPLINGS.day:
      return convertDateToString(value, DATETIME_FORMATS.short);
    case SAMPLINGS.week:
      return [
        convertDateToString(value, DATETIME_FORMATS.short),
        convertDateToString(addUnitToDate(value, 1, TIME_UNITS.week), DATETIME_FORMATS.short),
      ].join(' - \n');
  }

  return convertDateToString(value, DATETIME_FORMATS.yearWithMonth);
};

/**
 * Check is user criteria
 *
 * @param {string} criteria
 * @returns {boolean}
 */
const isUserCriteria = criteria => KPI_RATING_USER_CRITERIA.includes(criteria);

/**
 * Get all metrics by criteria
 *
 * @param {string} criteria
 * @returns {string[]}
 */
export const getAvailableMetricsByCriteria = criteria => (
  isUserCriteria(criteria)
    ? KPI_RATING_USER_METRICS
    : KPI_RATING_ENTITY_METRICS
);

/**
 * If metric available for criteria return metric, else return first available metric
 *
 * @param {string} metric
 * @param {string} [criteria]
 * @returns {string}
 */
export const getAvailableMetricByCriteria = (metric, criteria) => {
  const metrics = getAvailableMetricsByCriteria(criteria);

  if (criteria && metrics.includes(metric)) {
    return metric;
  }

  const [firstMetric] = metrics;

  return firstMetric;
};

/**
 * Check query is changed with interval
 *
 * @param {Object} query
 * @param {Object} oldQuery
 * @param {number} minDate
 * @returns {boolean}
 */
export const isMetricsQueryChanged = (query, oldQuery, minDate) => {
  const isFromChanged = query.interval.from !== oldQuery.interval.from;
  const isFromEqualMinDate = query.interval.from === minDate;
  const isToChanged = query.interval.to !== oldQuery.interval.to;
  const isQueryWithoutIntervalChanged = !isOmitEqual(query, oldQuery, ['interval']);

  return isQueryWithoutIntervalChanged || (isFromChanged && !isFromEqualMinDate) || isToChanged;
};

/**
 * Convert metrics timestamps to local timezone
 *
 * @param {Metric[]} metrics
 * @param {string} timezone
 * @returns {Metric[]}
 */
export const convertMetricsToTimezone = (metrics, timezone) => metrics.map(metric => ({
  ...metric,

  timestamp: convertDateToTimestampByTimezone(metric.timestamp, timezone),
}));

/**
 * Get default aggregate function by metric
 *
 * @param {string} metric
 * @returns {string}
 */
export const getDefaultAggregateFunctionByMetric = (metric) => {
  switch (metric) {
    case ALARM_METRIC_PARAMETERS.createdAlarms:
    case ALARM_METRIC_PARAMETERS.ackAlarms:
    case ALARM_METRIC_PARAMETERS.cancelAckAlarms:
    case ALARM_METRIC_PARAMETERS.ratioCorrelation:
    case ALARM_METRIC_PARAMETERS.ratioInstructions:
    case ALARM_METRIC_PARAMETERS.ratioNonDisplayed:
    case ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms:
      return AGGREGATE_FUNCTIONS.sum;
    default:
      return AGGREGATE_FUNCTIONS.avg;
  }
};

/**
 * Get all available aggregate functions by metric
 *
 * @param {string} metric
 * @returns {string[]}
 */
export const getAggregateFunctionsByMetric = (metric) => {
  switch (metric) {
    case ALARM_METRIC_PARAMETERS.activeAlarms:
    case ALARM_METRIC_PARAMETERS.ackActiveAlarms:
    case ALARM_METRIC_PARAMETERS.ticketActiveAlarms:
    case ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms:
    case ALARM_METRIC_PARAMETERS.averageAck:
    case ALARM_METRIC_PARAMETERS.averageResolve:
    case ALARM_METRIC_PARAMETERS.timeToAck:
    case ALARM_METRIC_PARAMETERS.timeToResolve:
    case ALARM_METRIC_PARAMETERS.notAckedAlarms:
    case ALARM_METRIC_PARAMETERS.notAckedInHourAlarms:
    case ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms:
    case ALARM_METRIC_PARAMETERS.notAckedInDayAlarms:
      return [
        AGGREGATE_FUNCTIONS.avg,
        AGGREGATE_FUNCTIONS.min,
        AGGREGATE_FUNCTIONS.max,
      ];
    case ALARM_METRIC_PARAMETERS.ratioCorrelation:
    case ALARM_METRIC_PARAMETERS.ratioInstructions:
    case ALARM_METRIC_PARAMETERS.ratioNonDisplayed:
    case ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms:
      return [AGGREGATE_FUNCTIONS.sum];
    case ALARM_METRIC_PARAMETERS.ratioTickets:
      return [AGGREGATE_FUNCTIONS.avg];
    default:
      return [
        AGGREGATE_FUNCTIONS.sum,
        AGGREGATE_FUNCTIONS.avg,
        AGGREGATE_FUNCTIONS.min,
        AGGREGATE_FUNCTIONS.max,
      ];
  }
};
