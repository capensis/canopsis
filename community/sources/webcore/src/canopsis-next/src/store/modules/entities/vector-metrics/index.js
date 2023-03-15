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
    getMetaByWidgetId: state => widgetId => state.widgets[widgetId]?.meta ?? {},
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: true, error: null });
    },

    [types.FETCH_LIST_COMPLETED]: (state, { widgetId, metrics, meta }) => {
      Vue.setSeveral(state.widgets, widgetId, { widgetId, metrics, meta, pending: false });
    },
  },
  actions: {
    async fetchList({ commit }, { widgetId, params } = {}) {
      commit(types.FETCH_LIST, { widgetId });

      const { data: metrics, meta } = await request.get(API_ROUTES.metrics.alarm, { params });

      commit(types.FETCH_LIST_COMPLETED, {
        widgetId,
        metrics,
        meta,
      });
    },
  },
};
