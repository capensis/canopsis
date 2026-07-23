import {
  ALARM_METRIC_PARAMETERS,
  KPI_RATING_USER_CRITERIA,
  USER_METRIC_PARAMETERS,
} from '@/constants';

/**
 * @typedef { 'sum' | 'avg' | 'min' | 'max' } MetricAggregateFunctions
 */

/**
 * @typedef { 'hour' | 'day' | 'week' | 'month' } Sampling
 */

/**
 * @typedef { 'ms' | 'us' | 's' } ExternalMetricTimeUnit
 */

/**
 * @typedef { 'c' | 'B' | 'KB' | 'MB' | 'GB' | 'TB' } ExternalMetricDataSizeUnit
 */

/**
 * @typedef { ExternalMetricTimeUnit | ExternalMetricDataSizeUnit | '%' } ExternalMetricUnit
 */

/**
 * @typedef {Object} MetricPreset
 * @property {string} metric
 * @property {string} [color]
 * @property {string} [label]
 * @property {boolean} [auto]
 * @property {boolean} [external]
 * @property {MetricAggregateFunctions} [aggregate_func]
 */

/**
 * @typedef {MetricPreset[]} MetricPresets
 */

/**
 * @typedef {MetricPreset & ObjectKey} MetricPresetForm
 */

/**
 * @typedef {MetricPresetForm[]} MetricPresetsForm
 */

/**
 * @typedef {Object} Metric
 * @property {number} timestamp
 */

/**
 * Check metric is time-based.
 *
 * @param {string} metric
 * @returns {boolean}
 */
export const isTimeMetric = (metric) =>
  [
    USER_METRIC_PARAMETERS.totalUserActivity,
    USER_METRIC_PARAMETERS.minUserSession,
    USER_METRIC_PARAMETERS.maxUserSession,
    USER_METRIC_PARAMETERS.averageUserSession,
    ALARM_METRIC_PARAMETERS.averageAck,
    ALARM_METRIC_PARAMETERS.averageResolve,
    ALARM_METRIC_PARAMETERS.timeToAck,
    ALARM_METRIC_PARAMETERS.timeToResolve,
    ALARM_METRIC_PARAMETERS.minAck,
    ALARM_METRIC_PARAMETERS.maxAck,
    ALARM_METRIC_PARAMETERS.minResolve,
    ALARM_METRIC_PARAMETERS.maxResolve,
  ].includes(metric);

/**
 * Check metric is percentage/ratio-based.
 *
 * @param {string} metric
 * @returns {boolean}
 */
export const isRatioMetric = (metric) =>
  [
    ALARM_METRIC_PARAMETERS.ratioCorrelation,
    ALARM_METRIC_PARAMETERS.ratioInstructions,
    ALARM_METRIC_PARAMETERS.ratioTickets,
    ALARM_METRIC_PARAMETERS.ratioNonDisplayed,
    ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms,
  ].includes(metric);

/**
 * Check criteria belongs to user rating.
 *
 * @param {string} criteria
 * @returns {boolean}
 */
export const isUserCriteria = (criteria) => KPI_RATING_USER_CRITERIA.includes(criteria);
