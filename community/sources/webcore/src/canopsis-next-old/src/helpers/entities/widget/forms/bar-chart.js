import { DEFAULT_PERIODIC_REFRESH, QUICK_RANGES, SAMPLINGS } from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';
import { metricPresetsToForm, formToMetricPresets } from '@/helpers/entities/metric/form';

/**
 * @typedef {Object} BarChartWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {MetricPreset[]} metrics
 * @property {boolean} stacked
 * @property {string} chart_title
 * @property {string} default_time_range
 * @property {Sampling} default_sampling
 * @property {boolean} comparison
 * @property {string | null} mainFilter
 */

/**
 * @typedef {BarChartWidgetParameters} BarChartWidgetParametersForm
 */

/**
 * Convert bar chart widget parameters to form
 *
 * @param {BarChartWidgetParameters} parameters
 * @return {BarChartWidgetParametersForm}
 */
export const barChartWidgetParametersToForm = (parameters = {}) => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  metrics: metricPresetsToForm(parameters.metrics),
  stacked: parameters.stacked ?? false,
  chart_title: parameters.chart_title ?? '',
  default_time_range: parameters.default_time_range ?? QUICK_RANGES.last30Days.value,
  default_sampling: parameters.default_sampling ?? SAMPLINGS.day,
  comparison: parameters.comparison ?? false,
  mainFilter: parameters.mainFilter ?? null,
});

/**
 * Convert form to bar chart widget parameters
 *
 * @param {BarChartWidgetParametersForm} form
 * @return {BarChartWidgetParameters}
 */
export const formToBarChartWidgetParameters = form => ({
  ...form,

  metrics: formToMetricPresets(form.metrics),
});
