import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.metaAlarm, { params });
    },

    create(context, { data }) {
      return request.post(API_ROUTES.metaAlarm, data);
    },

    addAlarms(context, { id, data }) {
      return request.put(`${API_ROUTES.metaAlarm}/${id}/add`, data);
    },

    removeAlarms(context, { id, data }) {
      return request.put(`${API_ROUTES.metaAlarm}/${id}/remove`, data);
    },
  },
};
