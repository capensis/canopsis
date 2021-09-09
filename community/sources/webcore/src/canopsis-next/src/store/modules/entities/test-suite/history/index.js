import Vue from 'vue';
import { get } from 'lodash';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { testSuiteHistorySchema } from '@/store/schemas';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    testSuites: {},
  },
  getters: {
    getListByTestSuiteId: (state, getters, rootState, rootGetters) => testSuiteId => rootGetters['entities/getList'](
      ENTITIES_TYPES.testSuiteHistory,
      get(state.testSuites[testSuiteId], 'allIds', []),
    ),
    getPendingByTestSuiteId: state => testSuiteId => get(state.testSuites[testSuiteId], 'pending'),
    getMetaByTestSuiteId: state => testSuiteId => get(state.testSuites[testSuiteId], 'meta', {}),
  },
  mutations: {
    [types.FETCH_LIST](state, { id }) {
      Vue.setSeveral(state.testSuites, id, { pending: true });
    },
    [types.FETCH_LIST_COMPLETED](state, { id, meta, allIds }) {
      Vue.setSeveral(state.testSuites, id, { allIds, meta, pending: false });
    },
    [types.FETCH_LIST_FAILED](state, { id }) {
      Vue.setSeveral(state.testSuites, id, { pending: false });
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { id, params = {} }) {
      try {
        commit(types.FETCH_LIST, { id });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: `${API_ROUTES.junit.history}/${id}`,
          schema: [testSuiteHistorySchema],
          dataPreparer: d => d.data,
          params,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          id,
          allIds: normalizedData.result,
          meta: data.meta,
        });
      } catch (err) {
        commit(types.FETCH_LIST_FAILED, { id });

        throw err;
      }
    },
  },
};
