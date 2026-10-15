import { mapActions } from 'pinia';

import { useServiceStore } from '@/stores/modules/service';

export const entitiesServiceMixin = {
  methods: {
    ...mapActions(useServiceStore, {
      fetchServiceAlarmsWithoutStore: 'fetchAlarmsWithoutStore',
      fetchServiceItemWithoutStore: 'fetchItemWithoutStore',
      createService: 'create',
      editService: 'update',
      removeService: 'remove',
      bulkRemoveServices: 'bulkRemoveServices',
      fetchServiceDependenciesWithoutStore: 'fetchServiceDependenciesWithoutStore',
      fetchServiceImpactsWithoutStore: 'fetchServiceImpactsWithoutStore',
    }),
  },
};
