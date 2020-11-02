import request from '@/services/request';

import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    fetchItems(context, { data }) {
      return request.post(API_ROUTES.planning.timespan, data);
    },
  },
};
