import Sortable from 'sortablejs';
import { throttle } from 'lodash';

export const widgetResizeAlarmMixin = {
  data() {
    return {
      resizingMode: false,
      resizingColumnIndex: null,
      percentsInPixel: null,
      columnWidthByField: {},
      columnOrderByField: {},
      aggregatedMovementX: 0,
    };
  },
  created() {
    this.throttledResizeColumnByDiff = throttle(this.resizeColumnByDiff, 10);
  },
  beforeDestroy() {
    this.finishColumnResize();
    this.finishColumnSortable();
  },
  computed: {
    tableRow() {
      return this.tableHeader.querySelector('tr:first-of-type');
    },

    headerCells() {
      return this.tableRow.querySelectorAll('th');
    },
  },
  methods: {
    toggleResizingMode() {
      this.resizingMode = !this.resizingMode;

      this.calculateColumnsWidths();
      this.startColumnSortable();
    },

    calculateInitialSorting() {
      this.columnOrderByField = this.headers.reduce((acc, { value }, index) => {
        acc[value] = index;

        return acc;
      }, {});
    },

    getColumnPositionByField(field) {
      return this.columnOrderByField[field];
    },

    handleColumnSort({ related, dragged }) {
      const draggedItemIndex = this.headers.findIndex(({ value }) => value === dragged.dataset.value);
      const relatedItemIndex = this.headers.findIndex(({ value }) => value === related.dataset.value);

      const copiedHeaders = [...this.headers];

      const [item] = copiedHeaders.splice(draggedItemIndex, 1);

      copiedHeaders.splice(relatedItemIndex, 0, item);

      this.columnOrderByField = copiedHeaders.reduce((acc, { value }, index) => {
        acc[value] = index;

        return acc;
      }, {});
    },

    startColumnSortable() {
      this.calculateInitialSorting();

      this.sortableInstance = Sortable.create(this.tableRow, {
        draggable: '.alarms-list-table__draggable-column',
        onMove: this.handleColumnSort,
        direction: 'horizontal',
      });
    },

    finishColumnSortable() {
      this.sortableInstance?.destroy();
      this.sortableInstance = null;
    },

    getColumnWidthByField(field) {
      return this.columnWidthByField[field];
    },

    setPercentsInPixel() {
      const { width: rowWidth } = this.tableRow.getBoundingClientRect();

      this.percentsInPixel = 100 / rowWidth;
    },

    calculateColumnsWidths() {
      this.setPercentsInPixel();

      this.columnWidthByField = [...this.headerCells].reduce((acc, headerElement) => {
        if (headerElement.dataset?.value) {
          const { width } = headerElement.getBoundingClientRect();

          acc[headerElement.dataset.value] = width * this.percentsInPixel;
        }

        return acc;
      }, {});
    },

    getNormalizedWidth(newWidth) {
      return Math.max(newWidth, 8);
    },

    resizeColumnByDiff(index) {
      const diff = this.aggregatedMovementX;

      const toRight = diff > 0;

      const resizingLeftColumn = this.headers[index].value;
      const resizingRightColumn = this.headers[index + 1].value;

      const previousLeftColumnWidth = this.getColumnWidthByField(resizingLeftColumn);
      const previousRightColumnWidth = this.getColumnWidthByField(resizingRightColumn);

      let newLeftColumnWidth = this.getNormalizedWidth(previousLeftColumnWidth + diff);
      let newRightColumnWidth = this.getNormalizedWidth(previousRightColumnWidth - diff);

      const resultLeftDiff = Math.abs(newLeftColumnWidth - previousLeftColumnWidth);
      const resultRightDiff = Math.abs(newRightColumnWidth - previousRightColumnWidth);

      let remainderDiff = Math.abs(resultLeftDiff - resultRightDiff);

      const affectedHeadersWidths = {};

      if (remainderDiff > 0) {
        const affectedHeaders = toRight
          /**
           * We need check each column from right cell to last
           */
          ? this.headers.slice(this.resizingColumnIndex + 2)
          /**
           * We need check each column from first to current column
           */
          : this.headers.slice(0, this.resizingColumnIndex).reverse();

        /**
         * Try to find free space for remaining diff
         */
        for (const { value } of affectedHeaders) {
          /**
           * We can have first header for expand button and remediation icon.
           * This header doesn't have value and cannot to resize.
           */
          if (value) {
            const affectedHeaderWidth = this.getColumnWidthByField(value);
            const newAffectedHeaderWidth = this.getNormalizedWidth(affectedHeaderWidth - remainderDiff);

            affectedHeadersWidths[value] = newAffectedHeaderWidth;

            remainderDiff -= (affectedHeaderWidth - newAffectedHeaderWidth);

            if (remainderDiff <= 0) {
              break;
            }
          }
        }
      }

      /**
       * Normalize width, if we don't have available space
       */
      if (remainderDiff) {
        if (toRight) {
          newLeftColumnWidth -= remainderDiff;
        } else {
          newRightColumnWidth -= remainderDiff;
        }
      }

      this.columnWidthByField = {
        ...this.columnWidthByField,
        ...affectedHeadersWidths,
        [resizingLeftColumn]: newLeftColumnWidth,
        [resizingRightColumn]: newRightColumnWidth,
      };
      this.aggregatedMovementX = 0;
    },

    handleColumnResize(event) {
      this.aggregatedMovementX += event.movementX * this.percentsInPixel;

      this.throttledResizeColumnByDiff(this.resizingColumnIndex);
    },

    finishColumnResize() {
      document.body.removeEventListener('mousemove', this.handleColumnResize);
      document.body.removeEventListener('mouseup', this.finishColumnResize);
      document.body.addEventListener('mouseleave', this.finishColumnResize);
    },

    startColumnResize(columnName) {
      this.resizingColumnIndex = this.headers.findIndex(({ value }) => value === columnName);

      document.body.addEventListener('mousemove', this.handleColumnResize);
      document.body.addEventListener('mouseup', this.finishColumnResize);
      document.body.addEventListener('mouseleave', this.finishColumnResize);
    },
  },
};
