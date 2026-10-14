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
  const fetchAlarmsWithoutStore = ({ id, params = {} } = {}) =>
    request.get(`${API_ROUTES.serviceAlarms}/${encodeURIComponent(id)}`, { params });
  const fetchEntityUpstreamWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.entityUpstream, { params });
  const fetchEntityDownstreamsWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.entityDownstreams, { params });
  const bulkRemoveServices = ({ data } = {}) => request.delete(API_ROUTES.bulkService, { data });

  const create = ({ data } = {}) => request.post(API_ROUTES.service, data);
  const update = ({ id, data } = {}) =>
    request.put(`${API_ROUTES.service}/${encodeURIComponent(id)}`, data);
  const remove = ({ id } = {}) => request.delete(`${API_ROUTES.service}/${encodeURIComponent(id)}`);
  const fetchItemWithoutStore = ({ id, params } = {}) =>
    request.get(`${API_ROUTES.service}/${encodeURIComponent(id)}`, { params });

  return {
    fetchEntityInfosKeysWithoutStore,
    fetchServiceDependenciesWithoutStore,
    fetchServiceImpactsWithoutStore,
    fetchAlarmsWithoutStore,
    fetchEntityUpstreamWithoutStore,
    fetchEntityDownstreamsWithoutStore,
    bulkRemoveServices,
    create,
    update,
    remove,
    fetchItemWithoutStore,
  };
});
