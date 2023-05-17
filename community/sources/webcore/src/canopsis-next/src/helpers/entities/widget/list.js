import { ALARM_FIELDS } from '@/constants';

import i18n from '@/i18n';

import { getInfosWidgetColumn, isLinksWidgetColumn } from '@/helpers/entities/widget/column/form';

/**
 * Get translated label for widget column
 *
 * @param {WidgetColumn} [column = {}]
 * @param {Object<string, string>} [labelsMap = {}]
 * @returns {string}
 */
export const getWidgetColumnLabel = (column = {}, labelsMap = {}) => {
  if (column.label) {
    return column.label;
  }

  const infosColumn = getInfosWidgetColumn(column.value);

  if (infosColumn) {
    return i18n.tc(labelsMap[infosColumn], 2);
  }

  if (isLinksWidgetColumn(column.value)) {
    return i18n.tc(labelsMap[ALARM_FIELDS.links], 2);
  }

  return i18n.tc(labelsMap[column.value], 2);
};

/**
 * Get sortable property for widget column for table
 *
 * @param {WidgetColumn} [column = {}]
 * @param {string[]} [unsortableFields = []]
 * @returns {boolean}
 */
export const getWidgetColumnSortable = (column = {}, unsortableFields = []) => (
  !unsortableFields.some(field => (column.value ?? '').startsWith(field))
);
