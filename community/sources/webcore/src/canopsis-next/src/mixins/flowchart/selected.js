export const selectedShapesMixin = {
  props: {
    selected: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      selectedShapes: [...this.selected],
    };
  },
  computed: {
    hasSelected() {
      return !!this.selectedShapes.length;
    },
  },
  methods: {
    isSelected(id) {
      return this.selectedShapes.includes(id);
    },

    updateSelected() {
      this.$emit('update:selected', [...this.selectedShapes]);
    },

    setSelected(shape) {
      if (!this.isSelected(shape._id)) {
        this.selectedShapes.push(shape._id);

        this.updateSelected();
      }
    },

    clearSelected() {
      this.selectedShapes = [];
      this.updateSelected();
    },

    removeSelectedShape(shape) {
      this.selectedShapes = this.selectedShapes.filter(id => id !== shape._id);

      this.updateSelected();
    },
  },
};
