import { map, pick } from 'lodash';
import { computed, ref, unref } from 'vue';

import { EXTERNAL_DATA_TABLE_COLUMN_TAGS } from '@/constants';

import { externalDataTableColumnConfigsToForm } from '@/helpers/entities/external-data-table/form';

import { usePendingHandler } from '@/hooks/query/pending';
import { useExternalDataTable } from '@/hooks/store/modules/external-data-table';
import { useWidgetUserPreference } from '@/hooks/store/modules/user-preference';

/**
 * Hook that provides functionality for managing an external data table widget
 *
 * @param {Object} options - Configuration options
 * @param {Object|Ref<Object>} options.widget - The widget object containing configuration
 * @returns {Object} Object containing external data table state and methods
 * @property {ComputedRef<Object>} userPreference - User preferences for this widget
 * @property {Function} fetchUserPreference - Function to fetch user preferences
 * @property {Function} updateContentInUserPreference - Function to update user preference content
 * @property {Ref<Object>} externalDataTable - The external data table object
 * @property {ComputedRef<Object>} columns - Computed getter/setter for table columns
 * @property {Ref<boolean>} pending - Whether data is currently being fetched
 * @property {Function} fetchItem - Function to fetch the external data table
 */
export const useExternalDataTableWidgetTable = ({
  widget,
}) => {
  const widgetId = computed(() => unref(widget)._id);

  const {
    userPreference,
    fetchUserPreference,
    updateContentInUserPreference,
  } = useWidgetUserPreference({ widgetId });

  const {
    updateExternalDataTable,
    fetchExternalDataTableWithoutStore,
  } = useExternalDataTable();

  const externalDataTable = ref({});

  const columns = computed({
    get: () => {
      const wholeColumns = externalDataTableColumnConfigsToForm(
        externalDataTable.value.column_configs,
      );

      return pick({
        ...wholeColumns,
        _id: EXTERNAL_DATA_TABLE_COLUMN_TAGS.noType,
      }, map(unref(widget).parameters.widgetColumns, 'value'));
    },

    set: async (newColumns) => {
      externalDataTable.value = await updateExternalDataTable({
        id: externalDataTable.value._id,
        data: {
          ...pick(externalDataTable.value, ['type', 'name', 'description']),

          column_types: externalDataTable.value.columns.map((column, index) => (
            newColumns[column] ?? externalDataTable.value.column_types[index]
          )),
        },
      });
    },
  });

  const { pending, handler: fetchItem } = usePendingHandler(async () => {
    const [response] = await Promise.all([
      fetchExternalDataTableWithoutStore({ id: unref(widget).parameters.table }),
      fetchUserPreference(),
    ]);

    externalDataTable.value = response;
  }, true);

  return {
    userPreference,
    fetchUserPreference,
    updateContentInUserPreference,

    externalDataTable,

    columns,

    pending,
    fetchItem,
  };
};

/**
 * Hook that provides column settings functionality for an external data table widget
 *
 * @param {Object} options - Configuration options
 * @param {Object|Ref<Object>} options.widget - The widget object containing configuration
 * @param {ComputedRef<Object>} options.userPreference - User preferences for this widget
 * @param {Function} options.updateContentInUserPreference - Function to update user preference content
 * @returns {Object} Object containing column settings state and methods
 * @property {ComputedRef<boolean>} draggableColumn - Whether columns can be dragged/reordered
 * @property {ComputedRef<boolean>} resizableColumn - Whether columns can be resized
 * @property {ComputedRef<string>} cellsContentBehavior - How cell content behaves when resized
 * @property {ComputedRef<Object>} columnsSettings - Current column settings from user preferences
 * @property {Function} updateColumnsSettings - Function to update column settings
 */
export const useExternalDataTableWidgetColumnsSettings = ({
  widget,
  userPreference,
  updateContentInUserPreference,
}) => {
  const draggableColumn = computed(() => unref(widget).parameters?.columns?.draggable);
  const resizableColumn = computed(() => unref(widget).parameters?.columns?.resizable);
  const cellsContentBehavior = computed(() => unref(widget).parameters?.columns?.cells_content_behavior);
  const columnsSettings = computed(() => unref(userPreference).content.columns_settings ?? {});

  const updateColumnsSettings = newColumnsSettings => (
    updateContentInUserPreference({ columns_settings: newColumnsSettings })
  );

  return {
    draggableColumn,
    resizableColumn,
    cellsContentBehavior,
    columnsSettings,
    updateColumnsSettings,
  };
};

/**
 * Hook that provides dense mode functionality for an external data table widget
 *
 * @param {Object} options - Configuration options
 * @param {Object|Ref<Object>} options.widget - The widget object containing configuration
 * @param {ComputedRef<Object>} options.userPreference - User preferences for this widget
 * @param {Function} options.updateContentInUserPreference - Function to update user preference content
 * @returns {Object} Object containing dense mode state and methods
 * @property {ComputedRef<boolean>} dense - Whether dense mode is enabled
 * @property {Function} updateDense - Function to update dense mode setting
 */
export const useExternalDataTableWidgetDense = ({
  widget,
  userPreference,
  updateContentInUserPreference,
}) => {
  const dense = computed(() => unref(userPreference).content.dense ?? unref(widget).parameters.dense);
  const updateDense = newDense => updateContentInUserPreference({ dense: newDense });

  return {
    dense,
    updateDense,
  };
};
