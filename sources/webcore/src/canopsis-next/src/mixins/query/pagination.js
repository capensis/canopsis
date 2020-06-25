export default {
  methods: {
    updateQueryPage(page) {
      this.$emit('update:query', { ...this.query, page });
    },

    updateQueryLimit(limit) {
      this.$emit('update:query', { ...this.query, limit });
    },
  },
};
