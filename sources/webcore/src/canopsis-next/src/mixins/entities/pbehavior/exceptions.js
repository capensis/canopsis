import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('pbehaviorExceptions');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      pbehaviorExceptions: 'items',
      pbehaviorExceptionsPending: 'pending',
      pbehaviorExceptionsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchPbehaviorExceptionsList: 'fetchList',
      fetchPbehaviorExceptionsListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchPbehaviorExceptionsListWithoutStore: 'fetchListWithoutStore',
      createPbehaviorException: 'create',
      updatePbehaviorException: 'update',
      removePbehaviorException: 'remove',
    }),
  },
};
