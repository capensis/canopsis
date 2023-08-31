import { DEFAULT_PERIODIC_REFRESH, KPI_ENTITY_RATING_SETTINGS_CUSTOM_CRITERIA, QUICK_RANGES } from '@/constants';

import { uid } from '@/helpers/uid';
import { durationWithEnabledToForm } from '@/helpers/date/duration';
import { addKeyInEntities, removeKeyFromEntities } from '@/helpers/array';

/**
 * @typedef {Object} StatisticsWidgetColumn
 * @property {string} metric
 * @property {string} criteria
 * @property {string} label
 */

/**
 * @typedef {StatisticsWidgetColumn & ObjectKey} StatisticsWidgetColumnForm
 * @property {boolean} split
 */

/**
 * @typedef {Object} StatisticsWidgetMainParameter
 * @property {string | number} criteria
 * @property {string} columnName
 * @property {Filter[]} patterns
 */

/**
 * @typedef {StatisticsWidgetMainParameter} StatisticsWidgetMainParameterForm
 * @property {string | number} criteria
 * @property {string} columnName
 * @property {(Filter & ObjectKey)[]} patterns
 */

/**
 * @typedef {Object} StatisticsWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {StatisticsWidgetMainParameter} mainParameter
 * @property {StatisticsWidgetColumn[]} widgetColumns
 * @property {string} table_title
 * @property {string} default_time_range
 * @property {string | null} mainFilter
 */

/**
 * @typedef {StatisticsWidgetParameters} StatisticsWidgetParametersForm
 * @property {StatisticsWidgetMainParameterForm} mainParameter
 * @property {StatisticsWidgetColumn[]} StatisticsWidgetColumnForm
 */

/**
 * Convert statistics widget column to form
 *
 * @param {StatisticsWidgetColumn} [widgetColumn = {}]
 * @returns {StatisticsWidgetColumnForm}
 */
export const statisticsWidgetColumnToForm = (widgetColumn = {}) => ({
  metric: widgetColumn.metric ?? '',
  criteria: widgetColumn.criteria ?? '',
  label: widgetColumn.label ?? '',
  split: !!widgetColumn.criteria,
  key: uid(),
});

/**
 * Convert statistics widget main parameter to form
 *
 * @param {StatisticsWidgetMainParameter} mainParameter
 * @returns {StatisticsWidgetMainParameterForm}
 */
export const statisticsMainParameterToForm = (mainParameter = {}) => ({
  criteria: mainParameter.criteria ?? '',
  columnName: mainParameter.columnName ?? '',
  patterns: addKeyInEntities(mainParameter.patterns),
});

/**
 * Convert form to statistics widget parameters to form
 *
 * @param {StatisticsWidgetParameters} [parameters = {}]
 * @returns {StatisticsWidgetParametersForm}
 */
export const statisticsWidgetParametersToForm = (parameters = {}) => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  mainParameter: statisticsMainParameterToForm(parameters.mainParameter),
  widgetColumns: (parameters.widgetColumns ?? []).map(statisticsWidgetColumnToForm),
  table_title: parameters.table_title ?? '',
  default_time_range: parameters.default_time_range ?? QUICK_RANGES.last30Days.value,
  mainFilter: parameters.mainFilter ?? null,
});

/**
 * Convert form to statistics widget column
 *
 * @param {StatisticsWidgetColumnForm} form
 * @returns {StatisticsWidgetColumn}
 */
export const formToStatisticsWidgetColumn = form => ({
  metric: form.metric,
  label: form.label,
  criteria: form.split && form.criteria ? form.criteria : '',
});

/**
 * Convert form to statistics widget main parameter
 *
 * @param {StatisticsWidgetMainParameterForm} form
 * @returns {StatisticsWidgetMainParameter}
 */
export const formToStatisticsMainParameter = (form) => {
  const mainParameter = {
    ...form,

    patterns: removeKeyFromEntities(form.patterns),
  };

  if (form.criteria !== KPI_ENTITY_RATING_SETTINGS_CUSTOM_CRITERIA) {
    mainParameter.columnName = '';
    mainParameter.patterns = [];
  }

  return mainParameter;
};

/**
 * Convert form to statistics widget parameters to form
 *
 * @param {StatisticsWidgetParametersForm} form
 * @returns {StatisticsWidgetParameters}
 */
export const formToStatisticsWidgetParameters = form => ({
  ...form,

  mainParameter: formToStatisticsMainParameter(form.mainParameter),
  widgetColumns: form.widgetColumns.map(formToStatisticsWidgetColumn),
});
