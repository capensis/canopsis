import { isArray, isEmpty, isFunction, pick } from 'lodash';

import { CRUD_ACTIONS, POPUP_TYPES } from '@/constants';

import { usePopupsStore } from '@/plugins/popups/stores';

import { i18n } from '@/i18n';

import { useAuthStore } from '@/stores/modules/auth';

import { checkUserAnyAccess, groupedPermissionToPermission } from '@/helpers/permission';

import { useInfo } from '@/hooks/store/modules/info';

/**
 * @typedef {Object} RouteRequiresPermission
 * @property {string|string[]|function(import('vue-router').RouteLocationNormalized): string|string[]} id -
 *   Permission id, grouped ids, or a function of the target route (e.g. view kiosk `params.id`).
 * @property {string} [action] - CRUD action to check; defaults to `CRUD_ACTIONS.read` in user access checks.
 */

/**
 * Resolves `meta.requiresPermission.id`: calls it with `route` when it is a function, otherwise returns the
 * static id (string or grouped-permission array).
 *
 * @param {RouteRequiresPermission} requiresPermission - Route meta permission descriptor.
 * @param {import('vue-router').RouteLocationNormalized} route - Target route from the navigation guard.
 * @returns {string|string[]} Permission id or grouped ids.
 */
const resolvePermissionId = (requiresPermission, route) => isFunction(requiresPermission.id) ? requiresPermission.id(route) : requiresPermission.id;

/**
 * Guards a route by app edition: fetches app info when missing, then checks
 * `checkAppInfoAccessByPermission`. Used in `router.beforeResolve`.
 *
 * Routes without `meta.requiresPermission` pass. On denial, shows a not-found popup and throws.
 *
 * @param {import('vue-router').RouteLocationNormalized} [route={}] - Target route (`to`).
 * @returns {Promise<true>} Resolves when the edition allows the page.
 * @throws {Error} When the application edition does not include this page.
 */
export const checkAppInfoAccessForRoute = async (route = {}) => {
  const { requiresPermission } = route.meta;

  if (!requiresPermission) {
    return true;
  }

  const { edition, fetchAppInfo, checkAppInfoAccessByPermission } = useInfo();

  let permissionId = resolvePermissionId(requiresPermission, route);

  if (isArray(permissionId)) {
    permissionId = groupedPermissionToPermission(permissionId);
  }

  if (!edition.value) {
    await fetchAppInfo();
  }

  if (checkAppInfoAccessByPermission(permissionId)) {
    return true;
  }

  const popupsStore = usePopupsStore();

  popupsStore.error({ text: i18n.global.t('common.notFound') });

  throw new Error('Application does not have access to the page');
};

/**
 * Guards a logged-in route by user permission: fetches the current user when needed, then checks
 * `checkUserAnyAccess` for `requiresPermission.action` (default `CRUD_ACTIONS.read`). Used in
 * `router.beforeResolve`.
 *
 * Routes without `meta.requiresLogin` or `meta.requiresPermission` pass. On denial, shows a forbidden
 * popup and throws.
 *
 * @param {import('vue-router').RouteLocationNormalized} [route={}] - Target route (`to`).
 * @returns {Promise<true>} Resolves when the user may open the page.
 * @throws {Error} When the user lacks the required permission for `route.path`.
 */
export const checkUserAccessForRoute = async (route = {}) => {
  const { requiresPermission, requiresLogin } = route.meta;

  if (!requiresLogin || !requiresPermission) {
    return true;
  }

  const authStore = useAuthStore();
  const popupsStore = usePopupsStore();
  const permissionAction = requiresPermission.action || CRUD_ACTIONS.read;
  let permissionId = resolvePermissionId(requiresPermission, route);

  if (isEmpty(authStore.currentUser) && authStore.isLoggedIn) {
    await authStore.fetchCurrentUser();
  }

  if (!isArray(permissionId)) {
    permissionId = [permissionId];
  }

  if (
    checkUserAnyAccess(pick(authStore.currentUserPermissionsById, permissionId), permissionAction)
  ) {
    return true;
  }

  popupsStore.add({
    text: i18n.global.t('common.forbidden'),
    type: POPUP_TYPES.error,
  });

  throw new Error(
    `User does not have access to page '${route.path}' with permission id '${permissionId}'`,
  );
};
