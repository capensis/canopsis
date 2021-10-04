import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    fetchItemWithoutStore() {
      return request.get(API_ROUTES.healthcheck.parameters);
    },

    update(context, { data }) {
      return request.put(API_ROUTES.healthcheck.parameters, data);
    },
  },
};
