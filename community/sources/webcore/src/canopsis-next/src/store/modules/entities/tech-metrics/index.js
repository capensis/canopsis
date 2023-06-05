import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    async createTechMetricsExport() {
      return request.post(API_ROUTES.techMetrics);
    },

    async fetchTechMetricsExport() {
      return request.get(API_ROUTES.techMetrics);
    },
  },
};
