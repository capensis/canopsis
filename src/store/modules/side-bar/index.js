import { VUETIFY_ANIMATION_DELAY } from '@/config';

export const types = {
  SHOW: 'SHOW',
  HIDE: 'HIDE',
  HIDE_COMPLETED: 'HIDE_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    name: null,
    config: {},
    hidden: false,
  },
  getters: {
    name: state => state.name,
    config: state => state.config,
    hidden: state => state.hidden,
  },
  mutations: {
    [types.SHOW](state, { name, config = {} }) {
      state.name = name;
      state.config = config;
      state.hidden = false;
    },
    [types.HIDE](state) {
      state.hidden = true;
    },
    [types.HIDE_COMPLETED](state) {
      state.name = null;
      state.config = {};
      state.hidden = false;
    },
  },
  actions: {
    show({ commit, state }, payload) {
      if (!state.name) {
        commit(types.SHOW, payload);
      }
    },
    hide({ commit, state }) {
      commit(types.HIDE);

      /**
       * This function added for vuetify animation waiting
       */
      setTimeout(() => {
        if (state.hidden) {
          commit(types.HIDE_COMPLETED);
        }
      }, VUETIFY_ANIMATION_DELAY);
    },
  },
};
