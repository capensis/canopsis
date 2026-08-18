import { throttle, isEmpty } from 'lodash';

export const widgetColumnResizingAlarmMixin = {
  props: {
    resizingColumnThrottleDelay: {
      type: Number,
      default: 10,
    },
    minColumnWidth: {
      type: Number,
      default: 40,
    },
  },
  data() {
    return {
      resizingMode: false,
      resizingColumnIndex: null,
      columnsWidthByField: {},
    };
  },
  created() {
    this.aggregatedMovementDiff = 0;
    this.throttledResizeColumnByDiff = throttle(
      this.resizeColumnByDiff,
      this.resizingColumnThrottleDelay,
    );
  },
  beforeUnmount() {
    this.finishColumnResize({ suppressClick: false });
  },
  computed: {
    tableHeader() {
      return (
        this.$refs.dataTable?.$el?.querySelector('thead') ??
        this.$el?.querySelector('.alarms-list-widget__table thead')
      );
    },

    tableRow() {
      return this.tableHeader?.querySelector('tr:first-of-type');
    },

    headerCells() {
      return this.tableRow?.querySelectorAll('th') ?? [];
    },

    sumOfColumnsWidth() {
      return this.calculateFullColumnsWidth(this.columnsWidthByField);
    },

    isEmptyColumnsWidthByField() {
      return isEmpty(this.columnsWidthByField);
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
      return Math.max(newWidth, this.minColumnWidth);
    },

    calculateFullColumnsWidth(columnsWidth) {
      return Object.values(columnsWidth).reduce((acc, width) => acc + width, 0);
    },

    calculateElementNormalizedWidth(element, field) {
      const { width } = element.getBoundingClientRect();

      return this.getNormalizedWidth(field, width);
    },

    calculateColumnsWidths() {
      this.columnsWidthByField = [...this.headerCells].reduce((acc, headerElement) => {
        if (headerElement.dataset?.value) {
          const { value } = headerElement.dataset;

          acc[value] = this.calculateElementNormalizedWidth(headerElement, value);
        }

        return acc;
      }, {});
    },

    resizeColumnByDiff(index) {
      const diff = this.aggregatedMovementDiff;

      if (!diff) {
        return;
      }

      const resizingLeftColumn = this.headers[index]?.value;
      const previousLeftColumnWidth = this.getColumnWidthByField(resizingLeftColumn);
      const newLeftColumnWidth = this.getNormalizedWidth(
        resizingLeftColumn,
        previousLeftColumnWidth + diff,
      );

      this.columnsWidthByField = {
        ...this.columnsWidthByField,
        [resizingLeftColumn]: newLeftColumnWidth,
      };

      if (newLeftColumnWidth !== previousLeftColumnWidth) {
        this.aggregatedMovementDiff = 0;
      }
    },

    handleColumnResize(event) {
      this.aggregatedMovementDiff += event.movementX;

      this.throttledResizeColumnByDiff(this.resizingColumnIndex);
    },

    suppressNextClickAfterResize() {
      const suppressNextClick = (clickEvent) => {
        clickEvent.preventDefault();
        clickEvent.stopPropagation();
        clickEvent.stopImmediatePropagation();
        document.removeEventListener('click', suppressNextClick, true);
      };

      document.addEventListener('click', suppressNextClick, true);
      window.setTimeout(() => {
        document.removeEventListener('click', suppressNextClick, true);
      }, 0);
    },

    startColumnResize(columnName) {
      this.finishColumnResize({ suppressClick: false });

      const body = document.querySelector('body');

      this.resizingColumnIndex = this.headers.findIndex(({ value }) => value === columnName);

      if (!body) {
        return;
      }

      this.sortableInstance?.option('disabled', true);

      body.addEventListener('pointermove', this.handleColumnResize);
      body.addEventListener('pointerup', this.finishColumnResize);
      body.addEventListener('pointercancel', this.finishColumnResize);
    },

    finishColumnResize(event) {
      this.aggregatedMovementDiff = 0;

      const body = document.querySelector('body');

      if (body) {
        body.removeEventListener('pointermove', this.handleColumnResize);
        body.removeEventListener('pointerup', this.finishColumnResize);
        body.removeEventListener('pointercancel', this.finishColumnResize);
      }

      if (this.draggingMode) {
        this.sortableInstance?.option('disabled', false);
      }

      if (event?.suppressClick === false) {
        return;
      }

      this.suppressNextClickAfterResize();
    },
  },
};
