import { omit } from 'lodash';

/**
 * @mixin Add searching logic, need Pagination mixins
 * @see src/mixins/query.js
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

        page: 1,
        [this.requestParam]: this.requestData,
      });
    },
  },
};
