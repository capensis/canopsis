import Vue from 'vue';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import uid from '@/helpers/uid';

export const types = {
  SHOW: 'SHOW',
  HIDE: 'HIDE',
  HIDE_COMPLETED: 'HIDE_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    byId: {},
  },
  getters: {
    modals: state => state.allIds.map(id => state.byId[id]),
    hasModals: state => Boolean(state.allIds.length),
  },
  mutations: {
    [types.SHOW](state, { id, name, config = {} }) {
      Vue.set(state.byId, id, {
        id,
        name,
        config,
        hidden: false,
      });

      state.allIds.push(id);
    },
    [types.HIDE](state, { id }) {
      Vue.set(state.byId, id, {
        ...state.byId[id],
        hidden: true,
      });
    },
    [types.HIDE_COMPLETED](state, { id }) {
      state.allIds = state.allIds.filter(value => value !== id);

      Vue.delete(state.byId, id);
    },
  },
  actions: {
    /**
     * Show modal window by name with unique id
     *
     * @param {function} commit
     * @param {string} name
     * @param {Object} [config={}]
     * @param {string} [id=uid()]
     */
    show({ commit }, { name, config = {}, id = uid('modal') } = {}) {
      commit(types.SHOW, { id, name, config });
    },
    /**
     * Hide modal by id
     *
     * @param {function} commit
     * @param {Object} state
     * @param {string} [id=uid()]
     */
    hide({ commit, state }, { id } = {}) {
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
  },
};
