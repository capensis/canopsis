import {
  ALARM_METRIC_PARAMETERS,
  KPI_RATING_USER_CRITERIA,
  USER_METRIC_PARAMETERS,
} from '@/constants';
import { uid } from '@/helpers/uid';

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
export const isTimeMetric = metric =>
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
export const isRatioMetric = metric =>
  [
    ALARM_METRIC_PARAMETERS.ratioCorrelation,
    ALARM_METRIC_PARAMETERS.ratioInstructions,
    ALARM_METRIC_PARAMETERS.ratioTickets,
    ALARM_METRIC_PARAMETERS.ratioNonDisplayed,
    ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms,
  ].includes(metric);

/**
 * Converts one metric preset to its editable form representation.
 *
 * @param {MetricPreset} [preset={}]
 * @returns {MetricPresetForm}
 */
export const metricPresetToForm = (preset = {}) => ({
  metric: preset.metric ?? '',
  color: preset.color ?? '',
  aggregate_func: preset.aggregate_func ?? '',
  label: preset.label ?? '',
  auto: preset.auto ?? false,
  external: preset.external ?? false,
  key: uid(),
});

/**
 * Converts metric presets to editable form items.
 *
 * @param {MetricPresets} [presets=[]]
 * @returns {MetricPresetsForm}
 */
export const metricPresetsToForm = (presets = []) => presets.map(metricPresetToForm);

/**
 * Converts editable metric preset form item to widget value.
 *
 * @param {MetricPresetForm} form
 * @returns {MetricPreset}
 */
export const formToMetricPreset = (form) => {
  const {
    key,
    color,
    aggregate_func: aggregateFunc,
    label,
    auto,
    external,
    ...metricPreset
  } = form;

  if (color) {
    metricPreset.color = color;
  }

  if (aggregateFunc) {
    metricPreset.aggregate_func = aggregateFunc;
  }

  if (label) {
    metricPreset.label = label;
  }

  if (auto) {
    metricPreset.auto = auto;
  }

  if (external) {
    metricPreset.external = external;
  }

  return metricPreset;
};

/**
 * Converts editable metric preset form items to widget values.
 *
 * @param {MetricPresetsForm} form
 * @returns {MetricPresets}
 */
export const formToMetricPresets = form => form.map(formToMetricPreset);

/**
 * Check criteria belongs to user rating.
 *
 * @param {string} criteria
 * @returns {boolean}
 */
export const isUserCriteria = criteria => KPI_RATING_USER_CRITERIA.includes(criteria);
