import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  methods: {
    fetchListWithoutStore(context, { params, id }) {
      return request.get(`${API_ROUTES.alarmLinks}/${id}`, { params });
    },
  },
};
