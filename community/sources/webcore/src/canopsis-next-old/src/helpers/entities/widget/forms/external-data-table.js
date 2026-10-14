import { PAGINATION_LIMIT } from '@/config';
import {
  DENSE_TYPES,
  CSV_SEPARATORS,
  DEFAULT_PERIODIC_REFRESH,
  EXPORT_CSV_DATETIME_FORMATS,
  SORT_ORDERS,
} from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';
import { columnsParametersToForm } from '@/helpers/entities/widget/forms/alarm';
import { formToWidgetColumns, widgetColumnsToForm } from '@/helpers/entities/widget/column/form';

/**
 * @typedef {Object} ExternalDataTableWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {number} itemsPerPage
 * @property {string} table
 * @property {WidgetSort} sort
 * @property {WidgetColumn[]} widgetColumns
 * @property {WidgetColumn[]} widgetExportColumns
 * @property {WidgetColumnsParameters} [columns]
 * @property {WidgetDenseParameters} dense
 * @property {WidgetCsvSeparator} exportCsvSeparator
 * @property {string} exportCsvDatetimeFormat
 */

/**
 * @typedef {ExternalDataTableWidgetParameters} ExternalDataTableWidgetParametersForm
 * @property {WidgetColumnForm[]} widgetColumns
 * @property {WidgetColumnForm[]} widgetExportColumns
 */

/**
 * Convert widget sort parameters to form
 *
 * @param {WidgetSort} [sort = {}]
 * @return {WidgetSort}
 */
export const widgetSortToForm = (sort = {}) => ({
  order: sort.order ?? SORT_ORDERS.asc,
  column: sort.column ?? '',
});

/**
 * Convert external data table widget parameters to form
 *
 * @param {ExternalDataTableWidgetParameters} parameters
 * @returns {ExternalDataTableWidgetParametersForm}
 */
export const externalDataTableWidgetParametersToForm = parameters => ({
  itemsPerPage: parameters.itemsPerPage ?? PAGINATION_LIMIT,
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  table: parameters.table || '',
  sort: widgetSortToForm(parameters.sort),
  widgetColumns:
    widgetColumnsToForm(parameters.widgetColumns ?? []),
  widgetExportColumns:
    widgetColumnsToForm(parameters.widgetExportColumns ?? []),
  columns: columnsParametersToForm(parameters.columns),
  dense: parameters.dense ?? DENSE_TYPES.large,
  exportCsvSeparator: parameters.exportCsvSeparator ?? CSV_SEPARATORS.comma,
  exportCsvDatetimeFormat: parameters.exportCsvDatetimeFormat ?? EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
});

/**
 * Convert form to external data table widget parameters
 *
 * @param {ExternalDataTableWidgetParametersForm} form
 * @returns {ExternalDataTableWidgetParameters}
 */
export const formToExternalDataTableWidgetParameters = form => ({
  ...form,

  widgetColumns: formToWidgetColumns(form.widgetColumns),
  widgetExportColumns: formToWidgetColumns(form.widgetExportColumns),
});
