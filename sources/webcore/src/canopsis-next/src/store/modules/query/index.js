import Vue from 'vue';

export const types = {
  UPDATE: 'UPDATE',
  MERGE: 'MERGE',
  REMOVE: 'REMOVE',
};

export default {
  namespaced: true,
  state: {
    queries: {},
  },
  getters: {
    getQueryById: state => id => state.queries[id] || {},
  },
  mutations: {
    [types.UPDATE](state, { id, query }) {
      Vue.set(state.queries, id, query);
    },

    [types.MERGE](state, { id, query }) {
      Vue.set(state.queries, id, { ...state.queries[id], ...query });
    },

    [types.REMOVE](state, { id }) {
      Vue.delete(state.queries, id);
    },
  },
  actions: {
    update({ commit }, { id, query }) {
      commit(types.UPDATE, { id, query });
    },

    merge({ commit }, { id, query }) {
      commit(types.MERGE, { id, query });
    },

    remove({ commit }, { id }) {
      commit(types.REMOVE, { id });
    },
  },
};
