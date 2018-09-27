import Vue from 'vue';
import get from 'lodash/get';
import i18n from '@/i18n';

import { API_ROUTES } from '@/config';
import { statSchema } from '@/store/schemas';
import { ENTITIES_TYPES } from '@/constants';

export const types = {
  FETCH_STATS: 'FETCH_STATS',
  FETCH_STATS_COMPLETED: 'FETCH_STATS_COMPLETED',
  FETCH_STATS_FAILED: 'FETCH_STATS_FAILED',
};

export default {
  namespaced: true,
  state: {
    widgets: {},
  },
  getters: {
    getListByWidgetId: (state, getters, rootState, rootGetters) => widgetId =>
      rootGetters['entities/getList'](ENTITIES_TYPES.stat, get(state.widgets[widgetId], 'allIds', [])),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending'),
  },
  mutations: {
    [types.FETCH_STATS](state, { widgetId, params }) {
      Vue.set(state.widgets, widgetId, {
        ...state.widgets[widgetId],
        pending: true,
        fetchingParams: params,
      });
    },
    [types.FETCH_STATS_COMPLETED](state, { widgetId, allIds }) {
      Vue.set(state.widgets, widgetId, {
        ...state.widgets[widgetId],
        pending: false,
        allIds,
      });
    },
    [types.FETCH_STATS_FAILED](state, { widgetId }) {
      Vue.set(state.widgets, widgetId, {
        ...state.widgets[widgetId],
        pending: false,
      });
    },
  },
  actions: {
    async fetchStats({ commit, dispatch }, { widgetId, params } = {}) {
      commit(types.FETCH_STATS, { widgetId, params });
      try {
        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.stats,
          schema: [statSchema],
          body: params,
          method: 'POST',
          dataPreparer: d => d.values,
        }, { root: true });


        commit(types.FETCH_STATS_COMPLETED, {
          widgetId,
          allIds: normalizedData.result,
        });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: `${i18n.t('errors.default')} : ${err.description}` }, { root: true });

        commit(types.FETCH_STATS_FAILED, { widgetId });
      }
    },
  },
};
