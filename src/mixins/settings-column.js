export default {
  props: {
    prefixFormatter: {
      type: Function,
      default: value => value,
    },
    value: {
      type: [Array, Object],
      default: () => [],
    },
  },
  data() {
    return {
      columns: this.value,
    };
  },
  methods: {
    updateValue(index, key, value) {
      const items = [...this.columns];

      items[index] = {
        ...items[index],
        [key]: this.prefixFormatter(value),
      };

      this.$emit('input', items);
    },
  },
};
