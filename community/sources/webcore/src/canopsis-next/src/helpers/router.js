import {
  isEmpty,
  isFunction,
  isMatch,
  isArray,
  pick,
} from 'lodash';

import { USER_PERMISSIONS_TO_PAGES_RULES, POPUP_TYPES, CRUD_ACTIONS } from '@/constants';

import store from '@/store';
import i18n from '@/i18n';

import { checkUserAnyAccess } from '@/helpers/entities/permissions/list';
import { groupedPermissionToPermission } from '@/helpers/permission';

/**
 * Get app info value promise by key
 *
 * @param {string} key
 * @returns {Promise<unknown>|*}
 */
export const getAppInfoValuePromiseByKey = (key) => {
  const getterKey = `info/${key}`;

  if (store.getters[getterKey]) {
    return Promise.resolve(store.getters[getterKey]);
  }

  return store.watchOnce(state => state.info.appInfo[key], v => !isEmpty(v));
};

/**
 * Check app info access for route
 *
 * @param {Object} [to={}] - Route object
 * @returns {Promise<boolean>}
 */
export async function checkAppInfoAccessForRoute(to = {}) {
  const { requiresPermission } = to.meta;

  if (!requiresPermission) {
    return true;
  }

  let permissionId = isFunction(requiresPermission.id) ? requiresPermission.id(to) : requiresPermission.id;

  if (isArray(permissionId)) {
    permissionId = groupedPermissionToPermission(permissionId);
  }

  const permissionAppInfoRules = USER_PERMISSIONS_TO_PAGES_RULES[permissionId];

  if (!permissionAppInfoRules) {
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

  if (isMatch(appInfo, permissionAppInfoRules)) {
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
  const { requiresPermission, requiresLogin } = to.meta;

  if (!requiresLogin || !requiresPermission) {
    return true;
  }

  const permissionAction = requiresPermission.action ? requiresPermission.action : CRUD_ACTIONS.read;
  let permissionId = isFunction(requiresPermission.id) ? requiresPermission.id(to) : requiresPermission.id;

  const currentUser = store.getters['auth/currentUser'];

  if (isEmpty(currentUser)) {
    await store.watchOnce(state => state.auth.currentUser, v => !isEmpty(v));
  }

  if (!store.getters['info/popupTimeout']) {
    await store.watchOnce(state => state.info.popupTimeout, v => v);
  }

  const currentUserPermissionsById = store.getters['auth/currentUserPermissionsById'];

  if (!isArray(permissionId)) {
    permissionId = [permissionId];
  }

  if (checkUserAnyAccess(pick(currentUserPermissionsById, permissionId), permissionAction)) {
    return true;
  }

  store.dispatch('popups/add', {
    text: i18n.t('common.forbidden'),
    type: POPUP_TYPES.error,
  });

  throw new Error(`User doesn't have access to page '${to.path}' with permission id '${permissionId}'`);
}
