import request from '@/services/request';

import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore(context, { data }) {
      return request.post(API_ROUTES.pbehavior.timespan, data);
    },
  },
};
