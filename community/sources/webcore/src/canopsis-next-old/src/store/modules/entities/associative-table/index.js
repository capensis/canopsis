import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    async fetch(context, { name } = {}) {
      const { content } = await request.get(API_ROUTES.associativeTable, { params: { name } });

      return content;
    },

    async create(context, { name, data } = {}) {
      const { content } = await request.post(API_ROUTES.associativeTable, { name, content: data });

      return content;
    },

    update({ dispatch }, payload) {
      return dispatch('create', payload);
    },

    remove(context, { name } = {}) {
      return request.delete(API_ROUTES.associativeTable, { params: { name } });
    },
  },
};
