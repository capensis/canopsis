import uid from '@/helpers/uid';

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
    add({ commit }, {
      id = uid('popup'), type, text, autoClose,
    } = {}) {
      commit(types.ADD, {
        popup: {
          id, type, text, autoClose,
        },
      });
    },
    remove({ commit }, { id }) {
      commit(types.REMOVE, { id });
    },
  },
};
