import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('healthcheckParameters');

export const entitiesHealthcheckParametersMixin = {
  methods: {
    ...mapActions({
      fetchHealthcheckParametersWithoutStore: 'fetchItemWithoutStore',
      updateHealthcheckParameters: 'update',
    }),
  },
};
