import Cookies from 'js-cookie';
import qs from 'qs';

import request from '@/services/request';
import { API_ROUTES, COOKIE_SESSION_KEY } from '@/config';

const types = {
  LOGIN: 'LOGIN',
  LOGIN_COMPLETED: 'LOGIN_COMPLETED',
  LOGIN_FAILED: 'LOGIN_FAILED',

  LOGOUT: 'LOGOUT',

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
    [types.LOGOUT](state) {
      state.isLoggedIn = false;
      state.currentUser = {};
      state.pending = false;
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
      try {
        await request.post(API_ROUTES.auth, qs.stringify({ ...credentials, json_response: true }), {
          headers: { 'content-type': 'application/x-www-form-urlencoded' },
        });

        commit(types.LOGIN_COMPLETED);

        return dispatch('fetchCurrentUser');
      } catch (err) {
        commit(types.LOGOUT);

        throw err;
      }
    },
    async fetchCurrentUser({ commit, dispatch, state }) {
      if (!state.isLoggedIn) {
        return commit(types.LOGOUT);
      }

      try {
        commit(types.FETCH_USER);

        const { data: [currentUser] } = await request.get(API_ROUTES.currentUser);

        if (currentUser.ui_language) {
          dispatch('i18n/setLocale', currentUser.ui_language, { root: true });
        }

        return commit(types.FETCH_USER_COMPLETED, currentUser);
      } catch (err) {
        dispatch('logout');

        throw err;
      }
    },
    async logout({ commit }) {
      try {
        commit(types.LOGOUT);
        Cookies.remove(COOKIE_SESSION_KEY);

        window.location.reload();
      } catch (err) {
        console.error(err);
      }
    },
  },
};
