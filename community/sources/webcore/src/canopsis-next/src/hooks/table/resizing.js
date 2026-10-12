import { throttle, sum } from 'lodash';
import {
  computed,
  set,
  ref,
  unref,
  onBeforeUnmount,
} from 'vue';

import { useHTMLElements } from '@/hooks/html-elements';

/**
 * Function to handle table resizing functionality.
 *
 * @param {Object} options - The options object.
 * @param {Array} options.headers - The array of table headers.
 * @param {HTMLElement} options.tableHeaderElement - The table header element.
 * @param {string} options.cellsSelector - The selector for table cells.
 * @param {string} options.resizableColumn - The resizable column.
 * @param {number} [options.throttleDelay = 10] - The throttle delay value.
 * @param {number} [options.minColumnWidth = 40] - The minimum column width.
 * @returns {Object} An object containing functions and reactive properties for table resizing.
 */
export const useTableResizing = ({
  headers,
  tableHeaderElement,
  cellsSelector,
  resizableColumn,
  throttleDelay = 10,
  minColumnWidth = 40,
}) => {
  const resizingMode = ref(false);
  const columnsWidthByFields = ref({});

  let resizingColumnIndex = null;
  let aggregatedMovementDiff = 0;

  const { elements: tableHeaderCells } = useHTMLElements({
    parentElement: tableHeaderElement,
    selector: cellsSelector,
  });

  const sumOfColumnsWidth = computed(() => sum(Object.values(columnsWidthByFields.value)));

  /**
   * HEADERS WITH ACTUAL WIDTHS
   */
  const headersWithWidth = computed(() => (
    !unref(resizableColumn)
      ? unref(headers)
      : unref(headers).map((header) => {
        const width = columnsWidthByFields.value[header.value];

        return {
          ...header,
          width: header.width
            ? header.width
            : width && `${width}px`,
        };
      })
  ));

  /**
   * Function to get normalized width
   *
   * @param {string} field
   * @param {number} newWidth
   * @returns {number}
   */
  const getNormalizedWidth = (field, newWidth) => Math.max(newWidth, unref(minColumnWidth));

  /**
   * Function to calculate element normalized width
   *
   * @param {HTMLElement} element
   * @param {string} field
   * @returns {number}
   */
  const calculateElementNormalizedWidth = (element, field) => {
    const { width: headerWidth } = element.getBoundingClientRect();

    return getNormalizedWidth(field, headerWidth);
  };

  /**
   * Function to set whole columns width by fields object
   *
   * @param {Object<string, number>} columnsWidth
   */
  const setColumnsWidth = columnsWidth => columnsWidthByFields.value = { ...columnsWidth };

  /**
   * Function to set column width by key
   *
   * @param {string} itemKey
   * @param {number} width
   */
  const setColumnsWidthItem = (itemKey, width) => set(columnsWidthByFields.value, itemKey, width);

  /**
   * Function to calculate all table columns widths
   */
  const calculateColumnsWidths = () => (
    setColumnsWidth([...tableHeaderCells.value].reduce((acc, headerElement) => {
      if (headerElement.dataset?.value) {
        const { value } = headerElement.dataset;

        acc[value] = calculateElementNormalizedWidth(headerElement, value);
      }

      return acc;
    }, {}))
  );

  /**
   * Function to get column width by field
   *
   * @param {string} field
   * @returns {*}
   */
  const getColumnWidthByField = field => columnsWidthByFields.value[field];

  /**
   * Function to enable resizing mode
   */
  const enableResizingMode = () => {
    resizingMode.value = true;

    calculateColumnsWidths();
  };

  /**
   * Function to disable resizing mode
   */
  const disableResizingMode = () => resizingMode.value = false;

  /**
   * Function to toggle resizing mode
   */
  const toggleResizingMode = () => (resizingMode.value ? disableResizingMode() : enableResizingMode());

  /**
   * Function to resize column by difference
   *
   * @param {number} index
   */
  const resizeColumnByDiff = (index) => {
    if (!aggregatedMovementDiff) {
      return;
    }

    const resizingLeftColumn = unref(headers)?.[index]?.value;
    const previousLeftColumnWidth = getColumnWidthByField(resizingLeftColumn);
    const newLeftColumnWidth = getNormalizedWidth(
      resizingLeftColumn,
      previousLeftColumnWidth + aggregatedMovementDiff,
    );

    setColumnsWidthItem(resizingLeftColumn, newLeftColumnWidth);

    if (newLeftColumnWidth !== previousLeftColumnWidth) {
      aggregatedMovementDiff = 0;
    }
  };

  /**
   * Throttled function to resize column by difference
   */
  const throttledResizeColumnByDiff = throttle(resizeColumnByDiff, throttleDelay);

  /**
   * Function to handle tick of column resize
   */
  const handleColumnResize = (event) => {
    aggregatedMovementDiff += event.movementX;

    throttledResizeColumnByDiff(resizingColumnIndex);
  };

  /**
   * Function to finish column resize
   */
  const finishColumnResize = () => {
    aggregatedMovementDiff = 0;

    const body = document.querySelector('body');

    if (!body) {
      return;
    }

    body.removeEventListener('mousemove', handleColumnResize);
    body.removeEventListener('mouseup', finishColumnResize);
    body.removeEventListener('mouseleave', finishColumnResize);
  };

  /**
   * Function to start column resize
   *
   * @param {string} columnName
   */
  const startColumnResize = (columnName) => {
    const body = document.querySelector('body');

    resizingColumnIndex = unref(headers)?.findIndex(({ value }) => value === columnName);

    if (!body) {
      return;
    }

    body.addEventListener('mousemove', handleColumnResize);
    body.addEventListener('mouseup', finishColumnResize);
    body.addEventListener('mouseleave', finishColumnResize);
  };

  onBeforeUnmount(finishColumnResize);

  return {
    resizingMode,
    columnsWidthByFields,
    sumOfColumnsWidth,
    headers: headersWithWidth,

    enableResizingMode,
    disableResizingMode,
    toggleResizingMode,
    setColumnsWidth,
    calculateColumnsWidths,
    startColumnResize,
    finishColumnResize,
  };
};
