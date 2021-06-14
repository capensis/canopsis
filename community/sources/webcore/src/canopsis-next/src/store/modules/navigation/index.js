export const types = {
  TOGGLE_EDITING_MODE: 'TOGGLE_EDITING_MODE',
};

export default {
  namespaced: true,
  state: {
    isEditingMode: false,
  },
  getters: {
    isEditingMode: state => state.isEditingMode,
  },
  mutations: {
    [types.TOGGLE_EDITING_MODE](state) {
      state.isEditingMode = !state.isEditingMode;
    },
  },
  actions: {
    toggleEditingMode({ commit }) {
      commit(types.TOGGLE_EDITING_MODE);
    },
  },
};
