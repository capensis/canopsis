import Vue from 'vue';
import { merge } from 'lodash';

export const types = {
  UPDATE: 'UPDATE',
  MERGE: 'MERGE',
  REMOVE: 'REMOVE',

  FORCE_UPDATE: 'FORCE_UPDATE',

  UPDATE_LOCKED: 'UPDATE_LOCKED',
  REMOVE_LOCKED: 'REMOVE_LOCKED',
};

export default {
  namespaced: true,
  state: {
    queries: {},
    queriesNonces: {},
    lockedQueries: {},
  },
  getters: {
    getQueryById: state => id => ({ ...state.queries[id], ...state.lockedQueries[id] }),
    getQueryNonceById: state => id => state.queriesNonces[id] || 0,
  },
  mutations: {
    [types.UPDATE](state, { id, query }) {
      Vue.set(state.queries, id, query);
    },

    [types.MERGE](state, { id, query }) {
      Vue.set(state.queries, id, merge({}, state.queries[id], query));
    },

    [types.REMOVE](state, { id }) {
      Vue.delete(state.queries, id);
      Vue.delete(state.queriesNonces, id);
      Vue.delete(state.lockedQueries, id);
    },

    [types.FORCE_UPDATE](state, { id }) {
      const now = new Date();

      Vue.set(state.queriesNonces, id, now.getTime());
    },

    [types.UPDATE_LOCKED](state, { id, query }) {
      Vue.set(state.lockedQueries, id, query);
    },

    [types.REMOVE_LOCKED](state, { id }) {
      Vue.delete(state.lockedQueries, id);
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

    forceUpdate({ commit }, { id }) {
      commit(types.FORCE_UPDATE, { id });
    },

    updateLocked({ commit }, { id, query }) {
      commit(types.UPDATE_LOCKED, { id, query });
    },

    removeLocked({ commit }, { id }) {
      commit(types.REMOVE_LOCKED, { id });
    },
  },
};
