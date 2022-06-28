import { keyBy } from 'lodash';

import {
  API_ROUTES,
  DEFAULT_LOCALE,
  VUETIFY_ANIMATION_DELAY,
  LOCAL_STORAGE_ACCESS_TOKEN_KEY,
} from '@/config';
import { EXCLUDED_SERVER_ERROR_STATUSES } from '@/constants';

import router from '@/router';

import request from '@/services/request';
import localStorageService from '@/services/local-storage';

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
    isLoggedIn: localStorageService.has(LOCAL_STORAGE_ACCESS_TOKEN_KEY),
    currentUser: {},
    pending: true,
  },
  getters: {
    isLoggedIn: state => state.isLoggedIn,
    currentUser: state => state.currentUser,
    currentUserPermissionsById: state => keyBy(state.currentUser.permissions, '_id'),
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
        const { access_token: accessToken } = await request.post(API_ROUTES.login, credentials);

        await dispatch('applyAccessToken', accessToken);
      } catch (err) {
        console.error(err);
        commit(types.LOGOUT);

        throw err;
      }
    },

    applyAccessToken({ commit, dispatch }, accessToken) {
      localStorageService.set(LOCAL_STORAGE_ACCESS_TOKEN_KEY, accessToken);

      commit(types.LOGIN_COMPLETED);

      return Promise.all([
        dispatch('fetchCurrentUser'),
        dispatch('filesAccess'),
      ]);
    },

    filesAccess() {
      return request.get(API_ROUTES.fileAccess);
    },

    async fetchCurrentUser({ commit, dispatch, state }) {
      if (!state.isLoggedIn) {
        return commit(types.LOGOUT);
      }

      try {
        commit(types.FETCH_USER);

        const currentUser = await request.get(API_ROUTES.currentUser, { fullResponse: true });

        if (currentUser.ui_language) {
          dispatch('i18n/setPersonalLocale', currentUser.ui_language, { root: true });
        } else {
          dispatch('i18n/setDefaultLocale', DEFAULT_LOCALE, { root: true });
        }

        return commit(types.FETCH_USER_COMPLETED, currentUser);
      } catch (err) {
        if (EXCLUDED_SERVER_ERROR_STATUSES.includes(err.status)) {
          dispatch('logout');
        }

        throw err;
      }
    },

    async logout({ commit }, { redirectTo } = {}) {
      try {
        await request.post(API_ROUTES.logout);

        commit(types.LOGOUT);
        localStorageService.remove(LOCAL_STORAGE_ACCESS_TOKEN_KEY);

        if (redirectTo) {
          await router.replaceAsync(redirectTo);
        }

        /**
         * We've added timeout for the correct layout padding displaying with transition.
         * And we've added location.reload for refreshing every js objects (store, components states and etc.)
         */
        setTimeout(() => window.location.reload(), VUETIFY_ANIMATION_DELAY);
      } catch (err) {
        console.error(err);
      }
    },

    fetchLoggedUsersCountWithoutStore() {
      return request.get(API_ROUTES.loggedUserCount);
    },
  },
};
