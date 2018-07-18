import omit from 'lodash/omit';

/**
 * @mixin Add searching logic, need Pagination mixins
 * @see src/mixins/pagination.js
 */
export default {
  props: {
    query: {
      type: Object,
      required: true,
    },
  },
  methods: {
    clear() {
      this.searchingText = '';

      this.$emit('update:query', omit(this.query, [this.requestParam]));
    },
    submit() {
      this.$emit('update:query', {
        ...this.query,
        [this.requestParam]: this.requestData,
      });
    },
  },
};
