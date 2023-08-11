import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchItem({ dispatch }, { id, params } = {}) {
      return dispatch('view/group/fetchViewTab', { id, params }, { root: true });
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
