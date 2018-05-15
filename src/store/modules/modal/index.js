const types = {
  SHOW: 'SHOW',
  HIDE: 'HIDE',
  HIDE_COMPLETED: 'HIDE_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    name: null,
    config: {},
    animationPending: false,
  },
  getters: {
    name: state => state.name,
    config: state => state.config,
    animationPending: state => state.animationPending,
  },
  mutations: {
    [types.SHOW](state, { name, config = {} }) {
      state.name = name;
      state.config = config;
      state.animationPending = false;
    },
    [types.HIDE](state) {
      state.animationPending = true;
    },
    [types.HIDE_COMPLETED](state) {
      state.name = null;
      state.config = {};
      state.animationPending = false;
    },
  },
  actions: {
    show({ commit }, payload) {
      commit(types.SHOW, payload);
    },
    hide({ commit, state }) {
      commit(types.HIDE);

      /**
       * This function added for vuetify animation waiting
       */
      setTimeout(() => {
        if (state.animationPending) {
          commit(types.HIDE_COMPLETED);
        }
      }, 300);
    },
  },
};
