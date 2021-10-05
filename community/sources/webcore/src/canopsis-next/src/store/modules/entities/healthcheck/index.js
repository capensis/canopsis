import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchEnginesWithoutStore() {
      return request.get(API_ROUTES.healthcheck.engines);
    },

    fetchStatusWithoutStore() {
      return request.get(API_ROUTES.healthcheck.status);
    },
  },
};
