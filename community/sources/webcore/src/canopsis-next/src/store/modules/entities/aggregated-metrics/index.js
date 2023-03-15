import Vue from 'vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    widgets: {},
  },
  getters: {
    getListByWidgetId: state => widgetId => state.widgets[widgetId]?.metrics ?? [],
    getPendingByWidgetId: state => widgetId => state.widgets[widgetId]?.pending ?? false,
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: true, error: null });
    },

    [types.FETCH_LIST_COMPLETED]: (state, { widgetId, metrics }) => {
      Vue.setSeveral(state.widgets, widgetId, { widgetId, metrics, pending: false });
    },
  },
  actions: {
    async fetchList({ commit }, { widgetId, params } = {}) {
      commit(types.FETCH_LIST, { widgetId });

      const { data: metrics } = await request.post(API_ROUTES.metrics.aggregate, params);

      commit(types.FETCH_LIST_COMPLETED, {
        widgetId,
        metrics,
      });
    },
  },
};
