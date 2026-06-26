import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useServiceStore = defineStore('service', () => {
  const fetchEntityInfosKeysWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.entityInfosDictionaryKeys, { params })
  );
  const fetchServiceDependenciesWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.serviceDependencies, { params })
  );
  const fetchServiceImpactsWithoutStore = ({ params } = {}) => request.get(API_ROUTES.serviceImpacts, { params });
  const fetchEntityUpstreamWithoutStore = ({ params } = {}) => request.get(API_ROUTES.entityUpstream, { params });
  const fetchEntityDownstreamsWithoutStore = ({ params } = {}) => request.get(API_ROUTES.entityDownstreams, { params });
  const bulkRemoveServices = ({ data } = {}) => request.delete(API_ROUTES.bulkService, { data });

  return {
    fetchEntityInfosKeysWithoutStore,
    fetchServiceDependenciesWithoutStore,
    fetchServiceImpactsWithoutStore,
    fetchEntityUpstreamWithoutStore,
    fetchEntityDownstreamsWithoutStore,
    bulkRemoveServices,
  };
});
