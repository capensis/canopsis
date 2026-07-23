import { isArray, isEmpty, isFunction, pick } from 'lodash';

import { CRUD_ACTIONS, POPUP_TYPES } from '@/constants';

import { i18n } from '@/i18n';
import { usePopupsStore } from '@/plugins/popups/stores';
import { useAuthStore } from '@/stores/modules/auth';
import { useInfo } from '@/hooks/store/modules/info';

import { checkUserAnyAccess, groupedPermissionToPermission } from '@/helpers/permission';

const resolvePermissionId = (requiresPermission, route) =>
  isFunction(requiresPermission.id) ? requiresPermission.id(route) : requiresPermission.id;

export const checkAppInfoAccessForRoute = async (route = {}) => {
  const { requiresPermission } = route.meta;

  if (!requiresPermission) {
    return true;
  }

  const popupsStore = usePopupsStore();
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

  popupsStore.error({ text: i18n.t('common.notFound') });

  throw new Error('Application does not have access to the page');
};

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
    text: i18n.t('common.forbidden'),
    type: POPUP_TYPES.error,
  });

  throw new Error(
    `User does not have access to page '${route.path}' with permission id '${permissionId}'`,
  );
};
