import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('associativeTable');

export const entitiesAssociativeTableMixin = {
  methods: {
    ...mapActions({
      fetchAssociativeTable: 'fetch',
      createAssociativeTable: 'create',
      updateAssociativeTable: 'update',
      removeAssociativeTable: 'remove',
    }),
  },
};
