import Vue from 'vue';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { uid } from '@/helpers/uid';

export const types = {
  SHOW: 'SHOW',
  HIDE: 'HIDE',
  HIDE_COMPLETED: 'HIDE_COMPLETED',
  MINIMIZE: 'MINIMIZE',
  MAXIMIZE: 'MAXIMIZE',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    byId: {},
  },
  getters: {
    sidebars: state => state.allIds.map(id => state.byId[id]),
    sidebarsById: state => state.byId,
  },
  mutations: {
    [types.SHOW](state, {
      id,
      name,
      config = {},
    }) {
      Vue.set(state.byId, id, {
        id,
        name,
        config,
        hidden: false,
        minimized: false,
      });

      state.allIds.push(id);
    },
    [types.HIDE](state, { id }) {
      Vue.set(state.byId[id], 'hidden', true);
    },
    [types.HIDE_COMPLETED](state, { id }) {
      state.allIds = state.allIds.filter(value => value !== id);

      Vue.delete(state.byId, id);
    },
    [types.MINIMIZE](state, { id }) {
      Vue.set(state.byId[id], 'minimized', true);
    },
    [types.MAXIMIZE](state, { id }) {
      Vue.set(state.byId[id], 'minimized', false);
    },
  },
  actions: {
    /**
     * Show a sidebar by component name with a stable id (multiple sidebars may be open).
     *
     * @param {Function} commit
     * @param {Object} state
     * @param {string} name
     * @param {Object} [config = {}]
     * @param {string} [id = uid()]
     */
    show({ commit, state }, {
      name,
      config = {},
      id = uid('sidebar'),
    } = {}) {
      if (state.byId[id]) {
        return commit(types.MAXIMIZE, { id });
      }

      return commit(types.SHOW, {
        id,
        name,
        config,
      });
    },

    /**
     * Mark sidebar hidden; removes it from the store after the drawer animation.
     *
     * @param {Function} commit
     * @param {Object} state
     * @param {string} id
     */
    hide({ commit, state }, { id } = {}) {
      if (!id) {
        throw new Error('Missed required parameter');
      }

      if (!state.byId[id]) {
        return;
      }

      commit(types.HIDE, { id });

      setTimeout(() => {
        if (state.byId[id]) {
          commit(types.HIDE_COMPLETED, { id });
        }
      }, VUETIFY_ANIMATION_DELAY);
    },

    /**
     * @param {Function} commit
     * @param {string} id
     */
    minimize({ commit }, { id } = {}) {
      if (!id) {
        throw new Error('Missed required parameter');
      }

      commit(types.MINIMIZE, { id });
    },

    /**
     * @param {Function} commit
     * @param {string} id
     */
    maximize({ commit }, { id } = {}) {
      if (!id) {
        throw new Error('Missed required parameter');
      }

      commit(types.MAXIMIZE, { id });
    },
  },
};
