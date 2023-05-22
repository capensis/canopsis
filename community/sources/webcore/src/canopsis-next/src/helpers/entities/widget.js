import { pick } from 'lodash';

import {
  BAR_CHART_WIDGET_PRESET_PARAMETERS_BY_TYPE,
  BAR_CHART_WIDGET_PRESET_TYPES,
  CHART_WIDGET_PRESET_TYPES,
  NUMBERS_CHART_WIDGET_PRESET_PARAMETERS_BY_TYPE,
  NUMBERS_CHART_WIDGET_PRESET_TYPES,
  PIE_CHART_WIDGET_PRESET_PARAMETERS_BY_TYPE,
  PIE_CHART_WIDGET_PRESET_TYPES,
  WIDGET_TYPES,
} from '@/constants';

import { getDefaultAggregateFunctionByMetric } from '@/helpers/metrics';

import { metricPresetsToForm } from '../forms/metric';

/**
 * Get chart preset types by widget type
 *
 * @param {WidgetType} type
 * @returns {string[]}
 */
export const getWidgetChartPresetTypesByWidgetType = (type) => {
  const keys = {
    [WIDGET_TYPES.barChart]: BAR_CHART_WIDGET_PRESET_TYPES,
    [WIDGET_TYPES.pieChart]: PIE_CHART_WIDGET_PRESET_TYPES,
    [WIDGET_TYPES.numbers]: NUMBERS_CHART_WIDGET_PRESET_TYPES,
  }[type] ?? [];

  return Object.values(pick(CHART_WIDGET_PRESET_TYPES, keys));
};

/**
 * Get chart preset parameters by preset and widget id
 *
 * @param {WidgetType} type
 * @param {string} preset
 * @returns {Object}
 */
export const getWidgetChartPresetParameters = (type, preset) => {
  const parametersByPreset = {
    [WIDGET_TYPES.barChart]: BAR_CHART_WIDGET_PRESET_PARAMETERS_BY_TYPE,
    [WIDGET_TYPES.pieChart]: PIE_CHART_WIDGET_PRESET_PARAMETERS_BY_TYPE,
    [WIDGET_TYPES.numbers]: NUMBERS_CHART_WIDGET_PRESET_PARAMETERS_BY_TYPE,
  }[type];

  const { metrics, ...parameters } = parametersByPreset[preset] ?? {};

  if (metrics) {
    parameters.metrics = metricPresetsToForm(metrics);

    if (type === WIDGET_TYPES.numbers) {
      parameters.metrics = parameters.metrics.map(metric => ({
        ...metric,
        aggregate_func: metric.aggregate_func || getDefaultAggregateFunctionByMetric(metric.metric),
      }));
    }
  }

  return parameters;
};
