import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('pbehaviorDatesExceptions');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      pbehaviorDatesExceptions: 'items',
      pbehaviorDatesExceptionsPending: 'pending',
      pbehaviorDatesExceptionsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchPbehaviorDatesExceptionsList: 'fetchList',
      fetchPbehaviorDatesExceptionsListWithoutStore: 'fetchListWithoutStore',
      createPbehaviorDateException: 'create',
      updatePbehaviorDateException: 'update',
      removePbehaviorDateException: 'remove',
    }),
  },
};
