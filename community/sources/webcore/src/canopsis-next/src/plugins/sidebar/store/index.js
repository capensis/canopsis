import { VUETIFY_ANIMATION_DELAY } from '@/config';

export const types = {
  SHOW: 'SHOW',
  HIDE: 'HIDE',
  HIDE_COMPLETED: 'HIDE_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    sidebar: {
      hidden: true,
    },
  },
  getters: {
    sidebar: state => state.sidebar,
  },
  mutations: {
    [types.SHOW](state, { name, config = {} }) {
      state.sidebar = {
        name,
        config,

        hidden: false,
      };
    },
    [types.HIDE](state) {
      state.sidebar.hidden = true;
    },
    [types.HIDE_COMPLETED](state) {
      state.sidebar = {};
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
