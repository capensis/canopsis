import uid from 'uid';

export const types = {
  ADD: 'ADD',
  REMOVE: 'REMOVE',
};

export default {
  namespaced: true,
  state: {
    popups: [],
  },
  getters: {
    popups: state => state.popups,
  },
  mutations: {
    [types.ADD](state, { popup }) {
      state.popups.push(popup);
    },
    [types.REMOVE](state, { id }) {
      state.popups = state.popups.filter(v => v.id !== id);
    },
  },
  actions: {
    add({ commit }, { popup }) {
      commit(types.ADD, { popup: { id: uid(), ...popup } });
    },
    remove({ commit }, { id }) {
      commit(types.REMOVE, { id });
    },
  },
};
