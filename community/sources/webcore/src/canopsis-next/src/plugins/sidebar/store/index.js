import Vue from 'vue';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { uid } from '@/helpers/uid';

export const types = {
  SHOW: 'SHOW',
  HIDE: 'HIDE',
  HIDE_COMPLETED: 'HIDE_COMPLETED',
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
    }) {
      Vue.set(state.byId, id, {
        id,
        name,
        config,
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
    [types.UPDATE_CONFIG](state, { id, config = {} }) {
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
    show({ commit }, {
      name,
      config = {},
      id = uid('sidebar'),
    } = {}) {
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
  },
};
