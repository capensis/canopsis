import { cloneDeep } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import {
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_COLUMNS,
  DEFAULT_CONTEXT_ALARMS_COLUMNS,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_CSV_SEPARATORS,
  SORT_ORDERS,
} from '@/constants';

import { defaultColumnsToColumns } from '@/helpers/entities';

/**
 * @typedef {Object} ContextWidgetParameters
 * @property {string | null} mainFilter
 * @property {number} itemsPerPage
 * @property {WidgetColumn[]} widgetColumns
 * @property {WidgetColumn[]} serviceDependenciesColumns
 * @property {WidgetColumn[]} widgetExportColumns
 * @property {WidgetColumn[]} resolvedAlarmsColumns
 * @property {WidgetColumn[]} activeAlarmsColumns
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
  mainFilter: parameters.mainFilter ?? null,
  itemsPerPage: parameters.itemsPerPage ?? PAGINATION_LIMIT,
  widgetColumns: parameters.widgetColumns
    ? cloneDeep(parameters.widgetColumns)
    : defaultColumnsToColumns(DEFAULT_CONTEXT_WIDGET_COLUMNS),
  serviceDependenciesColumns: parameters.serviceDependenciesColumns
    ? cloneDeep(parameters.serviceDependenciesColumns)
    : defaultColumnsToColumns(DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
  resolvedAlarmsColumns: parameters.resolvedAlarmsColumns
    ? cloneDeep(parameters.resolvedAlarmsColumns)
    : defaultColumnsToColumns(DEFAULT_CONTEXT_ALARMS_COLUMNS),
  activeAlarmsColumns: parameters.resolvedAlarmsColumns
    ? cloneDeep(parameters.activeAlarmsColumns)
    : defaultColumnsToColumns(DEFAULT_CONTEXT_ALARMS_COLUMNS),
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
