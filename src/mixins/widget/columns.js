export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    columns() {
      if (this.widget.parameters.columnTranslations) {
        return this.widget.parameters.columnTranslations.map(v => ({ text: v.label, value: v.value }));
      }

      return [];
    },

    hasColumns() {
      return this.columns.length > 0;
    },
  },
};
