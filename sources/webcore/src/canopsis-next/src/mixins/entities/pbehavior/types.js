import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('pbehaviorTypes');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      pbehaviorTypes: 'items',
      pbehaviorTypesPending: 'pending',
      pbehaviorTypesMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchPbehaviorTypesList: 'fetchList',
      fetchPbehaviorTypesListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchPbehaviorTypesListWithoutStore: 'fetchListWithoutStore',
      createPbehaviorType: 'create',
      updatePbehaviorType: 'update',
      removePbehaviorType: 'remove',
      fetchPbehaviorTypeByEntityId: 'fetchListByEntityId',
    }),
  },
};
