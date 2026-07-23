import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const usePbehaviorEntitiesStore = defineStore('pbehaviorEntities', () => {
  const fetchListWithoutStore = ({ id, params } = {}) =>
    request.get(`${API_ROUTES.pbehavior.pbehaviors}/${id}/entities`, { params });

  return {
    fetchListWithoutStore,
  };
});
