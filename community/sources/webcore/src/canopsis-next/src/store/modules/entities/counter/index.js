import Vue from 'vue';
import { get } from 'lodash';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import i18n from '@/i18n';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    widgets: {},
  },
  getters: {
    getListByWidgetId: state => widgetId => get(state.widgets[widgetId], 'counters', []),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending', []),
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: true });
    },
    [types.FETCH_LIST_COMPLETED](state, { widgetId, counters }) {
      Vue.setSeveral(state.widgets, widgetId, { counters, pending: false });
    },
    [types.FETCH_LIST_FAILED](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: false });
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { widgetId, params = {}, filters = [] }) {
      try {
        commit(types.FETCH_LIST, { widgetId });

        const requests = filters.map(filter => request.get(API_ROUTES.counter, {
          params: {
            ...params,

            filters: [filter],
          },
        }));

        const counters = await Promise.all(requests);

        commit(types.FETCH_LIST_COMPLETED, {
          widgetId,
          counters,
        });
      } catch (err) {
        console.error(err);

        commit(types.FETCH_LIST_FAILED, { widgetId });

        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });
      }
    },
  },
};
