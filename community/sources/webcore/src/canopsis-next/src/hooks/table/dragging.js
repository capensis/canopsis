import { computed, ref, unref, onBeforeUnmount } from 'vue';
import Sortable from 'sortablejs';

/**
 * Custom hook for enabling table column dragging functionality.
 *
 * @param {Object} options - The options object.
 * @param {Element} options.tableHeaderElement - Table header element.
 * @param {Object[]} options.headers - The array of table headers.
 * @param {string} options.draggableClass - The class name for draggable columns.
 * @param {string} options.draggableColumn - The boolean property for enabling or disabling dragging.
 * @returns {Object} An object containing functions and reactive values for table dragging.
 */
export const useTableDragging = ({ tableHeaderElement, headers, draggableClass, draggableColumn }) => {
  const draggingMode = ref(false);
  const columnsPositionByFields = ref({});

  let sortableInstance;

  /**
   * HEADERS SORTED BY DRAGGING POSITIONS
   */
  const sortedHeaders = computed(() => (
    unref(draggableColumn)
      ? [...unref(headers)].sort((a, b) => (
        columnsPositionByFields.value[a.value] - columnsPositionByFields.value[b.value]
      ))
      : unref(headers)
  ));

  /**
   * Set columns position based on provided columns position.
   *
   * @param {Object} columnsPosition - The columns position object.
   */
  const setColumnsPosition = columnsPosition => columnsPositionByFields.value = ({ ...columnsPosition });

  /**
   * Calculate the sorting order of columns based on the provided headers for calculation.
   *
   * @param {Object[]} headersForCalculation - The array of headers to calculate the sorting order.
   */
  const calculateColumnsSortingOrderByHeaders = headersForCalculation => (
    setColumnsPosition(headersForCalculation.reduce((acc, { value }, index) => {
      acc[value] = index;

      return acc;
    }, {}))
  );

  /**
   * Generate JSDoc for the function `getHeadersByMovingIndexes`.
   *
   * @param {number} oldIndex - The old index of the header.
   * @param {number} newIndex - The new index where the header will be moved.
   * @returns {Array} The updated array of headers after moving the item.
   */
  const getHeadersByMovingIndexes = (oldIndex, newIndex) => {
    const copiedHeaders = [...sortedHeaders.value];

    const [item] = copiedHeaders.splice(oldIndex, 1);

    copiedHeaders.splice(newIndex, 0, item);

    return copiedHeaders;
  };

  /**
   * Handle the sorting of columns when dragging.
   *
   * @param {Object} params - The parameters object.
   * @param {HTMLElement} params.dragged - The dragged column element.
   * @param {HTMLElement} params.related - The related column element.
   */
  const handleColumnSort = ({ dragged, related }) => {
    const oldDraggableIndex = sortedHeaders.value.findIndex(({ value }) => value === dragged.dataset.value);
    const newDraggableIndex = sortedHeaders.value.findIndex(({ value }) => value === related.dataset.value);
    const copiedHeaders = getHeadersByMovingIndexes(oldDraggableIndex, newDraggableIndex);

    calculateColumnsSortingOrderByHeaders(copiedHeaders);
  };

  /**
   * Start column dragging by creating a Sortable instance.
   */
  const startColumnDragging = () => {
    calculateColumnsSortingOrderByHeaders(sortedHeaders.value);

    sortableInstance = Sortable.create(unref(tableHeaderElement), {
      draggable: `.${unref(draggableClass)}`,
      onMove: handleColumnSort,
      direction: 'horizontal',
    });
  };

  /**
   * Finish column dragging by destroying the sortable instance.
   */
  const finishColumnDragging = () => {
    sortableInstance?.destroy();
    sortableInstance = null;
  };

  /**
   * Enable dragging mode for table columns.
   */
  const enableDraggingMode = () => {
    draggingMode.value = true;
    startColumnDragging();
  };

  /**
   * Disable dragging mode for table columns.
   */
  const disableDraggingMode = () => {
    draggingMode.value = false;
    finishColumnDragging();
  };

  /**
   * Toggle dragging mode for table columns.
   */
  const toggleDraggingMode = () => (draggingMode.value ? disableDraggingMode() : enableDraggingMode());

  /**
   * Cleanup function to destroy sortable instance on component unmount.
   */
  onBeforeUnmount(finishColumnDragging);

  return {
    draggingMode,
    columnsPositionByFields,
    headers: sortedHeaders,

    setColumnsPosition,
    enableDraggingMode,
    disableDraggingMode,
    toggleDraggingMode,
  };
};
