export default {
  namespaced: true,
  state: {
    isSideBarOpen: false,
  },
  mutations: {
    toggleSideBar(state) {
      state.isSideBarOpen = !state.isSideBarOpen;
    },
  },
  actions: {
    toggleSideBar(context) {
      context.commit('toggleSideBar');
    },
  },
};
