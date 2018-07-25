import Cookies from 'js-cookie';

import request from '@/services/request';
import { API_ROUTES, COOKIE_SESSION_KEY } from '@/config';

const types = {
  LOGIN: 'LOGIN',
  LOGIN_COMPLETED: 'LOGIN_COMPLETED',
  LOGIN_FAILED: 'LOGIN_FAILED',

  LOGOUT: 'LOGOUT',

  FETCH_USER: 'FETCH_USER',
  FETCH_USER_COMPLETED: 'FETCH_USER_COMPLETED',
  FETCH_USER_FAILED: 'FETCH_USER_FAILED',
};

export default {
  namespaced: true,
  state: {
    isLoggedIn: !!Cookies.get(COOKIE_SESSION_KEY),
    user: {},
  },
  getter: {
    isLoggedIn: state => state.isLoggedIn,
    user: state => state.user,
  },
  mutations: {
    [types.LOGIN_COMPLETED](state) {
      state.isLoggedIn = true;
    },
    [types.LOGOUT](state) {
      state.isLoggedIn = false;
      state.user = {};
    },
    [types.FETCH_USER]() {
    },
    [types.FETCH_USER_FAILED]() {
    },
    [types.FETCH_USER_COMPLETED](state, user) {
      state.user = user;
    },
  },
  actions: {
    async login({ commit, dispatch }, credentials) {
      try {
        await request.post(API_ROUTES.auth, { ...credentials/* , json_response: true */ });
        commit(types.LOGIN_COMPLETED);

        return dispatch('getCurrentUser');
      } catch (err) {
        commit(types.LOGOUT);

        throw err;
      }
    },
    async getCurrentUser({ commit, dispatch, state }) {
      commit(types.FETCH_USER);

      if (!state.isLoggedIn) {
        commit(types.FETCH_USER_FAILED);

        return dispatch('logout');
      }

      try {
        const { user } = await request.get(API_ROUTES.currentUser);
        return commit(types.FETCH_USER_COMPLETED, user);
      } catch (err) {
        commit(types.FETCH_USER_FAILED);

        throw err;
      }
    },
    async logout({ commit }) {
      commit(types.LOGOUT);
    },
  },
};
