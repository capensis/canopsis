import Vue from 'vue';

export const types = {
  UPDATE: 'UPDATE',
  MERGE: 'MERGE',
  SET_PENDING: 'SET_PENDING',
};

export default {
  namespaced: true,
  state: {
    queries: {},
    pending: {},
  },
  getters: {
    getQueryById: state => id => state.queries[id] || {},
    getQueryPendingById: state => id => state.pending[id],
  },
  mutations: {
    [types.UPDATE](state, { id, query }) {
      Vue.set(state.queries, id, query);
    },
    [types.MERGE](state, { id, query }) {
      Vue.set(state.queries, id, { ...state.queries[id], ...query });
    },
    [types.SET_PENDING](state, { id, value }) {
      Vue.set(state.pending, id, value);
    },
  },
  actions: {
    update({ commit }, { id, query }) {
      commit(types.UPDATE, { id, query });
    },
    merge({ commit }, { id, query }) {
      commit(types.MERGE, { id, query });
    },
    startPending({ commit }, { id }) {
      commit(types.SET_PENDING, { id, value: true });
    },
    stopPending({ commit }, { id }) {
      commit(types.SET_PENDING, { id, value: false });
    },
  },
};
