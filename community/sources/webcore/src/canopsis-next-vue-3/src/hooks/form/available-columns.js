import { computed, unref } from 'vue';

import {
  ALARM_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_LIST_WIDGET_GROUPED_COLUMNS,
  CONTEXT_WIDGET_GROUPED_COLUMNS,
  ENTITIES_TYPES,
  ENTITY_FIELDS_TO_LABELS_KEYS,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';

/**
 * Hook for getting available columns for form fields.
 *
 * @param {Object} options
 * @param {string|import('vue').Ref<string>} [options.type] - Entity type
 * @param {Array|import('vue').Ref<Array>} [options.items] - Custom items array
 * @param {Array|import('vue').Ref<Array>} [options.excludedColumns] - Excluded column values
 * @param {boolean|import('vue').Ref<boolean>} [options.withInstructions] - Include assigned instructions column
 * @returns {{ availableColumns: import('vue').ComputedRef<Array> }}
 */
export const useAvailableColumns = ({
  type = ENTITIES_TYPES.alarm,
  items,
  excludedColumns = [],
  withInstructions = false,
}) => {
  const { t } = useI18n();

  /**
   * Convert grouped columns object to flat array with headers.
   *
   * @param {Object} [groupedColumns={}]
   * @param {Object<string, string>} [keys={}]
   * @returns {Array<{header: string} | {value: string, text: string}>}
   */
  const groupedColumnsToColumns = (groupedColumns = {}, keys = {}) => Object.entries(groupedColumns).reduce((acc, [group, groupItems]) => {
      const preparedItems = groupItems.map(value => ({
        value,
        text: t(keys[value], 2),
      }));

      if (!preparedItems.length) {
        return acc;
      }

      acc.push({ header: t(`settings.columnsGroups.${group}`) }, ...preparedItems);

      return acc;
    }, []);

  const isAlarmType = computed(() => unref(type) === ENTITIES_TYPES.alarm);

  const alarmListAvailableColumns = computed(() => {
    const columns = groupedColumnsToColumns(
      ALARM_LIST_WIDGET_GROUPED_COLUMNS,
      ALARM_FIELDS_TO_LABELS_KEYS,
    );

    return unref(withInstructions)
      ? columns.filter(({ value }) => value !== ALARM_FIELDS.assignedInstructions)
      : columns;
  });

  const contextAvailableColumns = computed(() => groupedColumnsToColumns(CONTEXT_WIDGET_GROUPED_COLUMNS, ENTITY_FIELDS_TO_LABELS_KEYS),
  );

  const availableColumns = computed(() => {
    const unwrappedItems = unref(items);

    if (unwrappedItems?.length) {
      return unwrappedItems;
    }

    const columns = isAlarmType.value
      ? alarmListAvailableColumns.value
      : contextAvailableColumns.value;
    const unwrappedExcludedColumns = unref(excludedColumns) ?? [];

    return columns.filter(({ value }) => !value || !unwrappedExcludedColumns.includes(value));
  });

  return {
    availableColumns,
  };
};
