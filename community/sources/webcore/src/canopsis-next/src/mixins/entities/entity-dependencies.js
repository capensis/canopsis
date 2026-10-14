import { mapActions } from 'pinia';

import { useServiceStore } from '@/stores/modules/service';

export const entitiesEntityDependenciesMixin = {
  methods: {
    ...mapActions(useServiceStore, [
      'fetchServiceDependenciesWithoutStore',
      'fetchServiceImpactsWithoutStore',
    ]),

    fetchDependenciesList(data) {
      return this.impact
        ? this.fetchServiceImpactsWithoutStore(data)
        : this.fetchServiceDependenciesWithoutStore(data);
    },
  },
};
