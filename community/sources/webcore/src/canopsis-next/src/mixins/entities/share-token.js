import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('shareToken');

export const entitiesShareTokenMixin = {
  computed: {
    ...mapGetters({
      shareTokens: 'items',
      shareTokensPending: 'pending',
      shareTokensMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchShareTokensList: 'fetchList',
      fetchShareTokensListWithPreviousParams: 'fetchListWithPreviousParams',
      createShareToken: 'create',
      updateShareToken: 'update',
      removeShareToken: 'remove',
    }),
  },
};
