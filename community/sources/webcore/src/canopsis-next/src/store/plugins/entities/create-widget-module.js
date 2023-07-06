import Vue from 'vue';
import { merge } from 'lodash';

import { REQUEST_METHODS } from '@/constants';

import request from '@/services/request';

export const DEFAULT_WIDGET_MODULE_TYPES = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default ({
  types = DEFAULT_WIDGET_MODULE_TYPES,
  method = REQUEST_METHODS.post,
  route,
}, module = {}) => {
  const moduleState = {
    widgets: {},
  };

  const moduleGetters = {
    getListByWidgetId: state => widgetId => state.widgets[widgetId]?.metrics ?? [],
    getPendingByWidgetId: state => widgetId => state.widgets[widgetId]?.pending ?? false,
    getMetaByWidgetId: state => widgetId => state.widgets[widgetId]?.meta ?? {},
  };

  const moduleMutations = {
    [types.FETCH_LIST](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: true, error: null });
    },

    [types.FETCH_LIST_COMPLETED]: (state, { widgetId, metrics, meta }) => {
      Vue.setSeveral(state.widgets, widgetId, { widgetId, metrics, meta, pending: false });
    },

    [types.FETCH_LIST_FAILED]: (state, { widgetId }) => {
      Vue.setSeveral(state.widgets, widgetId, { widgetId, pending: false });
    },
  };

  const moduleActions = {
    async fetchList({ commit }, { widgetId, params } = {}) {
      try {
        commit(types.FETCH_LIST, { widgetId });

        const {
          data: metrics,
          meta,
        } = await request({
          method,
          url: route,
          [method === REQUEST_METHODS.get ? 'params' : 'data']: params,
        });

        commit(types.FETCH_LIST_COMPLETED, {
          widgetId,
          metrics,
          meta,
        });
      } catch (err) {
        console.error(err);
      }
    },
  };

  return merge({
    namespaced: true,
    state: moduleState,
    getters: moduleGetters,
    mutations: moduleMutations,
    actions: moduleActions,
  }, module);
};
