import Vue from 'vue';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { uid } from '@/helpers/uid';

export const types = {
  SHOW: 'SHOW',
  HIDE: 'HIDE',
  HIDE_COMPLETED: 'HIDE_COMPLETED',
  REGISTER_ON_HIDE: 'REGISTER_ON_HIDE',
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
    modals: state => state.allIds.map(id => state.byId[id]),
    hasMaximizedModal: state => state.allIds.some(id => !state.byId[id].minimized),
  },
  mutations: {
    [types.SHOW](state, {
      id,
      name,
      onHide,
      config = {},
      dialogProps = {},
    }) {
      Vue.set(state.byId, id, {
        id,
        name,
        onHide,
        config,
        dialogProps,
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
    [types.REGISTER_ON_HIDE](state, { id, callback }) {
      Vue.set(state.byId[id], 'onHide', callback);
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
     * Show modal window by name with unique id
     *
     * @param {Function} commit
     * @param {Object} state
     * @param {string} name
     * @param {Object} [config = {}]
     * @param {Object} [dialogProps = {}]
     * @param {string} [id = uid()]
     */
    show({ commit, state }, {
      name,
      onHide,
      config = {},
      dialogProps = {},
      id = uid('modal'),
    } = {}) {
      if (state.byId[id]) {
        return commit(types.MAXIMIZE, { id });
      }

      return commit(types.SHOW, {
        id,
        name,
        onHide,
        config,
        dialogProps,
      });
    },

    /**
     * Hide modal by id
     *
     * @param {Function} commit
     * @param {Object} state
     * @param {string} [id]
     */
    hide({ commit, state }, { id } = {}) {
      if (!id) {
        throw new Error('Missed required parameter');
      }

      if (!state.byId[id]) {
        return;
      }

      commit(types.HIDE, { id });

      state.byId[id].onHide?.();

      /**
       * This function added for vuetify animation waiting
       */
      setTimeout(() => {
        if (state.byId[id]) {
          commit(types.HIDE_COMPLETED, { id });
        }
      }, VUETIFY_ANIMATION_DELAY);
    },

    /**
     * Register on hide callback
     *
     * @param {Function} commit
     * @param {Object} state
     * @param {string} id
     * @param {Function} callback
     */
    registerOnHide({ commit, state }, { id, callback } = {}) {
      if (!state.byId[id]) {
        return;
      }

      commit(types.REGISTER_ON_HIDE, { id, callback });
    },

    /**
     * Minimize modal by id
     *
     * @param {Function} commit
     * @param {string} [id]
     */
    minimize({ commit }, { id } = {}) {
      if (!id) {
        throw new Error('Missed required parameter');
      }

      commit(types.MINIMIZE, { id });
    },

    /**
     * Maximize modal by id
     *
     * @param {Function} commit
     * @param {string} [id]
     */
    maximize({ commit }, { id } = {}) {
      if (!id) {
        throw new Error('Missed required parameter');
      }

      commit(types.MAXIMIZE, { id });
    },
  },
};
