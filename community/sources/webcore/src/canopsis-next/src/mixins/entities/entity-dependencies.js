import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('service');

export const entitiesEntityDependenciesMixin = {
  props: {
    impact: {
      type: Boolean,
      default: false,
    },
  },
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
