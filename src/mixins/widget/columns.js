export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    columns() {
      if (this.widget.widget_columns) {
        return this.widget.widget_columns.map(v => ({ text: v.label, value: v.value }));
      }

      return [];
    },

    hasColumns() {
      return this.columns.length > 0;
    },
  },
};
