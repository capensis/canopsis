import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    /**
     * Fetch permissions list by params without store
     *
     * @param {VuexActionContext} context
     * @param {Object} params
     * @returns {Object}
     */
    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.permissions, { params });
    },
  },
};
