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
        return this.widget.parameters.widgetColumns.map(this.mapColumnEntity);
      }

      return [];
    },

    hasColumns() {
      return this.columns.length > 0;
    },
  },
  methods: {
    mapColumnEntity({ label, value, ...column }) {
      return {
        ...column,
        value,
        text: label,
      };
    },
  },
};
