import { get, isUndefined, omit, sortBy, keyBy } from 'lodash';
import flatten from 'flat';

import {
  NOT_COMPLETED_USER_PERMISSIONS,
  USERS_PERMISSIONS,
  PERMISSIONS_TYPES_TO_ACTIONS,
  USER_PERMISSIONS_PREFIXES, CRUD_ACTIONS,
} from '@/constants';

/**
 * Check user access for a permission
 *
 * @param {Object} permissionsById
 * @param {string} permissionId
 * @param {string} action
 * @returns {boolean}
 */
export function checkUserAccess(permissionsById, permissionId, action) {
  const permission = permissionsById[permissionId];

  if (permission && permission.actions) {
    const { actions } = permission;

    return action === CRUD_ACTIONS.can
      ? actions.length >= 0
      : actions.includes(action);
  }

  return false;
}

/**
 * Get actions for permission by type
 *
 * @param {Object} permission
 * @returns {*}
 */
export function getPermissionActions(permission) {
  return PERMISSIONS_TYPES_TO_ACTIONS[permission.type]
    ? PERMISSIONS_TYPES_TO_ACTIONS[permission.type]
    : [CRUD_ACTIONS.can];
}

/**
 * Get right checkbox value for the form
 *
 * @param {string} permissionId
 * @param {PermissionForm} permissions
 * @param {PermissionForm} changedPermissions
 * @param {string} [action=CRUD_ACTIONS.can]
 * @returns {boolean}
 */
export const getCheckboxValue = (
  permissionId,
  permissions,
  changedPermissions,
  action = CRUD_ACTIONS.can,
) => {
  const actions = permissions[permissionId] || [];

  const changedActions = get(changedPermissions, [permissionId]);
  const currentActions = isUndefined(changedActions) ? actions : changedActions;

  return currentActions.includes(action);
};

/**
 * Get prepared grouped permissions for the permissions page
 *
 * @param {Array<Object>} permissions
 * @param {Array<Object>} [views = []]
 * @param {Array<Object>} [playlists = []]
 * @returns {*}
 */
export function getGroupedPermissions(permissions, views = [], playlists = []) {
  const allBusinessPermissionsIds = flatten(USERS_PERMISSIONS.business);

  const { exploitation: exploitationTechnicalPermissions, ...adminTechnicalPermissions } = USERS_PERMISSIONS.technical;
  const adminTechnicalPermissionsValues = Object.values(adminTechnicalPermissions);
  const exploitationTechnicalPermissionsValues = Object.values(exploitationTechnicalPermissions);
  const apiPermissionsValues = Object.values(USERS_PERMISSIONS.api);
  const viewsById = keyBy(views, '_id');
  const playlistsById = keyBy(playlists, '_id');

  const groupedPermissions = permissions.reduce((acc, permission) => {
    const permissionId = String(permission._id, '\'');
    const view = viewsById[permissionId];
    const playlist = playlistsById[permissionId];

    if (view) {
      acc.view.push(permission);
    } else if (playlist) {
      acc.playlist.push(permission);
    } else if (adminTechnicalPermissionsValues.includes(permissionId)) {
      acc.technical.admin.push(permission);
    } else if (exploitationTechnicalPermissionsValues.includes(permissionId)) {
      acc.technical.exploitation.push(permission);
    } else if (
      Object.values(allBusinessPermissionsIds).includes(permissionId) ||
      NOT_COMPLETED_USER_PERMISSIONS.some(id => permissionId.startsWith(id))
    ) {
      const [parentKey] = permission._id.split('_');

      if (acc.business[parentKey]) {
        acc.business[parentKey].push(permission);
      }
    } else if (apiPermissionsValues.includes(permissionId)) {
      acc.api.push(permission);
    }

    return acc;
  }, {
    business: {
      [USER_PERMISSIONS_PREFIXES.business.common]: [],
      [USER_PERMISSIONS_PREFIXES.business.alarmsList]: [],
      [USER_PERMISSIONS_PREFIXES.business.context]: [],
      [USER_PERMISSIONS_PREFIXES.business.serviceWeather]: [],
      [USER_PERMISSIONS_PREFIXES.business.counter]: [],
      [USER_PERMISSIONS_PREFIXES.business.testingWeather]: [],
    },
    view: [],
    playlist: [],
    technical: {
      admin: [],
      exploitation: [],
    },
    api: [],
  });

  /**
   * We are using order which one we've defined on the reduce accumulator initial value.
   * For not `number`/`number string` object keys ordering is staying like we define
   */
  groupedPermissions.business = Object.entries(groupedPermissions.business)
    .map(([key, groupPermissions]) => ({
      key: `permissions.business.${key}`,
      permissions: sortBy(groupPermissions, ['description']),
    }));

  /**
   * Ordering behavior have the same behavior as for `business`
   */
  groupedPermissions.technical = Object.entries(groupedPermissions.technical)
    .map(([key, groupPermissions]) => ({
      key: `permissions.technical.${key}`,
      permissions: sortBy(groupPermissions, ['description']),
    }));

  groupedPermissions.view.push(...groupedPermissions.playlist);

  return omit(groupedPermissions, ['playlist']);
}