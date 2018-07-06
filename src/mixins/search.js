import omit from 'lodash/omit';

/**
 * @mixin Add searching logic, need Pagination mixins
 * @see src/mixins/pagination.js
 */
export default {
  methods: {
    clear() {
      const query = omit(this.$route.query, [this.requestParam]);
      this.$router.push({ query });
    },
    submit() {
      const query = {
        ...this.$route.query,
      };
      query[this.requestParam] = this.requestData;
      this.$router.replace({ query });
    },
  },
};
