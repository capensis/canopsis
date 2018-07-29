import Cookies from 'js-cookie';

import request from '@/services/request';
import { API_ROUTES, COOKIE_SESSION_KEY } from '@/config';

const types = {
  LOGIN: 'LOGIN',
  LOGIN_COMPLETED: 'LOGIN_COMPLETED',
  LOGIN_FAILED: 'LOGIN_FAILED',

  LOGOUT: 'LOGOUT',

  FETCH_USER_COMPLETED: 'FETCH_USER_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    isLoggedIn: !!Cookies.get(COOKIE_SESSION_KEY),
    currentUser: {},
  },
  getters: {
    isLoggedIn: state => state.isLoggedIn,
    currentUser: state => state.currentUser,
  },
  mutations: {
    [types.LOGIN_COMPLETED](state) {
      state.isLoggedIn = true;
    },
    [types.LOGOUT](state) {
      state.isLoggedIn = false;
      state.currentUser = {};
    },
    [types.FETCH_USER_COMPLETED](state, currentUser) {
      state.currentUser = currentUser;
    },
  },
  actions: {
    async login({ commit, dispatch }, credentials) {
      try {
        await request.post(API_ROUTES.auth, { ...credentials, json_response: true });
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
        const { data } = await request.get(API_ROUTES.currentUser);

        return commit(types.FETCH_USER_COMPLETED, data[0]);
      } catch (err) {
        dispatch('logout');

        throw err;
      }
    },
    async logout({ commit }) {
      try {
        commit(types.LOGOUT);
        await request.get(API_ROUTES.logout);
      } catch (err) {
        console.error(err);
      }
    },
  },
};
