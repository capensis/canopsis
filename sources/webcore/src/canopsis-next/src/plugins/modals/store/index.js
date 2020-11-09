import Vue from 'vue';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import uid from '@/helpers/uid';
import { isOmitEqual } from '@/helpers/is-omit-equal';

export const types = {
  SHOW: 'SHOW',
  HIDE: 'HIDE',
  HIDE_COMPLETED: 'HIDE_COMPLETED',
  MINIMIZE: 'MINIMIZE',
  MAXIMIZE: 'MAXIMIZE',
};

/**
 * Default modals config comparator
 *
 * @param {Object|Array} value
 * @param {Object|Array} other
 * @returns {boolean}
 */
export const defaultConfigComparator = (value, other) => isOmitEqual(value, other, [
  'action',
  'cancel',
  'afterSubmit',
  'continueAction',
  'continueWithTicketAction',
  'onCreate',
]);

export default {
  namespaced: true,
  state: {
    allIds: [],
    byId: {},
  },
  getters: {
    modals: state => state.allIds.map(id => state.byId[id]),
    hasMaximizedModal: state => state.allIds.some(id => !state.byId[id].minimized),
    hasMinimizedModal: state => state.allIds.some(id => state.byId[id].minimized),
  },
  mutations: {
    [types.SHOW](state, {
      id,
      name,
      config = {},
      dialogProps = {},
    }) {
      Vue.set(state.byId, id, {
        id,
        name,
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
     * @param {Function} [configComparator = defaultConfigComparator]
     */
    show({ commit, state }, {
      name,
      config = {},
      dialogProps = {},
      id = uid('modal'),
      configComparator = defaultConfigComparator,
    } = {}) {
      const sameMinimizedModalId = state.allIds.find(modalId =>
        state.byId[modalId]
        && state.byId[modalId].name === name
        && state.byId[modalId].minimized
        && configComparator(state.byId[modalId].config, config));

      if (sameMinimizedModalId) {
        return commit(types.MAXIMIZE, { id: sameMinimizedModalId });
      }

      return commit(types.SHOW, {
        id,
        name,
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

      commit(types.HIDE, { id });

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
