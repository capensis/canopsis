import qs from 'qs';

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
      return request.get(API_ROUTES.snmpMib, { params });
    },

    upload(context, { data } = {}) {
      return request.post(API_ROUTES.snmpMib, qs.stringify({
        filecontent: JSON.stringify([{
          filename: 'concatenatedMibFiles',
          data,
        }]),
      }), {
        headers: { 'content-type': 'application/x-www-form-urlencoded' },
      });
    },
  },
};
