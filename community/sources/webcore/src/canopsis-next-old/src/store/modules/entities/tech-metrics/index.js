import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    createTechMetricsExport() {
      return request.post(API_ROUTES.techMetrics);
    },

    fetchTechMetricsExport() {
      return request.get(API_ROUTES.techMetrics);
    },

    fetchTechMetricsSettings() {
      return request.get(API_ROUTES.techMetricsSettings);
    },

    updateTechMetricsSettings(context, { data }) {
      return request.put(API_ROUTES.techMetricsSettings, data);
    },
  },
};
