import Vue from 'vue';

export const types = {
  UPDATE: 'UPDATE',
  MERGE: 'MERGE',

  FORCE_UPDATE: 'FORCE_UPDATE',
};

export default {
  namespaced: true,
  state: {
    queries: {},
    queriesNonces: {},
  },
  getters: {
    getQueryById: state => id => state.queries[id] || {},
    getQueryNonceById: state => id => state.queriesNonces[id] || 0,
  },
  mutations: {
    [types.UPDATE](state, { id, query }) {
      Vue.set(state.queries, id, query);
    },
    [types.MERGE](state, { id, query }) {
      Vue.set(state.queries, id, { ...state.queries[id], ...query });
    },
    [types.FORCE_UPDATE](state, { id }) {
      const now = new Date();

      Vue.set(state.queriesNonces, id, now.getTime());
    },
  },
  actions: {
    update({ commit }, { id, query }) {
      commit(types.UPDATE, { id, query });
    },
    merge({ commit }, { id, query }) {
      commit(types.MERGE, { id, query });
    },
    forceUpdate({ commit }, { id }) {
      commit(types.FORCE_UPDATE, { id });
    },
  },
};
