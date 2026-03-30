import Vue from 'vue';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { SIDE_BARS_MINIMIZABLE_USER_FIELD } from '@/constants';

import { uid } from '@/helpers/uid';

export const types = {
  SHOW: 'SHOW',
  HIDE: 'HIDE',
  HIDE_COMPLETED: 'HIDE_COMPLETED',
  MINIMIZE: 'MINIMIZE',
  MAXIMIZE: 'MAXIMIZE',
  UPDATE_CONFIG: 'UPDATE_CONFIG',
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
      minimized = false,
    }) {
      Vue.set(state.byId, id, {
        id,
        name,
        config,
        minimized,
        hidden: false,
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
    [types.UPDATE_CONFIG](state, { id, config }) {
      Vue.setSeveral(state.byId[id], 'config', config);
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
    show({ commit, state, rootGetters }, {
      name,
      config = {},
      id = uid('sidebar'),
    } = {}) {
      if (state.byId[id]) {
        return commit(types.MAXIMIZE, { id });
      }

      const field = SIDE_BARS_MINIMIZABLE_USER_FIELD[name];
      let minimized = false;

      if (field) {
        minimized = rootGetters['auth/currentUser']?.ui_tours?.[field] ?? false;
      }

      return commit(types.SHOW, {
        id,
        name,
        config,
        minimized,
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
    minimize({ commit, dispatch }, { id } = {}) {
      if (!id) {
        throw new Error('Missed required parameter');
      }

      commit(types.MINIMIZE, { id });

      dispatch('setCurrentUserMinimized', { id, minimized: true });
    },

    /**
     * @param {Function} commit
     * @param {string} id
     */
    maximize({ commit, dispatch }, { id } = {}) {
      if (!id) {
        throw new Error('Missed required parameter');
      }

      commit(types.MAXIMIZE, { id });

      dispatch('setCurrentUserMinimized', { id, minimized: false });
    },

    /**
     * Shallow-merges `config` into the open sidebar’s `config` object (`UPDATE_CONFIG`).
     *
     * @param {Object} context
     * @param {Function} context.commit
     * @param {Object} [payload]
     * @param {string} payload.id - Sidebar instance id.
     * @param {Object} payload.config - Partial config merged into the existing entry.
     */
    updateConfig({ commit }, { id, config } = {}) {
      if (!id) {
        throw new Error('Missed required parameter');
      }

      commit(types.UPDATE_CONFIG, { id, config });
    },

    /**
     * Persists minimized / expanded state for minimizable sidebars on the current user
     * (`user/updateCurrentUserTours` root action, field from `SIDE_BARS_MINIMIZABLE_USER_FIELD`).
     * No-op if the sidebar id is unknown or the bar has no mapped user field.
     *
     * @param {Object} context
     * @param {Object} context.getters
     * @param {Function} context.dispatch
     * @param {Object} [payload]
     * @param {string} payload.id - Sidebar instance id.
     * @param {boolean} payload.minimized - Whether the drawer is minimized.
     */
    async setCurrentUserMinimized({ getters, dispatch }, { id, minimized } = {}) {
      if (!id) {
        throw new Error('Missed required parameter');
      }

      const { name } = getters.sidebarsById[id] ?? {};

      if (!name) {
        return;
      }

      await dispatch('user/updateCurrentUserTours', {
        data: { [SIDE_BARS_MINIMIZABLE_USER_FIELD[name]]: minimized },
      }, { root: true });
    },
  },
};
