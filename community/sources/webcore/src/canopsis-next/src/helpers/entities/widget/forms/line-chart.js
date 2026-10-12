import { DEFAULT_PERIODIC_REFRESH, QUICK_RANGES, SAMPLINGS } from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';
import { metricPresetsToForm, formToMetricPresets } from '@/helpers/entities/metric/form';

/**
 * @typedef {Object} LineChartWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {MetricPreset[]} metrics
 * @property {string} chart_title
 * @property {string} default_time_range
 * @property {Sampling} default_sampling
 * @property {boolean} comparison
 * @property {string | null} mainFilter
 */

/**
 * @typedef {LineChartWidgetParameters} LineChartWidgetParametersForm
 */

/**
 * Convert line chart widget parameters to form
 *
 * @param {LineChartWidgetParameters} parameters
 * @return {LineChartWidgetParametersForm}
 */
export const lineChartWidgetParametersToForm = (parameters = {}) => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  metrics: metricPresetsToForm(parameters.metrics),
  chart_title: parameters.chart_title ?? '',
  default_time_range: parameters.default_time_range ?? QUICK_RANGES.last30Days.value,
  default_sampling: parameters.default_sampling ?? SAMPLINGS.day,
  comparison: parameters.comparison ?? false,
  mainFilter: parameters.mainFilter ?? null,
});

/**
 * Convert form to line chart widget parameters
 *
 * @param {LineChartWidgetParametersForm} form
 * @return {LineChartWidgetParameters}
 */
export const formToLineChartWidgetParameters = form => ({
  ...form,

  metrics: formToMetricPresets(form.metrics),
});
