import { DEFAULT_PERIODIC_REFRESH, QUICK_RANGES } from '@/constants';

import { uid } from '@/helpers/uid';
import { durationWithEnabledToForm } from '@/helpers/date/duration';
import { addKeyInEntities, removeKeyFromEntities } from '@/helpers/entities';

/**
 * @typedef {Object} StatisticsWidgetColumn
 * @property {string} column
 * @property {boolean} split
 * @property {string} infos
 */

/**
 * @typedef {StatisticsWidgetColumn & ObjectKey} StatisticsWidgetColumnForm
 */

/**
 * @typedef {Object} StatisticsWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {string} mainParameter
 * @property {Filter[]} patterns
 * @property {StatisticsWidgetColumn[]} widgetColumns
 * @property {string} table_title
 * @property {string} default_time_range
 * @property {string | null} mainFilter
 */

/**
 * @typedef {StatisticsWidgetParameters} StatisticsWidgetParametersForm
 * @property {(Filter & ObjectKey)[]} patterns
 * @property {StatisticsWidgetColumn[]} StatisticsWidgetColumnForm
 */

/**
 * Convert statistics widget column to form
 *
 * @param {StatisticsWidgetColumn} [widgetColumn = {}]
 * @returns {StatisticsWidgetColumnForm}
 */
export const statisticsWidgetColumnToForm = (widgetColumn = {}) => ({
  column: widgetColumn.column ?? '',
  infos: widgetColumn.infos ?? '',
  split: !!widgetColumn.split,
  key: uid(),
});

/**
 * Convert form to statistics widget parameters to form
 *
 * @param {StatisticsWidgetParameters} [parameters = {}]
 * @returns {StatisticsWidgetParametersForm}
 */
export const statisticsWidgetParametersToForm = (parameters = {}) => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  mainParameter: parameters.mainParameter ?? '',
  patterns: addKeyInEntities(parameters.patterns ?? []),
  widgetColumns: (parameters.widgetColumns ?? []).map(statisticsWidgetColumnToForm),
  table_title: parameters.table_title ?? '',
  default_time_range: parameters.default_time_range ?? QUICK_RANGES.last30Days.value,
  mainFilter: parameters.mainFilter ?? null,
});

/**
 * Convert form to statistics widget parameters to form
 *
 * @param {StatisticsWidgetParametersForm} form
 * @returns {StatisticsWidgetParameters}
 */
export const formToStatisticsWidgetParameters = form => ({
  ...form,

  patterns: removeKeyFromEntities(form.patterns),
  widgetColumns: removeKeyFromEntities(form.patterns),
});
