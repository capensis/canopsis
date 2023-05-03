import { uid } from '@/helpers/uid';

/**
 * @typedef { 'sum' | 'avg' | 'min' | 'max' } MetricAggregateFunctions
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
