import { throttle } from 'lodash';

export const widgetColumnResizingAlarmMixin = {
  props: {
    resizingColumnThrottleDelay: {
      type: Number,
      default: 10,
    },
  },
  data() {
    return {
      resizingMode: false,
      resizingColumnIndex: null,
      percentsInPixel: null,
      columnsWidthByField: {},
      columnsMinWidthByField: {},
      aggregatedMovementX: 0,
    };
  },
  created() {
    this.throttledResizeColumnByDiff = throttle(this.resizeColumnByDiff, this.resizingColumnThrottleDelay);
  },
  beforeDestroy() {
    this.finishColumnResize();
  },
  watch: {
    resizingMode() {
      this.calculateColumnsWidths();
    },
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
    enableResizingMode() {
      this.resizingMode = true;
    },

    disableResizingMode() {
      this.resizingMode = false;
    },

    toggleResizingMode() {
      return this.resizingMode ? this.disableResizingMode() : this.enableResizingMode();
    },

    setColumnsWidth(columnsWidth) {
      this.columnsWidthByField = { ...columnsWidth };
    },

    getColumnWidthByField(field) {
      return this.columnsWidthByField[field];
    },

    getColumnMinWidthByField(field) {
      return this.columnsMinWidthByField[field];
    },

    setPercentsInPixel() {
      const { width: rowWidth } = this.tableRow.getBoundingClientRect();

      this.percentsInPixel = 100 / rowWidth;
    },

    calculateColumnsWidths() {
      this.setPercentsInPixel();

      const { columnsWidthByField, columnsMinWidthByField } = [...this.headerCells].reduce((acc, headerElement) => {
        if (headerElement.dataset?.value) {
          const { value } = headerElement.dataset;
          const { width: headerWidth } = headerElement.getBoundingClientRect();
          const headerContentElement = headerElement.querySelector('span');
          const { width: contentWidth } = headerContentElement.getBoundingClientRect();

          const minWidth = (24 + 22 + 16 + contentWidth) * this.percentsInPixel;
          const width = headerWidth * this.percentsInPixel;

          acc.columnsWidthByField[value] = Math.max(minWidth, width);
          /**
           * 24 - max padding size
           * 22 - max sort position icon width
           * 16 - max sort direction icon width
           */
          acc.columnsMinWidthByField[value] = minWidth;
        }

        return acc;
      }, {
        columnsWidthByField: {},
        columnsMinWidthByField: {},
      });

      this.columnsWidthByField = columnsWidthByField;
      this.columnsMinWidthByField = columnsMinWidthByField;
    },

    getNormalizedWidth(field, newWidth) {
      return Math.max(newWidth, this.getColumnMinWidthByField(field));
    },

    resizeColumnByDiff(index) {
      const diff = this.aggregatedMovementX;

      const toRight = diff > 0;

      const resizingLeftColumn = this.headers[index].value;
      const resizingRightColumn = this.headers[index + 1].value;

      const previousLeftColumnWidth = this.getColumnWidthByField(resizingLeftColumn);
      const previousRightColumnWidth = this.getColumnWidthByField(resizingRightColumn);

      let newLeftColumnWidth = this.getNormalizedWidth(resizingLeftColumn, previousLeftColumnWidth + diff);
      let newRightColumnWidth = this.getNormalizedWidth(resizingRightColumn, previousRightColumnWidth - diff);

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
            const newAffectedHeaderWidth = this.getNormalizedWidth(value, affectedHeaderWidth - remainderDiff);

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

      this.columnsWidthByField = {
        ...this.columnsWidthByField,
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
