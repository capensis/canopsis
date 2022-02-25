import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { viewTabSchema } from '@/store/schemas';

export default {
  namespaced: true,
  actions: {
    async fetchItem({ dispatch }, { id, params } = {}) {
      const { data } = await dispatch('entities/fetch', {
        params,
        route: `${API_ROUTES.view.tabs}/${id}`,
        schema: viewTabSchema,
      }, { root: true });

      return data;
    },

    create(context, { data } = {}) {
      return request.post(API_ROUTES.view.tabs, data);
    },

    clone(context, { data, id } = {}) {
      return request.put(`${API_ROUTES.view.tabs}/${id}/clone`, data);
    },

    update(context, { data, id } = {}) {
      return request.put(`${API_ROUTES.view.tabs}/${id}`, data);
    },

    remove(context, { id } = {}) {
      return request.delete(`${API_ROUTES.view.tabs}/${id}`);
    },

    copy(context, { id, data } = {}) {
      return request.post(`${API_ROUTES.view.tabCopy}/${id}`, data);
    },

    updatePositions(context, { data } = {}) {
      return request.put(API_ROUTES.view.tabPositions, data);
    },
  },
};
