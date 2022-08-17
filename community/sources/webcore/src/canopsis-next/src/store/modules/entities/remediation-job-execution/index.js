import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    create(context, { data } = {}) {
      return request.post(API_ROUTES.remediation.jobExecutions, data);
    },

    cancel(context, { id }) {
      return request.put(`${API_ROUTES.remediation.jobExecutions}/${id}/cancel`);
    },
  },
};
