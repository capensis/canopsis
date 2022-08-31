import { API_ROUTES } from '@/config';

import request from '@/services/request';

const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    pending: false,
    metrics: [],
    meta: {},
    fetchingParams: {},
  },
  getters: {
    metrics: state => state.metrics,
    meta: state => state.meta,
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_LIST]: (state, { params }) => {
      state.fetchingParams = params;
      state.pending = true;
    },

    [types.FETCH_LIST_COMPLETED]: (state, { metrics, meta }) => {
      state.metrics = metrics;
      state.meta = meta;
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ commit }, { params } = {}) {
      commit(types.FETCH_LIST, { params });

      const { data: metrics, meta } = await request.get(API_ROUTES.metrics.remediation, { params });

      commit(types.FETCH_LIST_COMPLETED, { metrics, meta });
    },

    fetchListWithPreviousParams({ dispatch, state }) {
      return dispatch('fetchList', { params: state.fetchingParams });
    },
  },
};
