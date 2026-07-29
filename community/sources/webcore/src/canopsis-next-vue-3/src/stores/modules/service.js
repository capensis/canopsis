import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useServiceStore = defineStore('service', () => {
  const fetchEntityInfosKeysWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.entityInfosDictionaryKeys, { params });
  const fetchServiceDependenciesWithoutStore = ({ id, params = {} } = {}) =>
    request.get(API_ROUTES.serviceDependencies, {
      params: { ...params, _id: id },
    });
  const fetchServiceImpactsWithoutStore = ({ id, params = {} } = {}) =>
    request.get(API_ROUTES.serviceImpacts, {
      params: { ...params, _id: id },
    });
  const fetchEntityUpstreamWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.entityUpstream, { params });
  const fetchEntityDownstreamsWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.entityDownstreams, { params });
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
