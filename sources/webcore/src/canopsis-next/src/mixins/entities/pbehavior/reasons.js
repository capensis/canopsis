import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('pbehaviorReasons');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      pbehaviorReasons: 'items',
      pbehaviorReasonsPending: 'pending',
      pbehaviorReasonsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchPbehaviorReasonsList: 'fetchList',
      fetchPbehaviorReasonsListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchPbehaviorReasonsListWithoutStore: 'fetchListWithoutStore',
      createPbehaviorReason: 'create',
      updatePbehaviorReason: 'update',
      removePbehaviorReason: 'remove',
      fetchPbehaviorReasonByEntityId: 'fetchListByEntityId',
    }),
  },
};
