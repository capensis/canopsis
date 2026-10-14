import { mapActions } from 'pinia';

import { useAssociativeTableStore } from '@/stores/modules/associative-table';

export const entitiesAssociativeTableMixin = {
  methods: {
    ...mapActions(useAssociativeTableStore, {
      fetchAssociativeTable: 'fetch',
      createAssociativeTable: 'create',
      updateAssociativeTable: 'update',
      removeAssociativeTable: 'remove',
    }),
  },
};
