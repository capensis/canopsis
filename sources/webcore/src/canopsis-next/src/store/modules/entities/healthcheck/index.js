import request from '@/services/request';

import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    fetchWithoutStore() {
      return request.get(API_ROUTES.healthcheck);
    },
  },
};
