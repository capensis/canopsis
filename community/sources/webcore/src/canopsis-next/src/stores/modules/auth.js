import { keyBy } from 'lodash';
import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import {
  API_ROUTES,
  VUETIFY_ANIMATION_DELAY,
  LOCAL_STORAGE_ACCESS_TOKEN_KEY,
  LOCAL_STORAGE_WARNING_POPUP_KEY,
  DEFAULT_LOCALE,
} from '@/config';
import { EXCLUDED_SERVER_ERROR_STATUSES } from '@/constants';

import localStorageService from '@/services/local-storage';
import request from '@/services/request';

import { i18n } from '@/i18n';

import { isUserHasOnlyApiRole } from '@/helpers/entities/user/entity';
import { viewPermissionsGroupedPermissions } from '@/helpers/permission';

import { useI18nStore } from './i18n';

export const useAuthStore = defineStore('auth', () => {
  const i18nStore = useI18nStore();

  const isLoggedIn = ref(localStorageService.has(LOCAL_STORAGE_ACCESS_TOKEN_KEY));
  const currentUser = ref({});
  const pending = ref(true);

  const currentUserPermissionsById = computed(() => keyBy(currentUser.value.permissions, '_id'));

  const currentUserViewPermissionsByViewId = computed(() =>
    viewPermissionsGroupedPermissions(currentUser.value.permissions),
  );

  const localLogout = () => {
    isLoggedIn.value = false;
    currentUser.value = {};
    pending.value = false;
  };

  const applyAccessToken = (accessToken, afterFetchCallback) => {
    localStorageService.set(LOCAL_STORAGE_ACCESS_TOKEN_KEY, accessToken);

    isLoggedIn.value = true;

    return Promise.all([fetchCurrentUser(afterFetchCallback), filesAccess()]);
  };

  const login = async (credentials, afterFetchCallback) => {
    try {
      const { access_token: accessToken } = await request.post(API_ROUTES.login, credentials, {
        fullResponse: true,
      });

      await applyAccessToken(accessToken, afterFetchCallback);
    } catch (err) {
      console.error(err);

      localLogout();

      throw err;
    }
  };

  const filesAccess = () => request.get(API_ROUTES.fileAccess);

  const fetchCurrentUser = async () => {
    if (!isLoggedIn.value) {
      localLogout();

      return;
    }

    try {
      pending.value = true;

      const user = await request.get(API_ROUTES.currentUser, {
        fullResponse: true,
      });

      currentUser.value = user;
      pending.value = false;

      await afterFetchCurrentUser(user);
    } catch (err) {
      if (EXCLUDED_SERVER_ERROR_STATUSES.includes(err.status)) {
        logout();
      }

      throw err;
    }
  };

  const afterFetchCurrentUser = async (user) => {
    if (isUserHasOnlyApiRole(user)) {
      localStorageService.set(
        LOCAL_STORAGE_WARNING_POPUP_KEY,
        i18n.global.t('warnings.userDoesNotHaveUiRole'),
      );

      logout();

      return;
    }

    if (currentUser.value.ui_language) {
      i18nStore.setPersonalLocale(currentUser.value.ui_language);
    } else {
      i18nStore.setDefaultLocale(DEFAULT_LOCALE);
    }
  };

  const setCurrentUser = (user = {}) => {
    currentUser.value = user;
  };

  const logout = async ({ redirect } = {}) => {
    try {
      await request.post(API_ROUTES.logout);

      localLogout();
      localStorageService.remove(LOCAL_STORAGE_ACCESS_TOKEN_KEY);

      if (redirect) {
        await redirect();
      }
    } catch (err) {
      console.error(err);
    } finally {
      setTimeout(() => window.location.reload(), VUETIFY_ANIMATION_DELAY);
    }
  };

  const setTours = tours => (currentUser.value.ui_tours = tours);

  const fetchLoggedUsersCountWithoutStore = () => request.get(API_ROUTES.loggedUserCount);

  return {
    isLoggedIn,
    currentUser,
    pending,

    currentUserPermissionsById,
    currentUserViewPermissionsByViewId,

    login,
    logout,
    localLogout,
    applyAccessToken,
    fetchCurrentUser,
    afterFetchCurrentUser,
    setCurrentUser,
    filesAccess,
    setTours,
    fetchLoggedUsersCountWithoutStore,
  };
});
