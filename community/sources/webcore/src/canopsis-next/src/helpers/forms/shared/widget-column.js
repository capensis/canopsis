import {
  ALARM_FIELDS,
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
 * @property {string} [onlyIcon]
 * @property {ColorIndicator} [colorIndicator]
 */

/**
 * @typedef {Object} WidgetColumnForm
 * @property {string} key
 * @property {string} column
 * @property {string} [rule]
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
 * @param {string} [value = '']
 * @returns {string}
 */
export const getInfosWidgetColumn = (value = '') => [
  ...ALARM_INFOS_FIELDS,
  ...ENTITY_INFOS_FIELDS,
].find(constantField => value.startsWith(constantField));

/**
 * Check if column is links
 *
 * @param {string} [value = '']
 * @returns {boolean}
 */
export const isLinksWidgetColumn = (value = '') => value.startsWith(ALARM_LIST_WIDGET_COLUMNS.links);

/**
 * Convert widget column to form
 *
 * @param {string} [value = '']
 * @param {string} [label = '']
 * @param {boolean} [onlyIcon = false]
 * @param {WidgetColumn & { value: undefined, label: undefined }} [rest]
 * @returns {WidgetColumnForm}
 */
export const widgetColumnToForm = ({ value = '', label = '', onlyIcon = false, ...rest } = {}) => {
  const result = {
    ...rest,

    label,
    key: uid(),
    column: value,
  };

  const infosColumn = getInfosWidgetColumn(value) ?? '';

  if (infosColumn === ALARM_FIELDS.infos) {
    const [rule, dictionary] = value.replace(`${infosColumn}.`, '').split('.');

    result.column = infosColumn;
    result.rule = rule;
    result.dictionary = dictionary;
  } else if (infosColumn) {
    const [dictionary, field] = value.replace(`${infosColumn}.`, '').split('.');

    result.column = infosColumn;
    result.dictionary = dictionary;
    result.field = field;
  } else if (isLinksWidgetColumn(value)) {
    const [column, field = ''] = value.split('.');

    result.column = column;
    result.field = field;
    result.onlyIcon = onlyIcon;
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
export const formToWidgetColumns = (form = []) => (
  form.map(({ key, column, dictionary, field, onlyIcon, rule, ...rest }) => {
    const result = {
      ...rest,

      value: column,
    };

    if (column === ALARM_FIELDS.infos) {
      result.value = `${column}.${rule}.${dictionary}`;
    } else if (getInfosWidgetColumn(column)) {
      result.value = `${column}.${dictionary}.${field}`;
    } else if (isLinksWidgetColumn(column)) {
      result.onlyIcon = onlyIcon;

      if (field) {
        result.value = `${column}.${field}`;
      }
    }

    return result;
  })
);
