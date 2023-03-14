import { omit } from 'lodash';

import { AGGREGATE_FUNCTIONS } from '@/constants';

import uid from '@/helpers/uid';

/**
 * @typedef { 'sum' | 'avg' | 'min' | 'max' } MetricAggregateFunctions
 */

/**
 * @typedef {Object} MetricPreset
 * @property {string} metric
 * @property {string} color
 * @property {MetricAggregateFunctions} aggregate_func
 */

/**
 * @typedef {MetricPreset} MetricPresetForm
 * @property {string} key
 */

/**
 * Convert metric preset to form
 *
 * @param {MetricPreset} preset
 * @returns {MetricPresetForm}
 */
export const metricPresetToForm = (preset = {}) => ({
  metric: preset.metric ?? '',
  color: preset.color ?? '',
  aggregate_func: preset.aggregate_func ?? AGGREGATE_FUNCTIONS.avg,
  key: uid(),
});

/**
 * Convert metric preset form object to API compatible object
 *
 * @param {MetricPresetForm} form
 * @returns {MetricPreset}
 */
export const formToMetricPreset = form => omit(form, ['key']);
