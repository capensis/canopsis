export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    columns() {
      if (this.widget.parameters.widgetColumns) {
        return this.widget.parameters.widgetColumns.map(v => ({ text: v.label, value: v.value }));
      }

      return [];
    },

    hasColumns() {
      return this.columns.length > 0;
    },
  },
};
