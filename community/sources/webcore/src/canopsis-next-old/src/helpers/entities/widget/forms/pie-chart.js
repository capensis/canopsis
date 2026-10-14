import {
  AGGREGATE_FUNCTIONS,
  DEFAULT_PERIODIC_REFRESH,
  KPI_PIE_CHART_SHOW_MODS,
  QUICK_RANGES,
  SAMPLINGS,
} from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';
import { metricPresetsToForm, formToMetricPresets } from '@/helpers/entities/metric/form';

/**
 * @typedef {Object} PieChartWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {MetricPreset[]} metrics
 * @property {string} chart_title
 * @property {string} default_time_range
 * @property {Sampling} default_sampling
 * @property {string} aggregate_func
 * @property {string} show_mode
 * @property {string | null} mainFilter
 */

/**
 * @typedef {PieChartWidgetParameters} PieChartWidgetParametersForm
 */

/**
 * Convert pie chart widget parameters to form
 *
 * @param {PieChartWidgetParameters} parameters
 * @return {PieChartWidgetParametersForm}
 */
export const pieChartWidgetParametersToForm = (parameters = {}) => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  metrics: metricPresetsToForm(parameters.metrics),
  chart_title: parameters.chart_title ?? '',
  show_mode: parameters.show_mode ?? KPI_PIE_CHART_SHOW_MODS.numbers,
  aggregate_func: parameters.aggregate_func ?? AGGREGATE_FUNCTIONS.avg,
  default_time_range: parameters.default_time_range ?? QUICK_RANGES.last30Days.value,
  default_sampling: parameters.default_sampling ?? SAMPLINGS.day,
  mainFilter: parameters.mainFilter ?? null,
});

/**
 * Convert form to pie chart widget parameters
 *
 * @param {PieChartWidgetParametersForm} form
 * @return {PieChartWidgetParameters}
 */
export const formToPieChartWidgetParameters = form => ({
  ...form,

  metrics: formToMetricPresets(form.metrics),
});
