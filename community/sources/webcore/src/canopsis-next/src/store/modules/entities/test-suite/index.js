import Vue from 'vue';
import { get } from 'lodash';

import { ENTITIES_TYPES } from '@/constants';
import { API_ROUTES } from '@/config';

import request from '@/services/request';
import { testSuiteSchema } from '@/store/schemas';

import historyModule from './history';
import entityGanttModule from './entity-gantt';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  modules: {
    history: historyModule,
    entityGantt: entityGanttModule,
  },
  state: {
    widgets: {},
  },
  getters: {
    getListByWidgetId: (state, getters, rootState, rootGetters) => widgetId => rootGetters['entities/getList'](
      ENTITIES_TYPES.testSuite,
      get(state.widgets[widgetId], 'allIds', []),
    ),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending'),
    getItem: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](
      ENTITIES_TYPES.testSuite,
      id,
    ),
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: true, error: null });
    },
    [types.FETCH_LIST_COMPLETED](state, { widgetId, allIds }) {
      Vue.setSeveral(state.widgets, widgetId, { allIds, pending: false });
    },
    [types.FETCH_LIST_FAILED](state, { widgetId, error = {} }) {
      Vue.setSeveral(state.widgets, widgetId, { error, pending: false });
    },
  },
  actions: {
    validateDirectory(context, { data }) {
      return request.post(API_ROUTES.junit.directory, data);
    },

    async fetchList({ dispatch, commit }, { widgetId, params } = {}) {
      commit(types.FETCH_LIST, { widgetId });

      const { normalizedData, data } = await dispatch('entities/fetch', {
        route: `${API_ROUTES.junit.widget}/${widgetId}`,
        schema: [testSuiteSchema],
        dataPreparer: d => d.data,
        params,
      }, { root: true });

      commit(types.FETCH_LIST_COMPLETED, {
        widgetId,
        allIds: normalizedData.result,
        ...data,
      });
    },

    fetchSummaryWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.junit.testSuites}/${id}/summary`);
    },

    fetchItemGanttIntervalsWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.junit.testSuites}/${id}/gantt`, { params });
    },

    fetchItemDetailsWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.junit.testSuites}/${id}/details`, { params });
    },
  },
};
