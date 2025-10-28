import Vue from 'vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import i18n from '@/i18n';

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
    async fetchListWithoutStore(
      { dispatch },
      { params: { with_history: withHistory, prev_from: prevFrom, prev_to: prevTo, ...params } = {} },
    ) {
      const { data: metrics } = await request.post(API_ROUTES.metrics.aggregate, params);

      if (!withHistory) {
        return { data: metrics };
      }

      try {
        const previousInterval = {
          from: prevFrom ?? params.from - (params.to - params.from),
          to: prevTo ?? params.from,
        };

        const { data: previousMetrics } = await request.post(API_ROUTES.metrics.aggregate, {
          ...params,
          ...previousInterval,
        });

        return {
          data: metrics.map((metric, index) => {
            const previousMetric = previousMetrics[index] ?? {};

            return {
              ...metric,
              previous_metric: previousMetric.value,
              previous_interval: previousInterval,
            };
          }),
        };
      } catch (error) {
        console.error(error);

        dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });

        return { data: metrics };
      }
    },

    async fetchList({ commit, dispatch }, { widgetId, params = {} } = {}) {
      commit(types.FETCH_LIST, { widgetId });

      const { data: metrics } = await dispatch('fetchListWithoutStore', { params });

      commit(types.FETCH_LIST_COMPLETED, {
        widgetId,
        metrics,
      });
    },
  },
};
