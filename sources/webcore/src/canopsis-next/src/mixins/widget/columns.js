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
        return this.widget.parameters.widgetColumns.map(({ label, value, ...column }) => ({
          ...column,

          value,
          text: label,
        }));
      }

      return [];
    },

    hasColumns() {
      return this.columns.length > 0;
    },
  },
};
