import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    removeAlarms(context, { id, data }) {
      return request.put(`${API_ROUTES.metaAlarm}/${id}/remove`, data);
    },
  },
};
