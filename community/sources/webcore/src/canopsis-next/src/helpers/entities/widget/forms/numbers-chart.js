import { DEFAULT_PERIODIC_REFRESH, QUICK_RANGES, SAMPLINGS } from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';
import { metricPresetsToForm, formToMetricPresets } from '@/helpers/entities/metric/form';

/**
 * @typedef {Object} NumbersWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {MetricPreset[]} metrics
 * @property {string} chart_title
 * @property {string} default_time_range
 * @property {Sampling} default_sampling
 * @property {string} show_trend
 * @property {number} [font_size]
 * @property {string | null} mainFilter
 */

/**
 * @typedef {NumbersWidgetParameters} NumbersWidgetParametersForm
 */

/**
 * Convert numbers widget parameters to form
 *
 * @param {NumbersWidgetParameters} parameters
 * @return {NumbersWidgetParametersForm}
 */
export const numbersWidgetParametersToForm = (parameters = {}) => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  metrics: metricPresetsToForm(parameters.metrics),
  chart_title: parameters.chart_title ?? '',
  default_time_range: parameters.default_time_range ?? QUICK_RANGES.last30Days.value,
  default_sampling: parameters.default_sampling ?? SAMPLINGS.day,
  show_trend: parameters.show_trend ?? false,
  font_size: parameters.font_size,
  mainFilter: parameters.mainFilter ?? null,
});

/**
 * Convert form to numbers widget parameters
 *
 * @param {NumbersWidgetParametersForm} form
 * @return {NumbersWidgetParameters}
 */
export const formToNumbersWidgetParameters = form => ({
  ...form,

  metrics: formToMetricPresets(form.metrics),
});
