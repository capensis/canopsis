import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    create(context, { data } = {}) {
      return request.post(API_ROUTES.entityComments, data);
    },

    update(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.entityComments}/${id}`, data);
    },

    fetchListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.entityComments, { params });
    },
  },
};
