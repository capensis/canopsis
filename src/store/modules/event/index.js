import qs from 'qs';

import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    async create(context, { data }) {
      try {
        await request.post(API_ROUTES.event, qs.stringify({ event: JSON.stringify(data) }), {
          headers: { 'content-type': 'application/x-www-form-urlencoded' },
        });
      } catch (e) {
        console.warn(e);
      }
    },
  },
};
