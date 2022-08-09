import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('pbehaviorTypes');

export const entitiesPbehaviorTypeMixin = {
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

    async fetchDefaultPbehaviorTypes() {
      const { data } = await this.fetchPbehaviorTypesListWithoutStore({
        params: { default: true },
      });

      return data;
    },
  },
};
