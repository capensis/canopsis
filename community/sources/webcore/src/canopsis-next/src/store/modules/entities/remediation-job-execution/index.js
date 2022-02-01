import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    async create(context, { data } = {}) {
      return request.post(API_ROUTES.remediation.jobExecutions, data);
    },
  },
};
