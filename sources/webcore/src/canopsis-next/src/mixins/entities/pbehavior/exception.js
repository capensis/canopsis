import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('pbehaviorException');

/**
 * @mixin
 */
export default {
  data() {
    return {
      pbehaviorExceptionsPending: false,
      pbehaviorExceptions: [],
      pbehaviorExceptionsMeta: {},
    };
  },
  methods: {
    ...mapActions({
      fetchPbehaviorExceptionsListWithoutStore: 'fetchListWithoutStore',
      removeException: 'remove',
    }),

    async fetchPbehaviorExceptionsList() {
      this.pbehaviorExceptionsPending = true;

      const result = await this.fetchPbehaviorExceptionsListWithoutStore({ params: this.getQuery() });

      if (result) {
        // this.pbehaviorExceptions = result.data;
        this.pbehaviorExceptionsMeta = result.meta;
      }

      this.pbehaviorExceptionsPending = false;
    },
  },
};
