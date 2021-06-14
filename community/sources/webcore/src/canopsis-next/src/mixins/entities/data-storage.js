import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('dataStorage');

export const entitiesDataStorageSettingsMixin = {
  methods: {
    ...mapActions({
      fetchDataStorageSettingsWithoutStore: 'fetchItemWithoutStore',
      updateDataStorageSettings: 'update',
    }),
  },
};
