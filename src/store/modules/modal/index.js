import Vue from 'vue';
import isEmpty from 'lodash/isEmpty';

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
    modals: {},
  },
  getters: {
    modals: state => state.modals,
    hasModals: state => !isEmpty(state.modals),
  },
  mutations: {
    [types.SHOW](state, { id, name, config = {} }) {
      Vue.set(state.modals, id, {
        id,
        name,
        config,
        hidden: false,
      });
    },
    [types.HIDE](state, { id }) {
      Vue.set(state.modals, id, {
        ...state.modals[id],
        hidden: true,
      });
    },
    [types.HIDE_COMPLETED](state, { id }) {
      Vue.delete(state.modals, id);
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
        if (state.modals[id]) {
          commit(types.HIDE_COMPLETED, { id });
        }
      }, VUETIFY_ANIMATION_DELAY);
    },
  },
};
