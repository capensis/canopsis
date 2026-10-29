import { ALARM_FIELDS } from '@/constants';

import { i18n } from '@/i18n';

import {
  getInfosWidgetColumn,
  isAlarmTicketDataWidgetColumn,
  isLinksWidgetColumn,
} from '@/helpers/entities/widget/column/form';

/**
 * Get translated label for widget column.
 *
 * @param {WidgetColumn} [column={}]
 * @param {Object<string, string>} [labelsMap={}]
 * @returns {string}
 */
export const getWidgetColumnLabel = (column = {}, labelsMap = {}) => {
  if (column.label) {
    return column.label;
  }

  const infosColumn = getInfosWidgetColumn(column.value);

  if (infosColumn) {
    return getTranslatedColumnLabel(labelsMap[infosColumn], column.value);
  }

  if (isLinksWidgetColumn(column.value)) {
    return getTranslatedColumnLabel(labelsMap[ALARM_FIELDS.links], column.value);
  }

  if (isAlarmTicketDataWidgetColumn(column.value)) {
    return getTranslatedColumnLabel(labelsMap[ALARM_FIELDS.ticketData], column.value);
  }

  return getTranslatedColumnLabel(labelsMap[column.value], column.value);
};

/**
 * Find a widget column corresponding to a sort field.
 * Infos sort templates contain a root field, while widget columns contain a concrete infos path.
 *
 * @param {WidgetColumn[]} [columns=[]]
 * @param {string} [sortBy='']
 * @returns {WidgetColumn|undefined}
 */
export const findWidgetColumnBySortBy = (columns = [], sortBy = '') =>
  columns.find(({ value }) =>
    (getInfosWidgetColumn(sortBy) ? value?.startsWith(sortBy) : value === sortBy),
  );

/**
 * Translate a widget column label with the current vue-i18n pluralization API.
 *
 * @param {string} key
 * @param {string} fallback
 * @returns {string}
 */
const getTranslatedColumnLabel = (key, fallback = '') => {
  if (!key) {
    return fallback;
  }

  return i18n.global.t(key, { count: 2 });
};

/**
 * Get sortable property for widget column for table.
 *
 * @param {WidgetColumn} [column={}]
 * @param {string[]} [unsortableFields=[]]
 * @returns {boolean}
 */
export const getWidgetColumnSortable = (column = {}, unsortableFields = []) =>
  !unsortableFields.some(field => (column.value ?? '').startsWith(field));
