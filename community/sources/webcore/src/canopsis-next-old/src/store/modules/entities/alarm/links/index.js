import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchItemWithoutStore(context, { params, id }) {
      return request.get(`${API_ROUTES.alarmLinks}/${id}`, { params });
    },
  },
};
