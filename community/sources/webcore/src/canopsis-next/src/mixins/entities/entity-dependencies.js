import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('service');

export const entitiesEntityDependenciesMixin = {
  methods: {
    ...mapActions({
      fetchServiceDependenciesWithoutStore: 'fetchDependenciesWithoutStore',
      fetchServiceImpactsWithoutStore: 'fetchImpactsWithoutStore',
    }),

    fetchDependenciesList(data) {
      return this.impact
        ? this.fetchServiceImpactsWithoutStore(data)
        : this.fetchServiceDependenciesWithoutStore(data);
    },
  },
};
