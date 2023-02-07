import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    items: {},
    pending: true,
  },
  getters: {
    items: state => state.items,
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_LIST]: (state) => {
      state.items = {};
      state.pending = true;
    },
    [types.FETCH_LIST_COMPLETED]: (state, { items }) => {
      state.items = items;
      state.pending = false;
    },
    [types.FETCH_LIST_FAILED]: (state) => {
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ commit }, { params } = {}) {
      try {
        commit(types.FETCH_LIST);

        const items = await request.get(API_ROUTES.templateVars, params);

        commit(types.FETCH_LIST_COMPLETED, { items });
      } catch (err) {
        console.warn(err);

        commit(types.FETCH_LIST_FAILED);
      }
    },
  },
};
