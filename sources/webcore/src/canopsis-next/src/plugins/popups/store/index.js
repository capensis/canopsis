import { POPUP_TYPES } from '@/constants';

import uid from '@/helpers/uid';
import { POPUP_AUTO_CLOSE_DELAY } from '@/config';

export const types = {
  ADD: 'ADD',
  REMOVE: 'REMOVE',
  SET_DEFAULT_CLOSE_TIME: 'SET_DEFAULT_CLOSE_TIME',
};

export default {
  namespaced: true,
  state: {
    defaultCloseTime: POPUP_AUTO_CLOSE_DELAY,
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
    [types.SET_DEFAULT_CLOSE_TIME](state, time) {
      state.defaultCloseTime = time;
    },
  },
  actions: {
    add({ commit, state }, {
      id = uid('popup'),
      type,
      text,
      autoClose = state.defaultCloseTime,
    } = {}) {
      commit(types.ADD, {
        popup: {
          id, type, text, autoClose,
        },
      });
    },

    setDefaultCloseTime({ commit }, defaultCloseTime = POPUP_AUTO_CLOSE_DELAY) {
      return commit(types.SET_DEFAULT_CLOSE_TIME, defaultCloseTime);
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
