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
        return Object.keys(this.widget.parameters.columnTranslations).map(
          v => ({ text: this.widget.parameters.columnTranslations[v], value: v }),
        );
      }

      return [];
    },

    hasColumns() {
      return this.columns.length > 0;
    },
  },
};
