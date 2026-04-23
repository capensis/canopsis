import { keyBy } from 'lodash';

import {
  API_ROUTES,
  DEFAULT_LOCALE,
  VUETIFY_ANIMATION_DELAY,
  LOCAL_STORAGE_ACCESS_TOKEN_KEY,
  LOCAL_STORAGE_WARNING_POPUP_KEY,
} from '@/config';
import { EXCLUDED_SERVER_ERROR_STATUSES } from '@/constants';

import request from '@/services/request';
import localStorageService from '@/services/local-storage';

import i18n from '@/i18n';

import { viewPermissionsGroupedPermissions } from '@/helpers/permission';
import { isUserHasOnlyApiRole } from '@/helpers/entities/user/entity';

const types = {
  LOGIN: 'LOGIN',
  LOGIN_COMPLETED: 'LOGIN_COMPLETED',
  LOGIN_FAILED: 'LOGIN_FAILED',

  LOGOUT: 'LOGOUT',

  FETCH_USER: 'FETCH_USER',
  FETCH_USER_COMPLETED: 'FETCH_USER_COMPLETED',

  SET_TOURS: 'SET_TOURS',
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
    currentUserViewPermissionsByViewId: state => viewPermissionsGroupedPermissions(state.currentUser.permissions),
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
    [types.SET_TOURS](state, tours) {
      state.currentUser.ui_tours = tours;
    },
  },
  actions: {
    /**
     * POSTs login credentials, stores the returned access token, and runs post-login loading.
     * On failure, logs the error, dispatches a local logout, and rethrows.
     * @param {import('vuex').ActionContext} context
     * @param {Object} credentials - Payload forwarded to the login API (e.g. username, password)
     * @returns {Promise<void>}
     */
    async login({ commit, dispatch }, credentials) {
      try {
        const { access_token: accessToken } = await request.post(
          API_ROUTES.login,
          credentials,
          { fullResponse: true },
        );

        await dispatch('applyAccessToken', accessToken);
      } catch (err) {
        console.error(err.data);
        commit(types.LOGOUT);

        throw err;
      }
    },

    /**
     * Persists the access token, marks the session as logged in, and loads current user and file access in parallel.
     * @param {import('vuex').ActionContext} context
     * @param {string} accessToken - JWT (or other) access token from login/refresh
     * @returns {Promise<[void, void]>}
     */
    applyAccessToken({ commit, dispatch }, accessToken) {
      localStorageService.set(LOCAL_STORAGE_ACCESS_TOKEN_KEY, accessToken);

      commit(types.LOGIN_COMPLETED);

      return Promise.all([
        dispatch('fetchCurrentUser'),
        dispatch('filesAccess'),
      ]);
    },

    /**
     * Fetches file-access rights for the current session (used after successful authentication).
     * @returns {Promise<*>}
     */
    filesAccess() {
      return request.get(API_ROUTES.fileAccess);
    },

    /**
     * Loads the current user from the API and commits them to state.
     * No-ops with logout if not logged in. On certain HTTP statuses, dispatches `logout` before rethrowing.
     * @param {import('vuex').ActionContext} context
     * @returns {Promise<void>}
     */
    async fetchCurrentUser({ commit, dispatch, state }) {
      if (!state.isLoggedIn) {
        commit(types.LOGOUT);

        return;
      }

      try {
        commit(types.FETCH_USER);

        const currentUser = await request.get(API_ROUTES.currentUser, { fullResponse: true });

        commit(types.FETCH_USER_COMPLETED, currentUser);

        await dispatch('afterFetchCurrentUser', currentUser);
      } catch (err) {
        if (EXCLUDED_SERVER_ERROR_STATUSES.includes(err.status)) {
          dispatch('logout');
        }

        throw err;
      }
    },

    /**
     * Post-processing after the current user is loaded: blocks API-only users (logout + warning) or applies UI locale.
     * @param {import('vuex').ActionContext} context
     * @param {Object} currentUser - User object returned from `API_ROUTES.currentUser`
     */
    afterFetchCurrentUser({ dispatch }, currentUser) {
      if (isUserHasOnlyApiRole(currentUser)) {
        localStorageService.set(LOCAL_STORAGE_WARNING_POPUP_KEY, i18n.t('warnings.userDoesNotHaveUiRole'));

        dispatch('logout');

        return;
      }

      if (currentUser.ui_language) {
        dispatch('i18n/setPersonalLocale', currentUser.ui_language, { root: true });
      } else {
        dispatch('i18n/setDefaultLocale', DEFAULT_LOCALE, { root: true });
      }
    },

    /**
     * Calls the logout API, clears auth state and token, optionally runs a redirect,
     * then reloads the page after a short delay.
     * @param {import('vuex').ActionContext} context
     * @param {Object} [options]
     * @param {function(): (void|Promise<void>)} [options.redirect] - After local logout, before the delayed reload
     * @returns {Promise<void>}
     */
    async logout({ commit }, { redirect } = {}) {
      try {
        await request.post(API_ROUTES.logout);

        commit(types.LOGOUT);
        localStorageService.remove(LOCAL_STORAGE_ACCESS_TOKEN_KEY);

        if (redirect) {
          await redirect();
        }
      } catch (err) {
        console.error(err);
      } finally {
        /**
         * We've added timeout for the correct layout padding displaying with transition.
         * And we've added location.reload for refreshing every js objects (store, components states and etc.)
         */
        setTimeout(() => window.location.reload(), VUETIFY_ANIMATION_DELAY);
      }
    },

    /**
     * Returns the number of currently logged-in users (direct request; does not read or update Vuex state).
     * @returns {Promise<*>}
     */
    fetchLoggedUsersCountWithoutStore() {
      return request.get(API_ROUTES.loggedUserCount);
    },
  },
};
