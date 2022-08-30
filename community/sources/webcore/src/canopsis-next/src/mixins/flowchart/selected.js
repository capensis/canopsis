import { SHAPES } from '@/constants';

export const selectedShapesMixin = {
  props: {
    selected: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      selectedIds: [],

      selection: false,
      selectionStart: {
        x: 0,
        y: 0,
      },
    };
  },
  computed: {
    selectionPath() {
      return [
        `M ${this.selectionStart.x} ${this.selectionStart.y}`,
        `L ${this.selectionStart.x} ${this.cursor.y}`,
        `L ${this.cursor.x} ${this.cursor.y}`,
        `L ${this.cursor.x} ${this.selectionStart.y}`,
        'Z',
      ];
    },

    hasSelected() {
      return !!this.selectedIds.length;
    },
  },
  watch: {
    selected: {
      immediate: true,
      handler(selected) {
        this.selectedIds = [...selected];
      },
    },
  },
  methods: {
    isSelected(id) {
      return this.selectedIds.includes(id);
    },

    updateSelected() {
      this.$emit('update:selected', [...this.selectedIds]);
    },

    setSelected(selected) {
      this.selectedIds = selected;
    },

    setSelectedShape(shape) {
      if (!this.isSelected(shape._id)) {
        this.selectedIds.push(shape._id);

        this.updateSelected();
      }
    },

    clearSelected() {
      this.selectedIds = [];
      this.updateSelected();
    },

    removeSelectedShape(shape) {
      this.selectedIds = this.selectedIds.filter(id => id !== shape._id);

      this.updateSelected();
    },

    isAreaIncludeShape(start, end, shape) {
      switch (shape.type) {
        case SHAPES.storage:
        case SHAPES.parallelogram:
        case SHAPES.image:
        case SHAPES.rhombus:
        case SHAPES.ellipse:
        case SHAPES.process:
        case SHAPES.document:
        case SHAPES.rect:
          return shape.x > start.x
            && shape.x + shape.width < end.x
            && shape.y > start.y
            && shape.y + shape.height < end.y;
        case SHAPES.circle:
          return shape.x > start.x
            && shape.x + shape.diameter < end.x
            && shape.y > start.y
            && shape.y + shape.diameter < end.y;
        case SHAPES.arrowLine:
        case SHAPES.bidirectionalArrowLine:
        case SHAPES.line:
          return shape.points.every(
            point => point.x > start.x
              && point.x < end.x
              && point.y > start.y
              && point.y < end.y,
          );
      }

      return false;
    },

    selectShapesByArea(start, end, shiftKey) {
      const normalizedStart = {
        x: Math.min(start.x, end.x),
        y: Math.min(start.y, end.y),
      };
      const normalizedEnd = {
        x: Math.max(start.x, end.x),
        y: Math.max(start.y, end.y),
      };

      this.selectedIds = Object.values(this.data).reduce((acc, shape) => {
        /**
         * Selection area not include a shape
         */
        if (!this.isAreaIncludeShape(normalizedStart, normalizedEnd, shape)) {
          return acc;
        }

        /**
         * Selection area include already selected shape and shift pressed.
         */
        if (shiftKey && acc.includes(shape._id)) {
          return acc.filter(id => id !== shape._id);
        }

        /**
         * Selection area include shape
         */
        acc.push(shape._id);

        return acc;
      }, shiftKey ? [...this.selectedIds] : []);

      this.updateSelected();
    },
  },
};
