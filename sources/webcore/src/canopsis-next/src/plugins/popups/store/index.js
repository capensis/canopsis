import { POPUP_TYPES } from '@/constants';

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

    success({ dispatch }, popup) {
      return dispatch('add', { ...popup, type: POPUP_TYPES.success });
    },
    info({ dispatch }, popup) {
      return dispatch('add', { ...popup, type: POPUP_TYPES.info });
    },
    warning({ dispatch }, popup) {
      return dispatch('add', { ...popup, type: POPUP_TYPES.warning });
    },
    error({ dispatch }, popup) {
      return dispatch('add', { ...popup, type: POPUP_TYPES.error });
    },
  },
};
