import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore() {
      return request.get(API_ROUTES.engineRunInfo);
    },
  },
};
