import {
  ALARM_METRIC_PARAMETERS,
  EXTERNAL_METRIC_UNITS,
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
 * @typedef {MetricPreset} MetricPresetForm
 * @property {string} key
 */

/**
 * @typedef {MetricPresetForm[]} MetricPresetsForm
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
  ALARM_METRIC_PARAMETERS.maxAck,
  ALARM_METRIC_PARAMETERS.minAck,
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
 * Check is user criteria
 *
 * @param {string} criteria
 * @returns {boolean}
 */
export const isUserCriteria = criteria => KPI_RATING_USER_CRITERIA.includes(criteria);

/**
 * Check is external percent metric unit
 *
 * @param {ExternalMetricUnit} [unit]
 * @returns {boolean}
 */
export const isExternalPercentMetricUnit = unit => unit === EXTERNAL_METRIC_UNITS.percent;

/**
 * Check is external time metric unit
 *
 * @param {ExternalMetricUnit} [unit]
 * @returns {boolean}
 */
export const isExternalTimeMetricUnit = unit => [
  EXTERNAL_METRIC_UNITS.microsecond,
  EXTERNAL_METRIC_UNITS.millisecond,
  EXTERNAL_METRIC_UNITS.second,
].includes(unit);

/**
 * Check is external data size metric unit
 *
 * @param {ExternalMetricUnit} [unit]
 * @returns {boolean}
 */
export const isExternalDataSizeMetricUnit = unit => [
  EXTERNAL_METRIC_UNITS.continuousCounter,
  EXTERNAL_METRIC_UNITS.byte,
  EXTERNAL_METRIC_UNITS.kilobyte,
  EXTERNAL_METRIC_UNITS.megabyte,
  EXTERNAL_METRIC_UNITS.gigabyte,
  EXTERNAL_METRIC_UNITS.terabyte,
].includes(unit);

/**
 * Convert metric preset to form
 *
 * @param {MetricPreset} preset
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
 * Convert metric presets to form
 *
 * @param {MetricPresets} presets
 * @returns {MetricPresetsForm}
 */
export const metricPresetsToForm = (presets = []) => presets.map(metricPresetToForm);

/**
 * Convert metric preset form object to API compatible object
 *
 * @param {MetricPresetForm} form
 * @returns {MetricPreset}
 */
export const formToMetricPreset = (form) => {
  const { key, color, aggregate_func: aggregateFunc, label, auto, external, ...metricPreset } = form;

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
 * Convert metric preset form object to API compatible object
 *
 * @param {MetricPresetsForm} form
 * @returns {MetricPresets}
 */
export const formToMetricPresets = form => form.map(formToMetricPreset);
