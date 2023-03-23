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
    getPreviousListByWidgetId: state => widgetId => state.widgets[widgetId]?.previousMetrics ?? [],
    getPreviousIntervalByWidgetId: state => widgetId => state.widgets[widgetId]?.previousInterval ?? {},

    getListByWidgetId: (state, getters) => (widgetId) => {
      const metrics = state.widgets[widgetId]?.metrics ?? [];
      const previousMetrics = getters.getPreviousListByWidgetId(widgetId);
      const previousInterval = getters.getPreviousIntervalByWidgetId(widgetId);

      return metrics.map((metric, index) => {
        const previousMetric = previousMetrics[index] ?? {};

        return {
          ...metric,
          previous_metric: previousMetric.value,
          previous_interval: previousInterval,
        };
      });
    },
    getPendingByWidgetId: state => widgetId => state.widgets[widgetId]?.pending ?? false,
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: true, error: null });
    },

    [types.FETCH_LIST_COMPLETED]: (state, { widgetId, metrics, previousMetrics, previousInterval }) => {
      Vue.setSeveral(state.widgets, widgetId, { widgetId, metrics, previousMetrics, previousInterval, pending: false });
    },
  },
  actions: {
    async fetchList({ commit }, { widgetId, trend, params } = {}) {
      commit(types.FETCH_LIST, { widgetId });

      const previousInterval = {
        from: params.from - (params.to - params.from),
        to: params.from,
      };

      const { data: metrics } = await request.post(API_ROUTES.metrics.aggregate, params);

      let previousMetrics = [];

      if (trend) {
        const { data } = await request.post(API_ROUTES.metrics.aggregate, {
          ...params,
          ...previousInterval,
        });

        previousMetrics = data;
      }

      commit(types.FETCH_LIST_COMPLETED, {
        widgetId,
        metrics,
        previousMetrics,
        previousInterval,
      });
    },
  },
};
