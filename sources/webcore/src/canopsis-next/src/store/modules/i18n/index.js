import Vue from 'vue';

import { DEFAULT_LOCALE, LOCALE_PRIORITIES } from '@/config';

import i18n from '@/i18n';

export const types = {
  SET_LOCALE_PRIORITY: 'SET_LOCALE_PRIORITY',
};

export default {
  namespaced: true,
  state: {
    localePriority: 0,
  },
  mutations: {
    [types.SET_LOCALE_PRIORITY]: (state, { priority }) => {
      state.localePriority = priority;
    },
  },
  actions: {
    setLocale({ state, commit }, { locale, priority = LOCALE_PRIORITIES.default }) {
      if (state.localePriority <= priority) {
        Vue.set(i18n, 'locale', locale);

        commit(types.SET_LOCALE_PRIORITY, { priority });
      }
    },

    setDefaultLocale({ dispatch }, locale = DEFAULT_LOCALE) {
      return dispatch('setLocale', { locale, priority: LOCALE_PRIORITIES.default });
    },

    setGlobalLocale({ dispatch }, locale) {
      return dispatch('setLocale', { locale, priority: LOCALE_PRIORITIES.global });
    },

    setPersonalLocale({ dispatch }, locale) {
      return dispatch('setLocale', { locale, priority: LOCALE_PRIORITIES.personal });
    },
  },
};
