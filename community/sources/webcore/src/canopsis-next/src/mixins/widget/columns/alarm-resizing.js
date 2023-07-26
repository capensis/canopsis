import { throttle } from 'lodash';

export const widgetColumnResizingAlarmMixin = {
  props: {
    resizingColumnThrottleDelay: {
      type: Number,
      default: 10,
    },
    minColumnWidth: {
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
      aggregatedMovementDiff: 0,
    };
  },
  created() {
    this.throttledResizeColumnByDiff = throttle(this.resizeColumnByDiff, this.resizingColumnThrottleDelay);
  },
  beforeDestroy() {
    this.finishColumnResize();
  },
  computed: {
    tableRow() {
      return this.tableHeader.querySelector('tr:first-of-type');
    },

    headerCells() {
      return this.tableRow.querySelectorAll('th');
    },

    sumOfColumnsWidth() {
      return this.calculateFullColumnsWidth(this.columnsWidthByField);
    },

    minColumnsWidth() {
      /**
       * 24 - max padding size
       * 22 - max sort position icon width
       * 16 - max sort direction icon width
       */
      return (24 + 22 + 16 + this.minColumnWidth);
    },

    minColumnsWidthInPercent() {
      return this.minColumnsWidth * this.percentsInPixel;
    },
  },
  methods: {
    enableResizingMode() {
      this.resizingMode = true;

      this.calculateColumnsWidths();
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

    getNormalizedWidth(field, newWidth) {
      return Math.max(newWidth, this.minColumnsWidthInPercent);
    },

    setPercentsInPixel() {
      const { width: rowWidth } = this.tableRow.getBoundingClientRect();

      this.percentsInPixel = 100 / rowWidth;
    },

    calculateFullColumnsWidth(columnsWidth) {
      return Object.values(columnsWidth).reduce((acc, width) => acc + width, 0);
    },

    calculateColumnsWidths() {
      this.setPercentsInPixel();

      this.columnsWidthByField = [...this.headerCells].reduce((acc, headerElement) => {
        if (headerElement.dataset?.value) {
          const { value } = headerElement.dataset;
          const { width: headerWidth } = headerElement.getBoundingClientRect();
          const width = headerWidth * this.percentsInPixel;

          acc[value] = this.getNormalizedWidth(value, width);
        }

        return acc;
      }, {});
    },

    resizeColumnByDiff(index) {
      const diff = this.aggregatedMovementDiff;

      if (!diff) {
        return;
      }

      const resizingLeftColumn = this.headers[index].value;
      const previousLeftColumnWidth = this.getColumnWidthByField(resizingLeftColumn);
      const newLeftColumnWidth = this.getNormalizedWidth(resizingLeftColumn, previousLeftColumnWidth + diff);

      this.columnsWidthByField = {
        ...this.columnsWidthByField,
        [resizingLeftColumn]: newLeftColumnWidth,
      };

      if (newLeftColumnWidth !== previousLeftColumnWidth) {
        this.aggregatedMovementDiff = 0;
      }
    },

    handleColumnResize(event) {
      this.aggregatedMovementDiff += event.movementX * this.percentsInPixel;

      this.throttledResizeColumnByDiff(this.resizingColumnIndex);
    },

    finishColumnResize() {
      this.aggregatedMovementDiff = 0;

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
