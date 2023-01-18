import { ALARM_LIST_WIDGET_COLUMNS, CONTEXT_WIDGET_COLUMNS } from '@/constants';

import uid from '@/helpers/uid';

/**
 * @typedef {'impact-state' | 'state'} ColorIndicator
 */

/**
 * @typedef {Object} WidgetColumn
 * @property {string} value
 * @property {boolean} [isHtml]
 * @property {string} [template]
 * @property {ColorIndicator} [colorIndicator]
 */

/**
 * @typedef {Object} WidgetColumnForm
 * @property {string} key
 * @property {string} column
 * @property {string} [dictionary]
 * @property {string} [field]
 * @property {boolean} [isHtml]
 * @property {string} [template]
 * @property {ColorIndicator} [colorIndicator]
 */

/**
 * Check if column is infos
 *
 * @param {string} column
 * @returns {boolean}
 */
export const isInfosWidgetColumn = column => [
  ALARM_LIST_WIDGET_COLUMNS.infos,
  ALARM_LIST_WIDGET_COLUMNS.entityInfos,
  ALARM_LIST_WIDGET_COLUMNS.entityComponentInfos,
  CONTEXT_WIDGET_COLUMNS.infos,
  CONTEXT_WIDGET_COLUMNS.componentInfos,
].some(constantField => column.startsWith(constantField));

/**
 * Check if column is links
 *
 * @param {string} column
 * @returns {boolean}
 */
export const isLinksWidgetColumn = column => column.startsWith(ALARM_LIST_WIDGET_COLUMNS.links);

/**
 * Convert widget columns to form
 *
 * @param {WidgetColumn[]} columns
 * @returns {WidgetColumnForm[]}
 */
export const widgetColumnsToForm = (columns = []) => columns.map(({ value, ...rest }) => {
  const result = {
    ...rest,

    key: uid(),
    column: value,
  };

  if (isInfosWidgetColumn(value)) {
    const [column, dictionary, field] = value.split('.');

    result.column = column;
    result.dictionary = dictionary;
    result.field = field;
  } else if (isLinksWidgetColumn(value)) {
    const [column, field] = value.split('.');

    result.column = column;
    result.field = field;
  }

  return result;
});

/**
 * Convert form to widget columns
 *
 * @param {WidgetColumnForm[]} form
 * @returns {WidgetColumn[]}
 */
export const formToWidgetColumns = (form = []) => form.map(({ key, column, dictionary, field, ...rest }) => {
  let value = column;

  if (isInfosWidgetColumn(column)) {
    value = `${column}.${dictionary}.${field}`;
  } else if (isLinksWidgetColumn(column) && field) {
    value = `${column}.${field}`;
  }

  return {
    ...rest,

    value,
  };
});
