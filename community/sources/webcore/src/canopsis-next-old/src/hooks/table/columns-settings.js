import { computed, unref, nextTick, watch } from 'vue';

import { RESIZING_CELLS_CONTENTS_BEHAVIORS } from '@/constants';

import { useTableDragging } from './dragging';
import { useTableResizing } from './resizing';

/**
 * Hook that provides table column settings functionality including dragging and resizing capabilities
 *
 * This hook combines the functionality of useTableDragging and useTableResizing to provide
 * a complete solution for managing table column settings. It handles:
 * - Column dragging (reordering)
 * - Column resizing
 * - Synchronization with external column settings
 * - Emitting updated settings to parent components
 *
 * @param {Object} options - Configuration options for the table columns settings
 * @param {Ref<HTMLElement>} options.tableHeaderElement - Reference to the table header element
 * @param {Ref<Array>} options.headers - Array of table header objects
 * @param {Ref<boolean>} options.draggable - Whether columns can be dragged/reordered
 * @param {Ref<boolean>} options.resizable - Whether columns can be resized
 * @param {Ref<Object>} options.columnsSettings - External column settings (width and position)
 * @param {Ref<string>} options.cellsContentBehavior - How cell content behaves when resized ('wrap' or 'truncate')
 * @param {string} [options.draggableClass='draggable'] - CSS class for draggable elements
 * @param {Function} emit - Vue emit function for emitting events to parent component
 *
 * @returns {Object} Object containing table column settings state and methods
 * @property {ComputedRef<Object>} tableStyle - Computed styles for the table
 * @property {ComputedRef<Object>} tableClass - Computed CSS classes for the table
 * @property {ComputedRef<boolean>} isColumnsEditing - Whether columns are currently being edited (dragged or resized)
 * @property {Ref<boolean>} draggingMode - Whether dragging mode is active
 * @property {Ref<boolean>} resizingMode - Whether resizing mode is active
 * @property {ComputedRef<Array>} sortedHeaders - Headers sorted according to dragging positions
 * @property {ComputedRef<Array>} sortedHeadersWithWidth - Headers sorted and with width properties
 * @property {Function} startColumnResize - Function to start resizing a column
 * @property {Function} toggleColumnEditingMode - Function to toggle column editing mode
 * @property {Function} resetColumnsSettings - Function to reset column settings to defaults
 */
export const useTableColumnsSettings = ({
  tableHeaderElement,
  headers,
  draggable,
  resizable,
  columnsSettings,
  cellsContentBehavior,
  draggableClass = 'draggable',
}, emit) => {
  /**
   * DRAGGABLE
   */
  const {
    draggingMode,
    columnsPositionByFields,
    headers: sortedHeaders,
    setColumnsPosition,
    toggleDraggingMode,
  } = useTableDragging({
    tableHeaderElement,
    headers,
    draggableClass,
    draggableColumn: draggable,
  });

  /**
   * RESIZABLE
   */
  const {
    resizingMode,
    sumOfColumnsWidth,
    columnsWidthByFields,
    headers: sortedHeadersWithWidth,
    setColumnsWidth,
    toggleResizingMode,
    calculateColumnsWidths,
    startColumnResize,
  } = useTableResizing({
    headers: sortedHeaders,
    tableHeaderElement,
    resizableColumn: resizable,
    cellsSelector: 'th',
  });

  /**
   * TABLE COMPUTED PROPERTIES
   */
  const isColumnsEditing = computed(() => draggingMode.value || resizingMode.value);

  const tableStyle = computed(() => (
    unref(resizable)
      ? { '--external-data-table-width': sumOfColumnsWidth.value ? `${sumOfColumnsWidth.value}px` : '100%' }
      : {}
  ));

  const tableClass = computed(() => {
    const unwrappedCellsContentBehavior = unref(cellsContentBehavior);

    return {
      table__grid: isColumnsEditing.value,
      'table--wrapped': unwrappedCellsContentBehavior === RESIZING_CELLS_CONTENTS_BEHAVIORS.wrap,
      'table--truncated': unwrappedCellsContentBehavior === RESIZING_CELLS_CONTENTS_BEHAVIORS.truncate,
      'external-data-table--fixed': unref(resizable),
    };
  });

  /**
   * Updates the column settings based on current resizing and dragging modes
   * Collects width and position settings and emits them to the parent component
   * @emits update:columns-settings - Emits an object containing columns width and/or position settings
   */
  const updateColumnsSettings = () => {
    const settings = {};

    if (resizingMode.value) {
      settings.columns_width = columnsWidthByFields.value;
    }

    if (draggingMode.value) {
      settings.columns_position = columnsPositionByFields.value;
    }

    emit('update:columns-settings', settings);
  };

  /**
   * Toggles column editing mode on/off
   * If editing is active, updates column settings before toggling
   * Toggles dragging and resizing modes based on component props
   */
  const toggleColumnEditingMode = () => {
    if (isColumnsEditing.value) {
      updateColumnsSettings();
    }

    if (unref(draggable)) {
      toggleDraggingMode();
    }

    if (unref(resizable)) {
      toggleResizingMode();
    }
  };

  /**
   * Resets all column settings to their default values
   * Clears column positions if draggable is enabled
   * Resets column widths and recalculates them if resizable is enabled
   */
  const resetColumnsSettings = () => {
    if (unref(draggable)) {
      setColumnsPosition({});
    }

    if (unref(resizable)) {
      setColumnsWidth({});
      nextTick(calculateColumnsWidths);
    }
  };

  /**
   * Watches for changes to columnsSettings and updates column positions and widths accordingly
   */
  watch(columnsSettings, (newColumnsSettings) => {
    if (!draggingMode.value && newColumnsSettings?.columns_position) {
      setColumnsPosition(newColumnsSettings?.columns_position);
    }

    if (!resizingMode.value && newColumnsSettings?.columns_width) {
      setColumnsWidth(newColumnsSettings?.columns_width);
    }
  }, { immediate: true, deep: true });

  return {
    tableStyle,
    tableClass,
    isColumnsEditing,
    draggingMode,
    resizingMode,
    sortedHeaders,
    sortedHeadersWithWidth,

    startColumnResize,
    toggleColumnEditingMode,
    resetColumnsSettings,
  };
};
