import {
  ALARM_INFOS_FIELDS,
  ALARM_LIST_WIDGET_COLUMNS,
  ENTITY_INFOS_FIELDS,
} from '@/constants';

import uid from '@/helpers/uid';

/**
 * @typedef {'impact-state' | 'state'} ColorIndicator
 */

/**
 * @typedef {Object} WidgetColumn
 * @property {string} value
 * @property {boolean} [label]
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
 * @property {string} [label]
 * @property {boolean} [isHtml]
 * @property {string} [template]
 * @property {ColorIndicator} [colorIndicator]
 */

/**
 * Check if column is infos
 *
 * @param {string} [column = '']
 * @returns {string}
 */
export const getInfosWidgetColumn = (column = '') => [
  ...ALARM_INFOS_FIELDS,
  ...ENTITY_INFOS_FIELDS,
].find(constantField => column.startsWith(constantField));

/**
 * Check if column is links
 *
 * @param {string} [column = '']
 * @returns {boolean}
 */
export const isLinksWidgetColumn = (column = '') => column.startsWith(ALARM_LIST_WIDGET_COLUMNS.links);

/**
 * Convert widget column to form
 *
 * @param {string} [value]
 * @param {string} [label]
 * @param {WidgetColumn & { value: undefined, label: undefined }} [rest]
 * @returns {WidgetColumnForm}
 */
export const widgetColumnToForm = ({ value = '', label = '', ...rest } = {}) => {
  const result = {
    ...rest,

    label,
    key: uid(),
    column: value,
  };

  const infosColumn = getInfosWidgetColumn(value) ?? '';

  if (infosColumn) {
    const [dictionary, field] = value.replace(`${infosColumn}.`, '').split('.');

    result.column = infosColumn;
    result.dictionary = dictionary;
    result.field = field;
  } else if (isLinksWidgetColumn(value)) {
    const [column, field] = value.split('.');

    result.column = column;
    result.field = field;
  }

  return result;
};

/**
 * Convert widget columns array to form array
 *
 * @param {WidgetColumn[]} columns
 * @returns {WidgetColumnForm[]}
 */
export const widgetColumnsToForm = (columns = []) => columns.map(widgetColumnToForm);

/**
 * Convert form array to widget columns array
 *
 * @param {WidgetColumnForm[]} form
 * @returns {WidgetColumn[]}
 */
export const formToWidgetColumns = (form = []) => form.map(({ key, column, dictionary, field, ...rest }) => {
  let value = column;

  if (getInfosWidgetColumn(column)) {
    value = `${column}.${dictionary}.${field}`;
  } else if (isLinksWidgetColumn(column) && field) {
    value = `${column}.${field}`;
  }

  return {
    ...rest,

    value,
  };
});
