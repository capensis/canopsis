import Sortable from 'sortablejs';
import { computed, ref, unref, onBeforeUnmount } from 'vue';

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
export const useTableDragging = ({
  tableHeaderElement,
  headers,
  draggableClass,
  draggableColumn,
}) => {
  const draggingMode = ref(false);
  const columnsPositionByFields = ref({});

  let sortableInstance;
  let draggingField = null;
  let lastPointerX = null;
  let dragDirection = 0;

  /**
   * HEADERS SORTED BY DRAGGING POSITIONS
   */
  const sortedHeaders = computed(() =>
    unref(draggableColumn)
      ? [...unref(headers)].sort(
          (a, b) => columnsPositionByFields.value[a.value] - columnsPositionByFields.value[b.value],
        )
      : unref(headers),
  );

  /**
   * Set columns position based on provided columns position.
   *
   * @param {Object} columnsPosition - The columns position object.
   */
  const setColumnsPosition = columnsPosition =>
    (columnsPositionByFields.value = { ...columnsPosition });

  /**
   * Calculate the sorting order of columns based on the provided headers for calculation.
   *
   * @param {Object[]} headersForCalculation - The array of headers to calculate the sorting order.
   */
  const calculateColumnsSortingOrderByHeaders = headersForCalculation =>
    setColumnsPosition(
      headersForCalculation.reduce((acc, { value }, index) => {
        acc[value] = index;

        return acc;
      }, {}),
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
   * Resolve the column field from a header cell or a nested drag handle.
   *
   * @param {HTMLElement} element
   * @returns {string|null}
   */
  const getColumnValueFromElement = element =>
    element?.dataset?.value
    ?? element?.closest?.('th')?.dataset?.value
    ?? element?.querySelector?.('[data-value]')?.dataset?.value
    ?? element?.closest?.('th')?.querySelector?.('[data-value]')?.dataset?.value
    ?? null;

  /**
   * Handle the sorting of columns when dragging.
   * Vue re-renders the whole column (content + width). Sortable must not move DOM.
   * Only move in the pointer direction so a left/right drag can pass every column
   * without snapping back or skipping the next one.
   *
   * @param {Object} evt - Sortable onMove event.
   * @param {MouseEvent} originalEvent - Native pointer event.
   * @returns {false}
   */
  const handleColumnSort = (evt, originalEvent) => {
    const clientX = originalEvent?.clientX;

    if (clientX != null) {
      if (lastPointerX != null) {
        const pointerDeltaX = clientX - lastPointerX;

        if (Math.abs(pointerDeltaX) > 1) {
          dragDirection = pointerDeltaX < 0 ? -1 : 1;
        }
      }

      lastPointerX = clientX;
    }

    const draggedValue = draggingField ?? getColumnValueFromElement(evt.dragged);
    const relatedValue = getColumnValueFromElement(evt.related);

    if (!draggedValue || !relatedValue || draggedValue === relatedValue) {
      return false;
    }

    const oldDraggableIndex = sortedHeaders.value.findIndex(
      ({ value }) => value === draggedValue,
    );
    const newDraggableIndex = sortedHeaders.value.findIndex(
      ({ value }) => value === relatedValue,
    );

    if (
      oldDraggableIndex === -1
      || newDraggableIndex === -1
      || oldDraggableIndex === newDraggableIndex
    ) {
      return false;
    }

    if (dragDirection < 0 && newDraggableIndex > oldDraggableIndex) {
      return false;
    }

    if (dragDirection > 0 && newDraggableIndex < oldDraggableIndex) {
      return false;
    }

    const copiedHeaders = getHeadersByMovingIndexes(oldDraggableIndex, newDraggableIndex);

    calculateColumnsSortingOrderByHeaders(copiedHeaders);

    return false;
  };

  /**
   * Start column dragging by creating a Sortable instance.
   */
  const startColumnDragging = () => {
    calculateColumnsSortingOrderByHeaders(sortedHeaders.value);

    const headerRow = unref(tableHeaderElement);

    if (!headerRow) {
      return;
    }

    sortableInstance = Sortable.create(headerRow, {
      draggable: `th:has(.table__dragging-handler), th.${unref(draggableClass)}`,
      onStart: ({ item }) => {
        draggingField = getColumnValueFromElement(item);
        lastPointerX = null;
        dragDirection = 0;
      },
      onMove: handleColumnSort,
      onEnd: () => {
        draggingField = null;
        lastPointerX = null;
        dragDirection = 0;
      },
      direction: 'horizontal',
      forceFallback: true,
    });
  };

  /**
   * Finish column dragging by destroying the sortable instance.
   */
  const finishColumnDragging = () => {
    sortableInstance?.destroy();
    sortableInstance = null;
    draggingField = null;
    lastPointerX = null;
    dragDirection = 0;
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
  const toggleDraggingMode = () =>
    draggingMode.value ? disableDraggingMode() : enableDraggingMode();

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
