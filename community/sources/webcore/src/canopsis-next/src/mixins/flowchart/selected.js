import { isAreaIncludeShape } from '@/helpers/flowchart/points';

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

      selecting: false,
      selectionStart: {
        x: 0,
        y: 0,
      },
    };
  },
  computed: {
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
      this.setSelected([]);
      this.updateSelected();
    },

    removeSelectedShape(shape) {
      this.selectedIds = this.selectedIds.filter(id => id !== shape._id);

      this.updateSelected();
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
        if (!isAreaIncludeShape(normalizedStart, normalizedEnd, shape)) {
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
