import { cloneDeep } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import {
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_COLUMNS,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_CSV_SEPARATORS,
  FILTER_DEFAULT_VALUES,
  SORT_ORDERS,
} from '@/constants';

import { defaultColumnsToColumns } from '@/helpers/entities';

/**
 * @typedef {Object} ContextWidgetParameters
 * @property {WidgetFilter[]} viewFilters
 * @property {string | null} mainFilter
 * @property {WidgetFilterCondition} mainFilterCondition
 * @property {number} itemsPerPage
 * @property {WidgetColumn[]} widgetColumns
 * @property {WidgetColumn[]} serviceDependenciesColumns
 * @property {WidgetColumn[]} widgetExportColumns
 * @property {string[]} selectedTypes
 * @property {WidgetSort} sort
 * @property {WidgetCsvSeparator} exportCsvSeparator
 * @property {string} exportCsvDatetimeFormat
 */

/**
 * Convert context widget parameters to form
 *
 * @param {ContextWidgetParameters | {}} [parameters = {}]
 * @return {ContextWidgetParameters}
 */
export const contextWidgetParametersToForm = (parameters = {}) => ({
  viewFilters: parameters.viewFilters
    ? cloneDeep(parameters.viewFilters)
    : [],
  mainFilter: parameters.mainFilter ?? null,
  mainFilterCondition: parameters.mainFilterCondition ?? FILTER_DEFAULT_VALUES.condition,
  itemsPerPage: parameters.itemsPerPage ?? PAGINATION_LIMIT,
  widgetColumns: parameters.widgetColumns
    ? cloneDeep(parameters.widgetColumns)
    : defaultColumnsToColumns(DEFAULT_CONTEXT_WIDGET_COLUMNS),
  serviceDependenciesColumns: parameters.serviceDependenciesColumns
    ? cloneDeep(parameters.serviceDependenciesColumns)
    : defaultColumnsToColumns(DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
  selectedTypes: parameters.selectedTypes
    ? cloneDeep(parameters.selectedTypes)
    : [],
  sort: parameters.sort ? { ...parameters.sort } : { order: SORT_ORDERS.asc },
  exportCsvSeparator: parameters.exportCsvSeparator ?? EXPORT_CSV_SEPARATORS.comma,
  exportCsvDatetimeFormat: parameters.exportCsvDatetimeFormat ?? EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
  widgetExportColumns: parameters.widgetExportColumns
    ? cloneDeep(parameters.widgetExportColumns)
    : defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_COLUMNS),
});
