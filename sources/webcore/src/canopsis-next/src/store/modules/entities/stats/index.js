import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    async fetchListWithoutStore(context, { params }) {
      return request.post(API_ROUTES.stats, params);
    },

    async fetchEvolutionWithoutStore(context, { params }) {
      return request.post(`${API_ROUTES.stats}/evolution`, params);
    },
  },
};
