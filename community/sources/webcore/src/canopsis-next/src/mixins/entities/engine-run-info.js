import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('engineRunInfo');

export const entitiesEngineRunInfoMixin = {
  methods: {
    ...mapActions({
      fetchEnginesListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
