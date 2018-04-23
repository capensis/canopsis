const types = {
  SHOW: 'SHOW',
  HIDE: 'HIDE',
};

export default {
  namespaced: true,
  state: {
    component: null,
    config: {},
  },
  getters: {
    component: state => state.component,
    config: state => state.config,
  },
  mutations: {
    [types.SHOW](state, { componentName, config }) {
      state.component = componentName;
      state.config = config;
    },
    [types.HIDE](state) {
      state.component = null;
      state.config = {};
    },
  },
  actions: {
    show({ commit }, payload) {
      commit(types.SHOW, payload);
    },
    hide({ commit }) {
      commit(types.HIDE);
    },
  },
};
