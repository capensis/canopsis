import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchStatusWithoutStore() {
      return request.get(API_ROUTES.healthcheck);
    },
  },
};
