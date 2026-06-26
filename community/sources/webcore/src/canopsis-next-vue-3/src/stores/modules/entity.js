import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useEntityStore = defineStore('entity', () => {
  const fetchContextEntitiesListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.entity, { params });
  const fetchEntityInfosLogsListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.entityInfosLog, { params });
  const archiveDisabledEntitiesData = ({ data } = {}) => request.post(`${API_ROUTES.entity}/archive-disabled`, data);
  const archiveUnlinkedEntitiesData = ({ data } = {}) => request.post(`${API_ROUTES.entity}/archive-unlinked`, data);
  const cleanArchivedEntitiesData = () => request.post(`${API_ROUTES.entity}/clean-archived`);
  const checkStateSetting = ({ data } = {}) => request.post(API_ROUTES.entityCheckStateSetting, data);

  return {
    fetchContextEntitiesListWithoutStore,
    fetchEntityInfosLogsListWithoutStore,
    archiveDisabledEntitiesData,
    archiveUnlinkedEntitiesData,
    cleanArchivedEntitiesData,
    checkStateSetting,
  };
});
