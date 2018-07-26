const types = {
  TOGGLE: 'TOGGLE',
};

export default {
  namespaced: true,
  state: {
    isSideBarOpen: false,
  },
  getters: {
    isSideBarOpen: state => state.isSideBarOpen,
  },
  mutations: {
    [types.TOGGLE](state) {
      state.isSideBarOpen = !state.isSideBarOpen;
    },
  },
  actions: {
    toggleSideBar({ commit }) {
      commit(types.TOGGLE);
    },
  },
};
