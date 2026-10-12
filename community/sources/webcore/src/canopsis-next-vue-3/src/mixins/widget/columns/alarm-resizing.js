import { isEmpty } from 'lodash';

const TABLE_WIDTH_CSS_VARIABLE = '--alarms-list-table-width';
const COLUMN_WIDTH_CSS_VARIABLE_PREFIX = '--alarms-list-table-column-width';

export const widgetColumnResizingAlarmMixin = {
  props: {
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
    this.resizeColumnAnimationFrame = null;
    this.resizingColumnWidth = null;
    this.resizingTableWidth = null;
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

    isColumnResizing() {
      return this.resizingColumnIndex !== null;
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
      const { offsetWidth: width } = element;

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

    getColumnWidthCssVariable(index) {
      return `${COLUMN_WIDTH_CSS_VARIABLE_PREFIX}-${index}`;
    },

    applyColumnResizePreview(index, columnWidth, tableWidth) {
      const tableElement = this.$refs.dataTable?.$el;

      if (!tableElement) {
        return;
      }

      tableElement.style.setProperty(this.getColumnWidthCssVariable(index), `${columnWidth}px`);
      tableElement.style.setProperty(TABLE_WIDTH_CSS_VARIABLE, `${tableWidth}px`);
    },

    clearColumnResizePreview(index) {
      const tableElement = this.$refs.dataTable?.$el;

      if (!tableElement || index === null) {
        return;
      }

      tableElement.style.removeProperty(this.getColumnWidthCssVariable(index));

      if (this.sumOfColumnsWidth) {
        tableElement.style.setProperty(TABLE_WIDTH_CSS_VARIABLE, `${this.sumOfColumnsWidth}px`);
      }
    },

    resizeColumnByDiff(index) {
      const diff = this.aggregatedMovementDiff;

      if (!diff) {
        return;
      }

      const previousLeftColumnWidth = this.resizingColumnWidth;

      if (!Number.isFinite(previousLeftColumnWidth) || !Number.isFinite(this.resizingTableWidth)) {
        return;
      }

      const newLeftColumnWidth = this.getNormalizedWidth(
        this.headers[index]?.value,
        previousLeftColumnWidth + diff,
      );
      const appliedDiff = newLeftColumnWidth - previousLeftColumnWidth;

      if (appliedDiff) {
        this.resizingColumnWidth = newLeftColumnWidth;
        this.resizingTableWidth += appliedDiff;
        this.aggregatedMovementDiff = 0;

        this.applyColumnResizePreview(index, newLeftColumnWidth, this.resizingTableWidth);
      }
    },

    handleColumnResize(event) {
      this.aggregatedMovementDiff += event.movementX;

      if (this.resizeColumnAnimationFrame !== null) {
        return;
      }

      this.resizeColumnAnimationFrame = window.requestAnimationFrame(() => {
        this.resizeColumnAnimationFrame = null;
        this.resizeColumnByDiff(this.resizingColumnIndex);
      });
    },

    flushColumnResize() {
      if (this.resizeColumnAnimationFrame !== null) {
        window.cancelAnimationFrame(this.resizeColumnAnimationFrame);
        this.resizeColumnAnimationFrame = null;
      }

      this.resizeColumnByDiff(this.resizingColumnIndex);
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
      this.resizingColumnWidth = this.getColumnWidthByField(columnName);
      this.resizingTableWidth = this.sumOfColumnsWidth;

      if (!body) {
        return;
      }

      this.sortableInstance?.option('disabled', true);

      body.addEventListener('pointermove', this.handleColumnResize);
      body.addEventListener('pointerup', this.finishColumnResize);
      body.addEventListener('pointercancel', this.finishColumnResize);
    },

    finishColumnResize(event) {
      const resizingColumnIndex = this.resizingColumnIndex;

      this.flushColumnResize();

      const resizingColumn = this.headers[resizingColumnIndex]?.value;

      if (resizingColumn && Number.isFinite(this.resizingColumnWidth)) {
        this.columnsWidthByField = {
          ...this.columnsWidthByField,
          [resizingColumn]: this.resizingColumnWidth,
        };
      }

      this.clearColumnResizePreview(resizingColumnIndex);

      this.aggregatedMovementDiff = 0;
      this.resizingColumnIndex = null;
      this.resizingColumnWidth = null;
      this.resizingTableWidth = null;

      if (resizingColumnIndex !== null) {
        this.refreshHorizontalScrollbarAfterColumnResize?.();
      }

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
