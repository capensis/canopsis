import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchItemWithoutStore() {
      return request.get(API_ROUTES.notification);
    },
    update(context, { data }) {
      return request.put(API_ROUTES.notification, data);
    },
  },
};
