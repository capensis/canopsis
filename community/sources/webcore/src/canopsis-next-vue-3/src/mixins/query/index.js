import { mapActions, mapStores } from 'pinia';

import { useQueryStore } from '@/stores/modules/query';

export const queryMixin = {
  computed: {
    ...mapStores(useQueryStore),

    getQueryById() {
      return this.queryStore.getQueryById;
    },

    getQueryNonceById() {
      return this.queryStore.getQueryNonceById;
    },
  },

  methods: {
    ...mapActions(useQueryStore, {
      updateQuery: 'update',
      mergeQuery: 'merge',
      removeQuery: 'remove',

      updateLockedQuery: 'updateLocked',
      removeLockedQuery: 'removeLocked',

      forceUpdateQuery: 'forceUpdate',
    }),
  },
};
