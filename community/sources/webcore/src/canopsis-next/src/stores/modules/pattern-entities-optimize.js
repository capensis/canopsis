import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const usePatternEntitiesOptimizeStore = defineStore('pattern-entities-optimize', () => {
  const optimize = ({ data, cancelToken } = {}) =>
    request.post(API_ROUTES.pattern.entitiesOptimize, data, { cancelToken });

  const fetchOptimizeStatus = ({ id, cancelToken } = {}) =>
    request.get(`${API_ROUTES.pattern.entitiesOptimize}/${id}`, { cancelToken });

  const update = ({ id, data } = {}) =>
    request.put(`${API_ROUTES.pattern.entitiesOptimize}/${id}`, data);

  const remove = ({ id } = {}) => request.delete(`${API_ROUTES.pattern.entitiesOptimize}/${id}`);

  return {
    optimize,
    fetchOptimizeStatus,
    update,
    remove,
  };
});
