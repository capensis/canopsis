import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore(context, { data }) {
      return request.post(API_ROUTES.pbehavior.timespan, data);
    },
  },
};
