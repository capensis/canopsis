import { isEmpty, isFunction, isMatch } from 'lodash';

import {
  USER_PERMISSIONS_TO_PAGES_RULES,
  POPUP_TYPES, CRUD_ACTIONS,
} from '@/constants';

import store from '@/store';
import i18n from '@/i18n';

import { checkUserAccess } from '@/helpers/permission';

/**
 * Get app info value promise by key
 *
 * @param {string} key
 * @returns {Promise<unknown>|*}
 */
export function getAppInfoValuePromiseByKey(key) {
  const getterKey = `info/${key}`;

  if (store.getters[getterKey]) {
    return Promise.resolve(store.getters[getterKey]);
  }

  return store.watchOnce(state => state.info[key], v => !isEmpty(v));
}

/**
 * Check app info access for route
 *
 * @param {Object} [to={}] - Route object
 * @returns {Promise<boolean>}
 */
export async function checkAppInfoAccessForRoute(to = {}) {
  const { requiresRight } = to.meta;

  if (!requiresRight) {
    return true;
  }

  const rightId = isFunction(requiresRight.id) ? requiresRight.id(to) : requiresRight.id;
  const rightAppInfoRules = USER_PERMISSIONS_TO_PAGES_RULES[rightId];

  if (!rightAppInfoRules) {
    return true;
  }

  const [edition, stack] = await Promise.all([
    getAppInfoValuePromiseByKey('edition'),
    getAppInfoValuePromiseByKey('stack'),
  ]);

  const appInfo = {
    edition,
    stack,
  };

  if (isMatch(appInfo, rightAppInfoRules)) {
    return true;
  }

  store.dispatch('popups/error', { text: i18n.t('common.notFound') });

  throw new Error('Application don\'t have access to the page');
}

/**
 * Check user access for route
 *
 * @param {Object} [to={}] - Route object
 * @returns {Promise<boolean>}
 */
export async function checkUserAccessForRoute(to = {}) {
  const { requiresRight, requiresLogin } = to.meta;

  if (!requiresLogin || !requiresRight) {
    return true;
  }

  const permissionId = isFunction(requiresRight.id) ? requiresRight.id(to) : requiresRight.id;
  const permissionAction = requiresRight.action ? requiresRight.action : CRUD_ACTIONS.read;

  let currentUser = store.getters['auth/currentUser'];

  if (isEmpty(currentUser)) {
    currentUser = await store.watchOnce(state => state.auth.currentUser, v => !isEmpty(v));
  }

  if (!store.getters['info/popupTimeout']) {
    await store.watchOnce(state => state.info.popupTimeout, v => v);
  }

  const currentUserPermissionsById = store.getters['auth/currentUserPermissionsById'];

  if (checkUserAccess(currentUserPermissionsById, permissionId, permissionAction)) {
    return true;
  }

  store.dispatch('popups/add', {
    text: i18n.t('common.forbidden'),
    type: POPUP_TYPES.error,
  });

  throw new Error('User don\'t have access to page');
}

/**
 * Get path array for keep alive requests by route
 *
 * @param {string} path
 * @param {Object} query
 * @returns {Array}
 */
export function getKeepalivePathByRoute({ path, query } = {}) {
  const { tabId } = query;

  if (tabId) {
    return [path, tabId];
  }

  return [path];
}