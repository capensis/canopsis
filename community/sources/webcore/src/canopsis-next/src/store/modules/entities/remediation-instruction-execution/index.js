import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchItemWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.remediation.executions}/${id}`, { params });
    },

    fetchPausedListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.remediation.pausedExecutions, { params });
    },

    create(context, { data } = {}) {
      return request.post(API_ROUTES.remediation.executions, data);
    },

    update(context, { path, id, data }) {
      return request.put(`${API_ROUTES.remediation.executions}/${id}/${path}`, data);
    },

    cancel({ dispatch }, { id }) {
      return dispatch('update', { path: 'cancel', id });
    },

    nextOperation({ dispatch }, { id }) {
      return dispatch('update', { path: 'next', id });
    },

    nextStep({ dispatch }, { id, data }) {
      return dispatch('update', { path: 'next-step', id, data });
    },

    pause({ dispatch }, { id }) {
      return dispatch('update', { path: 'pause', id });
    },

    previousOperation({ dispatch }, { id }) {
      return dispatch('update', { path: 'previous', id });
    },

    resume({ dispatch }, { id }) {
      return dispatch('update', { path: 'resume', id });
    },

    rate({ dispatch }, { id, data }) {
      return dispatch('update', { path: 'rate', id, data });
    },
  },
};
