import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  actions: {
    fetchList(context, { params = {} } = {}) {
      return request.post(API_ROUTES.snmpMib.list, params);
    },

    fetchDistinctList(context, { params = {} } = {}) {
      return request.post(API_ROUTES.snmpMib.distinct, params);
    },
  },
};
