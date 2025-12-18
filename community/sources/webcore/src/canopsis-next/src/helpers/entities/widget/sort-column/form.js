import { SORT_ORDERS } from '@/constants';

import { addKeyInEntity, removeKeyFromEntity } from '@/helpers/array';

/**
 * @typedef {Object} WidgetSortColumn
 * @property {string} sort_by
 * @property {string} sort
 */

/**
 * @typedef {Object & ObjectKey} WidgetSortColumnForm
 * @property {string} sort_by
 * @property {string} sort
 */

/**
 * Convert widget sort column to form
 *
 * @param {WidgetSortColumn} [sortColumn = {}]
 * @returns {WidgetSortColumnForm}
 */
export const widgetSortColumnToForm = (sortColumn = {}) => addKeyInEntity({
  sort_by: sortColumn.sort_by ?? '',
  sort: (sortColumn.sort ?? SORT_ORDERS.asc).toUpperCase(),
});

/**
 * Convert widget sort columns array to form array
 *
 * @param {WidgetSortColumn[]} [sortColumns = []]
 * @returns {WidgetSortColumnForm[]}
 */
export const widgetSortColumnsToForm = (sortColumns = []) => sortColumns.map(widgetSortColumnToForm);

/**
 * Convert form array to widget sort columns array
 *
 * @param {WidgetSortColumnForm[]} [sortColumns = []]
 * @returns {WidgetSortColumn[]}
 */
export const formToWidgetSortColumns = (sortColumns = []) => sortColumns.map(sortColumn => ({
  ...removeKeyFromEntity(sortColumn),

  sort: sortColumn.sort.toLowerCase(),
}));
