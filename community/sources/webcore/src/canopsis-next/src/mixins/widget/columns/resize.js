export const widgetResizeAlarmMixin = {
  data() {
    return {
      resizingColumnIndex: null,
      rowWidth: null,
      columnWidthByField: {},
    };
  },
  beforeDestroy() {
    this.finishColumnResize();
  },
  methods: {
    getColumnWidthByField(field) {
      return this.columnWidthByField[field];
    },

    calculateColumnsWidthsByTable(tableHeaderElement) {
      const tableRow = tableHeaderElement.querySelector('tr:first-of-type');
      const headerCells = tableRow.querySelectorAll('th');

      const { width: rowWidth } = tableRow.getBoundingClientRect();

      this.rowWidth = rowWidth;

      this.columnWidthByField = [...headerCells].reduce((acc, headerElement) => {
        if (headerElement.dataset?.value) {
          const { width } = headerElement.getBoundingClientRect();

          acc[headerElement.dataset.value] = (width / this.rowWidth) * 100;
        }

        return acc;
      }, {});
    },

    getNormalizedWidth(newWidth) {
      return Math.max(newWidth, 8);
    },

    resizeColumnByDiff(index, diff) {
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
    },

    handleColumnResize(event) {
      const diff = (event.movementX / this.rowWidth) * 100;

      this.resizeColumnByDiff(this.resizingColumnIndex, diff);
    },

    finishColumnResize() {
      document.body.removeEventListener('mousemove', this.handleColumnResize);
      document.body.removeEventListener('mouseup', this.finishColumnResize);
    },

    startColumnResize(columnName) {
      this.resizingColumnIndex = this.headers.findIndex(({ value }) => value === columnName);

      document.body.addEventListener('mousemove', this.handleColumnResize);
      document.body.addEventListener('mouseup', this.finishColumnResize);
    },
  },
};
