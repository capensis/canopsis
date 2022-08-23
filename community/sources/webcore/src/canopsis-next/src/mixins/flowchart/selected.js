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
      this.selectedIds = [];
      this.updateSelected();
    },

    removeSelectedShape(shape) {
      this.selectedIds = this.selectedIds.filter(id => id !== shape._id);

      this.updateSelected();
    },
  },
};
