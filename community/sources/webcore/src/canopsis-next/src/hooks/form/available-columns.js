import { computed, unref } from 'vue';

import {
  ENTITIES_TYPES,
  ALARM_FIELDS,
  ALARM_LIST_WIDGET_GROUPED_COLUMNS,
  CONTEXT_WIDGET_GROUPED_COLUMNS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';

/**
 * Hook for getting available columns for form fields
 *
 * @param {Object} options
 * @param {string} [options.type] - Entity type (alarm or entity)
 * @param {Array} [options.items] - Custom items array (if provided, returns this)
 * @param {Array} [options.excludedColumns] - Array of column values to exclude
 * @param {boolean} [options.withInstructions] - Whether to include instructions column
 * @returns {Object} - Object with availableColumns computed property
 */
export const useAvailableColumns = ({
  type = ENTITIES_TYPES.alarm,
  items,
  excludedColumns = [],
  withInstructions = false,
}) => {
  const { t, tc } = useI18n();

  /**
   * Convert grouped columns object to flat array with headers
   *
   * @param {Object} [groupedColumns={}] - Object where keys are group names and values are arrays of column values
   * @param {Object} [keys={}] - Object mapping column values to i18n translation keys
   * @returns {Array<{header: string} | {value: string, text: string}>} - Array of column items with group headers
   */
  const groupedColumnsToColumns = (groupedColumns = {}, keys = {}) => (
    Object.entries(groupedColumns).reduce((acc, [group, groupItems]) => {
      const preparedItems = groupItems.map(value => ({
        value,
        text: tc(keys[value], 2),
      }));

      if (!preparedItems.length) {
        return acc;
      }

      acc.push(
        { header: t(`settings.columnsGroups.${group}`) },
        ...preparedItems,
      );

      return acc;
    }, [])
  );

  const isAlarmType = computed(() => unref(type) === ENTITIES_TYPES.alarm);

  const alarmListAvailableColumns = computed(() => {
    const columns = groupedColumnsToColumns(ALARM_LIST_WIDGET_GROUPED_COLUMNS, ALARM_FIELDS_TO_LABELS_KEYS);
    const withInstructionsValue = unref(withInstructions);

    return withInstructionsValue
      ? columns.filter(({ value }) => value !== ALARM_FIELDS.assignedInstructions)
      : columns;
  });

  const contextAvailableColumns = computed(() => (
    groupedColumnsToColumns(CONTEXT_WIDGET_GROUPED_COLUMNS, ENTITY_FIELDS_TO_LABELS_KEYS)
  ));

  const availableColumns = computed(() => {
    const itemsValue = unref(items);
    if (itemsValue) {
      return itemsValue;
    }

    const excludedColumnsValue = unref(excludedColumns);
    const columns = isAlarmType.value
      ? alarmListAvailableColumns.value
      : contextAvailableColumns.value;

    return columns.filter(({ value }) => !excludedColumnsValue.includes(value));
  });

  return {
    availableColumns,
  };
};
