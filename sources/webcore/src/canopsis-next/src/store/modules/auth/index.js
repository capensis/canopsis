import Cookies from 'js-cookie';
import qs from 'qs';

import request from '@/services/request';
import { BASE_URL, API_ROUTES, COOKIE_SESSION_KEY, DEFAULT_LOCALE } from '@/config';

const types = {
  LOGIN: 'LOGIN',
  LOGIN_COMPLETED: 'LOGIN_COMPLETED',
  LOGIN_FAILED: 'LOGIN_FAILED',

  FETCH_USER: 'FETCH_USER',
  FETCH_USER_COMPLETED: 'FETCH_USER_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    isLoggedIn: !!Cookies.get(COOKIE_SESSION_KEY),
    currentUser: {},
    pending: true,
  },
  getters: {
    isLoggedIn: state => state.isLoggedIn,
    currentUser: state => state.currentUser,
    pending: state => state.pending,
  },
  mutations: {
    [types.LOGIN_COMPLETED](state) {
      state.isLoggedIn = true;
    },
    [types.FETCH_USER](state) {
      state.pending = true;
    },
    [types.FETCH_USER_COMPLETED](state, currentUser) {
      state.currentUser = currentUser;
      state.pending = false;
    },
  },
  actions: {
    async login({ commit, dispatch }, credentials) {
      await request.post(API_ROUTES.auth, qs.stringify({ ...credentials, json_response: true }), {
        headers: { 'content-type': 'application/x-www-form-urlencoded' },
      });

      commit(types.LOGIN_COMPLETED);

      return dispatch('fetchCurrentUser');
    },
    async fetchCurrentUser({ commit, dispatch, state }) {
      try {
        if (state.isLoggedIn) {
          commit(types.FETCH_USER);

          const { data: [currentUser] } = await request.get(API_ROUTES.currentUser);

          if (currentUser.ui_language) {
            dispatch('i18n/setLocale', currentUser.ui_language, { root: true });
          } else {
            dispatch('i18n/setLocale', DEFAULT_LOCALE, { root: true });
          }

          commit(types.FETCH_USER_COMPLETED, currentUser);
        }
      } catch (err) {
        dispatch('logout');
      }
    },
    logout() {
      window.location = `${BASE_URL}${API_ROUTES.logout}`;
    },
  },
};
