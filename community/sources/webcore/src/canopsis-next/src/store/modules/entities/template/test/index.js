import { API_ROUTES } from '@/config';

import request from '@/services/request';

// TODO: move the similar logic into special store creator
export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.templateTests, { params });
    },

    create(context, { data }) {
      return request.post(API_ROUTES.templateTests, data);
    },

    update(context, { id, data }) {
      return request.put(`${API_ROUTES.templateTests}/${id}`, data);
    },

    remove(context, { id } = {}) {
      return request.delete(`${API_ROUTES.templateTests}/${id}`);
    },
  },
};
