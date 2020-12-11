import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore() {
      return request.get(API_ROUTES.engineRunInfo);
    },
  },
};
