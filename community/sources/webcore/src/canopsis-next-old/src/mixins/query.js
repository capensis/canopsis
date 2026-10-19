import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('query');

export const queryMixin = {
  computed: {
    ...mapGetters(['getQueryById', 'getQueryNonceById']),
  },
  methods: {
    ...mapActions({
      updateQuery: 'update',
      mergeQuery: 'merge',
      removeQuery: 'remove',

      updateLockedQuery: 'updateLocked',
      removeLockedQuery: 'removeLocked',

      forceUpdateQuery: 'forceUpdate',
    }),
  },
};
