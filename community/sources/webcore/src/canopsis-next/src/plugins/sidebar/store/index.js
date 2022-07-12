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
    hidden: true,
  },
  getters: {
    sidebar: state => state,
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
      state.hidden = true;
    },
  },
  actions: {
    show({ commit }, { name, config = {} } = {}) {
      commit(types.SHOW, { name, config });
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
