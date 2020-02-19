import { DEFAULT_WIDGET_GROUP_COLUMNS } from '@/constants';

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

    groupColumns() {
      if (this.widget.parameters.widgetGroupColumns) {
        return this.widget.parameters.widgetGroupColumns.map(this.mapColumnEntity);
      }

      return DEFAULT_WIDGET_GROUP_COLUMNS.map(({ label, value }) => ({
        text: this.$t(label),
        value,
      }));
    },

    hasColumns() {
      return this.columns.length > 0;
    },

    hasGroupColumns() {
      return this.columns.length > 0;
    },
  },
  methods: {
    mapColumnEntity({ label, value, ...column }) {
      return ({
        ...column,
        value,
        text: label,
      });
    },
  },
};
