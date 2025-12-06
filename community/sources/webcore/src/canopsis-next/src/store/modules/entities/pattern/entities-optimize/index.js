import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    optimize(context, { data }) {
      return request.post(API_ROUTES.pattern.entitiesOptimize, data);
    },

    fetchOptimizeStatus(context, { id }) {
      return request.get(`${API_ROUTES.pattern.entitiesOptimize}/${id}`);
    },

    update(context, { id, data }) {
      return request.put(`${API_ROUTES.pattern.entitiesOptimize}/${id}`, data);
    },

    remove(context, { id }) {
      return request.delete(`${API_ROUTES.pattern.entitiesOptimize}/${id}`);
    },
  },
};
