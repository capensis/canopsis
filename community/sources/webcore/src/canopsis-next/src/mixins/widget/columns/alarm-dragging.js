import Sortable from 'sortablejs';

export const widgetColumnDraggingAlarmMixin = {
  props: {
    columnsSettings: {
      type: Object,
      default: () => ({}),
    },
    draggableClass: {
      type: String,
      default: 'alarms-list-table__draggable-column',
    },
  },
  data() {
    return {
      draggingMode: false,
      columnsPositionByField: {},
    };
  },
  beforeDestroy() {
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

    handleColumnSort({ dragged, related }) {
      const oldDraggableIndex = this.headers.findIndex(({ value }) => value === dragged.dataset.value);
      const newDraggableIndex = this.headers.findIndex(({ value }) => value === related.dataset.value);
      const copiedHeaders = this.getHeadersByMovingIndexes(oldDraggableIndex, newDraggableIndex);

      this.calculateColumnsSortingOrderByHeaders(copiedHeaders);
    },

    startColumnDragging() {
      this.calculateColumnsSortingOrderByHeaders(this.headers);

      this.sortableInstance = Sortable.create(this.tableRow, {
        draggable: `.${this.draggableClass}`,
        onMove: this.handleColumnSort,
        direction: 'horizontal',
      });
    },

    finishColumnDragging() {
      this.sortableInstance?.destroy();
      this.sortableInstance = null;
    },
  },
};
