import Sortable from 'sortablejs';

const COLUMN_SWAP_REVERSE_THRESHOLD = 4;

export const widgetColumnDraggingAlarmMixin = {
  props: {
    draggableClass: {
      type: String,
      default: 'alarms-list-widget__draggable-column',
    },
  },
  data() {
    return {
      draggingMode: false,
      columnsPositionByField: {},
    };
  },
  beforeUnmount() {
    this.finishColumnDragging();
  },
  methods: {
    enableDraggingMode() {
      this.draggingMode = true;

      this.startColumnDragging();
    },

    disableDraggingMode() {
      this.draggingMode = false;

      this.finishColumnDragging();
    },

    toggleDraggingMode() {
      return this.draggingMode ? this.disableDraggingMode() : this.enableDraggingMode();
    },

    setColumnsPosition(columnsPosition) {
      this.columnsPositionByField = { ...columnsPosition };
    },

    calculateColumnsSortingOrderByHeaders(headers) {
      this.columnsPositionByField = headers.reduce((acc, { value }, index) => {
        acc[value] = index;

        return acc;
      }, {});
    },

    getColumnPositionByField(field) {
      return this.columnsPositionByField[field];
    },

    getHeadersByMovingIndexes(oldIndex, newIndex) {
      const copiedHeaders = [...this.headers];

      const [item] = copiedHeaders.splice(oldIndex, 1);

      copiedHeaders.splice(newIndex, 0, item);

      return copiedHeaders;
    },

    getColumnValueFromElement(element) {
      return (
        element?.dataset?.value ??
        element?.closest?.('th')?.dataset?.value ??
        element?.querySelector?.('[data-value]')?.dataset?.value ??
        element?.closest?.('th')?.querySelector?.('[data-value]')?.dataset?.value ??
        null
      );
    },

    getHeaderCellByValue(value) {
      return [...(this.tableRow?.querySelectorAll('th') ?? [])].find(
        cell => this.getColumnValueFromElement(cell) === value,
      );
    },

    isDraggableHeaderCell(cell) {
      return !!cell?.matches(`th:has(.table__dragging-handler), th.${this.draggableClass}`);
    },

    handlePointerColumnSwap(event) {
      if (this.isColumnSwapPending || !this.draggingField) {
        return;
      }

      const clientX = event?.clientX;

      if (clientX == null) {
        return;
      }

      const oldIndex = this.headers.findIndex(({ value }) => value === this.draggingField);

      if (oldIndex === -1) {
        return;
      }

      const draggedCell = this.getHeaderCellByValue(this.draggingField);

      if (!draggedCell) {
        return;
      }

      const { left, right } = draggedCell.getBoundingClientRect();
      let direction = 0;
      let targetIndex = oldIndex;

      if (clientX > right) {
        direction = 1;
        targetIndex = oldIndex + 1;
      } else if (clientX < left) {
        direction = -1;
        targetIndex = oldIndex - 1;
      } else {
        this.lastColumnSwapDirection = 0;
        this.lastSwapClientX = null;

        return;
      }

      if (
        this.lastColumnSwapDirection &&
        this.lastSwapClientX != null &&
        direction === -this.lastColumnSwapDirection
      ) {
        const movedOpposite = (this.lastSwapClientX - clientX) * this.lastColumnSwapDirection;

        if (movedOpposite < COLUMN_SWAP_REVERSE_THRESHOLD) {
          return;
        }
      }

      const targetHeader = this.headers[targetIndex];
      const targetCell = this.getHeaderCellByValue(targetHeader?.value);

      if (!this.isDraggableHeaderCell(targetCell)) {
        return;
      }

      this.isColumnSwapPending = true;
      this.lastColumnSwapDirection = direction;
      this.lastSwapClientX = clientX;
      this.calculateColumnsSortingOrderByHeaders(
        this.getHeadersByMovingIndexes(oldIndex, targetIndex),
      );
      this.$nextTick(() => {
        this.isColumnSwapPending = false;
      });
    },

    handleColumnSort(evt, originalEvent) {
      this.handlePointerColumnSwap(originalEvent);

      return false;
    },

    resetColumnSwapPointerState() {
      this.isColumnSwapPending = false;
      this.lastColumnSwapDirection = 0;
      this.lastSwapClientX = null;
    },

    bindPointerColumnSwap() {
      document.addEventListener('pointermove', this.handlePointerColumnSwap);
    },

    unbindPointerColumnSwap() {
      document.removeEventListener('pointermove', this.handlePointerColumnSwap);
      this.resetColumnSwapPointerState();
    },

    startColumnDragging() {
      this.calculateColumnsSortingOrderByHeaders(this.headers);

      if (!this.tableRow) {
        return;
      }

      this.sortableInstance = Sortable.create(this.tableRow, {
        draggable: `th:has(.table__dragging-handler), th.${this.draggableClass}`,
        handle: '.table__dragging-handler',
        filter: '.table__resize-handler',
        preventOnFilter: true,
        onStart: ({ item }) => {
          this.draggingField = this.getColumnValueFromElement(item);
          this.resetColumnSwapPointerState();
          this.bindPointerColumnSwap();
        },
        onMove: this.handleColumnSort,
        onEnd: () => {
          this.unbindPointerColumnSwap();
          this.draggingField = null;
        },
        direction: 'horizontal',
        forceFallback: true,
        fallbackOnBody: true,
        invertSwap: true,
        swapThreshold: 0.5,
      });
    },

    finishColumnDragging() {
      this.unbindPointerColumnSwap();
      this.sortableInstance?.destroy();
      this.sortableInstance = null;
      this.draggingField = null;
    },
  },
};
